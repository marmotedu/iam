// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package watcher

import (
	"context"
	"fmt"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/robfig/cron/v3"

	genericoptions "github.com/marmotedu/iam/internal/pkg/options"
	"github.com/marmotedu/iam/internal/watcher/options"
	"github.com/marmotedu/iam/internal/watcher/watcher"

	// trigger init functions in `internal/watcher/watcher/`.
	_ "github.com/marmotedu/iam/internal/watcher/watcher/all"
	"github.com/marmotedu/iam/pkg/log"
	"github.com/marmotedu/iam/pkg/log/cronlog"
)

type watchJob struct {
	*cron.Cron
	config *options.WatcherOptions
	rs     *redsync.Redsync
}

func newWatchJob(redisOptions *genericoptions.RedisOptions, watcherOptions *options.WatcherOptions) *watchJob {
	logger := cronlog.NewLogger(log.SugaredLogger())

	client := goredislib.NewClient(&goredislib.Options{
		Addr:     fmt.Sprintf("%s:%d", redisOptions.Host, redisOptions.Port),
		Username: redisOptions.Username,
		Password: redisOptions.Password,
	})

	rs := redsync.New(goredis.NewPool(client))

	cronjob := cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(logger), cron.Recover(logger)),
	)

	return &watchJob{
		Cron:   cronjob,
		config: watcherOptions,
		rs:     rs,
	}
}

func (w *watchJob) addWatchers() *watchJob {
	for name, watch := range watcher.ListWatchers() {
		// log with `{"watcher": "counter"}` key-value to distinguish which watcher the log comes from.
		//nolint: golint,staticcheck
		ctx := context.WithValue(context.Background(), log.KeyWatcherName, name)

		if err := watch.Init(ctx, w.rs.NewMutex(name, redsync.WithExpiry(2*time.Hour)), w.config); err != nil {
			log.Panicf("construct watcher %s failed: %s", name, err.Error())
		}

		_, _ = w.AddJob(watch.Spec(), watch)
	}

	return w
}
