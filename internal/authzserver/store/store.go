// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

import (
	"context"
	"errors"
	"sync"
	"time"

	pb "github.com/marmotedu/api/proto/apiserver/v1"
	"github.com/ory/ladon"

	"github.com/marmotedu/iam/pkg/log"
)

// ErrSecretNotFound defines secret not found error.
var ErrSecretNotFound = errors.New("secret not found")

var secrets = make(map[string]*pb.SecretInfo)
var policies = make(map[string][]*ladon.DefaultPolicy)

var reloadMu sync.Mutex

// CacheService defines a cache service which can load secrets and policies timing.
type CacheService struct {
	ctx      context.Context
	addr     string
	clientCA string
}

// GetSecret returns the secret information of the given secret.
func GetSecret(secretID string) (*pb.SecretInfo, error) {
	secret, ok := secrets[secretID]
	if !ok {
		return nil, ErrSecretNotFound
	}

	return secret, nil
}

// New create a new cache service instance by the given configuration.
func New(ctx context.Context, addr, clientCA string) *CacheService {
	return &CacheService{
		ctx:      ctx,
		addr:     addr,
		clientCA: clientCA,
	}
}

// Start start a cache service.
func (s *CacheService) Start() {
	grpcClient := GrpcClient{
		Addr:     s.addr,
		ClientCA: s.clientCA,
	}
	grpcClient.Connect()

	go startPubSubLoop()
	// 1s is the minimum amount of time between hot reloads. The
	// interval counts from the start of one reload to the next.
	go reloadLoop(s.ctx)
	go reloadQueueLoop(s.ctx)
	DoReload()
}

// shouldReload returns true if we should perform any reload. Reloads happens if
// we have reload callback queued.
func shouldReload() ([]func(), bool) {
	requeueLock.Lock()
	defer requeueLock.Unlock()
	if len(requeue) == 0 {
		return nil, false
	}
	n := requeue
	requeue = []func(){}
	return n, true
}

func reloadLoop(ctx context.Context, complete ...func()) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		// We don't check for reload right away as the gateway peroms this on the
		// startup sequence. We expect to start checking on the first tick after the
		// gateway is up and running.
		case <-ticker.C:
			cb, ok := shouldReload()
			if !ok {
				continue
			}
			start := time.Now()
			log.Info("reload: initiating")
			DoReload()
			log.Info("reload: complete")
			for _, c := range cb {
				// most of the callbacks are nil, we don't want to execute nil functions to
				// avoid panics.
				if c != nil {
					c()
				}
			}
			if len(complete) != 0 {
				complete[0]()
			}
			log.Infof("reload: cycle completed in %v", time.Since(start))
		}
	}
}

// reloadQueue is used by reloadURLStructure to queue a reload. It's not
// buffered, as reloadQueueLoop should pick these up immediately.
var reloadQueue = make(chan func())

var requeueLock sync.Mutex

// This is a list of callbacks to execute on the next reload. It is protected by
// requeueLock for concurrent use.
var requeue []func()

func reloadQueueLoop(ctx context.Context, cb ...func()) {
	for {
		select {
		case <-ctx.Done():
			return
		case fn := <-reloadQueue:
			requeueLock.Lock()
			requeue = append(requeue, fn)
			requeueLock.Unlock()
			log.Info("Reload queued")
			if len(cb) != 0 {
				cb[0]()
			}
		}
	}
}

// done will be called when the reload is finished. Note that if a
// reload is already queued, another won't be queued, but done will
// still be called when said queued reload is finished.
func reloadURLStructure(done func()) {
	reloadQueue <- done
}

// DoReload reload secrets and policies.
func DoReload() {
	reloadMu.Lock()
	defer reloadMu.Unlock()

	grpcClient := &GrpcClient{}
	if !grpcClient.Connect() {
		log.Error("Failed connecting to grpc server")
		return
	}

	var err error
	secrets, err = grpcClient.GetSecrets()
	if err != nil {
		log.Errorf("Error during syncing secrets: %s", err.Error())
		return
	}

	policies, err = grpcClient.GetPolicies()
	if err != nil {
		log.Errorf("Error during syncing policies: %s", err.Error())
		return
	}

	log.Info("Secrets and policies reload complete")
}
