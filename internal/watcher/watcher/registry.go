// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package watcher

import (
	"context"
	"errors"
	"sync"

	"github.com/go-redsync/redsync/v4"
	"github.com/robfig/cron/v3"
)

// IWatcher is the interface for watchers.
type IWatcher interface {
	Init(ctx context.Context, rs *redsync.Mutex, config interface{}) error
	Spec() string
	cron.Job
}

var (
	registryLock = new(sync.Mutex)
	registry     = make(map[string]IWatcher)
)

var (
	// ErrRegistered will be returned when watcher is already been registered.
	ErrRegistered = errors.New("watcher has already been registered")
	// ErrConfigUnavailable will be returned when the configuration input is not the expected type.
	ErrConfigUnavailable = errors.New("configuration is not available")
)

// Register registers a watcher and save in global variable `registry`.
func Register(name string, watcher IWatcher) {
	registryLock.Lock()
	defer registryLock.Unlock()

	if _, ok := registry[name]; ok {
		panic("duplicate watcher entry: " + name)
	}

	registry[name] = watcher
}

// FindMonitor looks up a watcher in the registry. If not found, nil is returned.
func FindMonitor(name string) IWatcher {
	registryLock.Lock()
	defer registryLock.Unlock()

	return registry[name]
}

// Unregister removes a watcher from the registry. The watcher should have stopped.
func Unregister(name string) {
	registryLock.Lock()
	defer registryLock.Unlock()

	delete(registry, name)
}

// ListWatchers returns registered watchers in map format.
func ListWatchers() map[string]IWatcher {
	registryLock.Lock()
	defer registryLock.Unlock()

	return registry
}
