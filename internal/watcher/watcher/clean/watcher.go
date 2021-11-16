// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package clean

import (
	"context"

	"github.com/go-redsync/redsync/v4"

	"github.com/marmotedu/iam/internal/apiserver/store/mysql"
	"github.com/marmotedu/iam/internal/watcher/options"
	"github.com/marmotedu/iam/internal/watcher/watcher"
	"github.com/marmotedu/iam/pkg/log"
)

type cleanWatcher struct {
	ctx            context.Context
	mutex          *redsync.Mutex
	maxReserveDays int
}

// Run runs the watcher job.
func (cw *cleanWatcher) Run() {
	if err := cw.mutex.Lock(); err != nil {
		log.L(cw.ctx).Info("cleanWatcher already run.")

		return
	}

	defer func() {
		if _, err := cw.mutex.Unlock(); err != nil {
			log.L(cw.ctx).Errorf("could not release cleanWatcher lock. err: %v", err)

			return
		}
	}()

	db, _ := mysql.GetMySQLFactoryOr(nil)

	rowsAffected, err := db.PolicyAudits().ClearOutdated(cw.ctx, cw.maxReserveDays)
	if err != nil {
		log.L(cw.ctx).Errorw("clean data from policy_audit failed", "error", err)

		return
	}

	log.L(cw.ctx).Debugf("clean data from policy_audit succ, %d rows affected", rowsAffected)
}

// Spec is parsed using the time zone of clean Cron instance as the default.
func (cw *cleanWatcher) Spec() string {
	return "@every 1d"
}

// Init initializes the watcher for later execution.
func (cw *cleanWatcher) Init(ctx context.Context, rs *redsync.Mutex, config interface{}) error {
	cfg, ok := config.(*options.WatcherOptions)
	if !ok {
		return watcher.ErrConfigUnavailable
	}

	*cw = cleanWatcher{
		ctx:            ctx,
		mutex:          rs,
		maxReserveDays: cfg.Clean.MaxReserveDays,
	}

	return nil
}

func init() {
	watcher.Register("clean", &cleanWatcher{})
}
