// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package task

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"

	"github.com/marmotedu/iam/internal/apiserver/store/mysql"
	"github.com/marmotedu/iam/internal/watcher/options"
	"github.com/marmotedu/iam/internal/watcher/watcher"
	"github.com/marmotedu/iam/pkg/log"
)

type taskWatcher struct {
	ctx             context.Context
	mutex           *redsync.Mutex
	maxInactiveDays int
}

// Run runs the watcher job.
func (tw *taskWatcher) Run() {
	if err := tw.mutex.Lock(); err != nil {
		log.L(tw.ctx).Info("taskWatcher already run.")

		return
	}
	defer func() {
		if _, err := tw.mutex.Unlock(); err != nil {
			log.L(tw.ctx).Errorf("could not release taskWatcher lock. err: %v", err)

			return
		}
	}()

	db, _ := mysql.GetMySQLFactoryOr(nil)

	users, err := db.Users().List(tw.ctx, metav1.ListOptions{})
	if err != nil {
		log.L(tw.ctx).Errorf("list user failed", "error", err)

		return
	}

	for _, user := range users.Items {
		// if maxInactiveDays equal to 0, means never forbid
		if tw.maxInactiveDays == 0 {
			continue
		}

		if time.Since(user.LoginedAt) > time.Duration(tw.maxInactiveDays)*(24*time.Hour) {
			log.L(tw.ctx).Infof("user %s not active for %d days, disable his account", user.Name, tw.maxInactiveDays)

			user.Status = 0
			_ = db.Users().Update(tw.ctx, user, metav1.UpdateOptions{})
		}
	}
}

// Spec is parsed using the time zone of task Cron instance as the default.
func (tw *taskWatcher) Spec() string {
	return "@every 1d"
}

// Init initializes the watcher for later execution.
func (tw *taskWatcher) Init(ctx context.Context, rs *redsync.Mutex, config interface{}) error {
	cfg, ok := config.(*options.WatcherOptions)
	if !ok {
		return watcher.ErrConfigUnavailable
	}

	*tw = taskWatcher{
		ctx:             ctx,
		mutex:           rs,
		maxInactiveDays: cfg.Task.MaxInactiveDays,
	}

	return nil
}

func init() {
	watcher.Register("task", &taskWatcher{})
}
