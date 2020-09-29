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
	"github.com/marmotedu/iam/pkg/storage"
	"github.com/marmotedu/log"
	"github.com/ory/ladon"
)

// ErrSecretNotFound defines secret not found error.
var ErrSecretNotFound = errors.New("secret not found")

var secrets = make(map[string]*pb.SecretInfo)
var policies = make(map[string][]*ladon.DefaultPolicy)

var reloadMu sync.Mutex

// CacheService defines a cache service which can load secrets and policies timing.
type CacheService struct {
	config   *storage.Config
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
func New(config *storage.Config, addr, clientCA string) *CacheService {
	return &CacheService{
		config:   config,
		addr:     addr,
		clientCA: clientCA,
	}
}

// Start start a cache service.
func (s *CacheService) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcClient := GrpcClient{
		Addr:     s.addr,
		ClientCA: s.clientCA,
	}
	grpcClient.Connect()

	go startPubSubLoop()
	// 1s is the minimum amount of time between hot reloads. The
	// interval counts from the start of one reload to the next.
	ticker := time.NewTicker(time.Second)
	go reloadLoop(ticker.C)
	go reloadQueueLoop()
	go storage.ConnectToRedis(ctx, s.config)

	DoReload()
}

// startReloadChan and reloadDoneChan are used by the two reload loops
// running in separate goroutines to talk. reloadQueueLoop will use
// startReloadChan to signal to reloadLoop to start a reload, and
// reloadLoop will use reloadDoneChan to signal back that it's done with
// the reload. Buffered simply to not make the goroutines block each
// other.
var startReloadChan = make(chan struct{}, 1)
var reloadDoneChan = make(chan struct{}, 1)

func reloadLoop(tick <-chan time.Time) {
	<-tick

	for range startReloadChan {
		log.Info("reload: initiating")
		DoReload()
		log.Info("reload: complete")
		reloadDoneChan <- struct{}{}

		<-tick
	}
}

// reloadQueue is used by reloadURLStructure to queue a reload. It's not
// buffered, as reloadQueueLoop should pick these up immediately.
var reloadQueue = make(chan func())

func reloadQueueLoop() {
	reloading := false

	var fns []func()

	for {
		select {
		case <-reloadDoneChan:
			for _, fn := range fns {
				fn()
			}

			fns = fns[:0]
			reloading = false
		case fn := <-reloadQueue:
			if fn != nil {
				fns = append(fns, fn)
			}

			if !reloading {
				log.Info("Reload queued")
				startReloadChan <- struct{}{}

				reloading = true
			} else {
				log.Info("Reload already queued")
			}
		}
	}
}

// reloadURLStructure will queue an API reload. The reload will
// eventually create a new muxer, reload all the app configs for an
// instance and then replace the DefaultServeMux with the new one. This
// enables a reconfiguration to take place without stopping any requests
// from being handled.
//
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

	client := &GrpcClient{}
	if !client.Connect() {
		log.Error("Failed connecting to grpc server")
		return
	}

	var err error
	secrets, err = client.GetSecrets()
	if err != nil {
		log.Errorf("Error during syncing secrets: %s", err.Error())
		return
	}

	policies, err = client.GetPolicies()
	if err != nil {
		log.Errorf("Error during syncing policies: %s", err.Error())
		return
	}

	log.Info("Secrets and policies reload complete")
}
