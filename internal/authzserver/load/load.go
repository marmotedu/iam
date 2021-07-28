// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package load

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/marmotedu/iam/pkg/log"
	"github.com/marmotedu/iam/pkg/storage"
)

// Loader defines function to reload storage.
type Loader interface {
	Reload() error
}

// Load is used to reload given storage.
type Load struct {
	ctx    context.Context
	lock   *sync.RWMutex
	loader Loader
}

// NewLoader return a loader with a loader implement.
func NewLoader(ctx context.Context, loader Loader) *Load {
	return &Load{
		ctx:    ctx,
		lock:   new(sync.RWMutex),
		loader: loader,
	}
}

// Start start a loop service.
func (l *Load) Start() {
	go startPubSubLoop()
	go l.reloadQueueLoop()
	// 1s is the minimum amount of time between hot reloads. The
	// interval counts from the start of one reload to the next.
	go l.reloadLoop()
	l.DoReload()
}

func startPubSubLoop() {
	cacheStore := storage.RedisCluster{}
	cacheStore.Connect()
	// On message, synchronize
	for {
		err := cacheStore.StartPubSubHandler(RedisPubSubChannel, func(v interface{}) {
			handleRedisEvent(v, nil, nil)
		})
		if err != nil {
			if !errors.Is(err, storage.ErrRedisIsDown) {
				log.Errorf("Connection to Redis failed, reconnect in 10s: %s", err.Error())
			}

			time.Sleep(10 * time.Second)
			log.Warnf("Reconnecting: %s", err.Error())
		}
	}
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

func (l *Load) reloadLoop(complete ...func()) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-l.ctx.Done():
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
			l.DoReload()
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

// reloadQueue used to queue a reload. It's not
// buffered, as reloadQueueLoop should pick these up immediately.
var reloadQueue = make(chan func())

var requeueLock sync.Mutex

// This is a list of callbacks to execute on the next reload. It is protected by
// requeueLock for concurrent use.
var requeue []func()

func (l *Load) reloadQueueLoop(cb ...func()) {
	for {
		select {
		case <-l.ctx.Done():
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

// DoReload reload secrets and policies.
func (l *Load) DoReload() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if err := l.loader.Reload(); err != nil {
		log.Errorf("faild to refresh target storage: %s", err.Error())
	}

	log.Debug("refresh target storage succ")
}
