// Copyright 2020 Lingfei Kong <marmotedu@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package options contains flags and options for initializing an apiserver
package options

import (
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
	"github.com/marmotedu/component-base/pkg/json"

	genericoptions "github.com/marmotedu/iam/internal/pkg/options"
	"github.com/marmotedu/iam/pkg/log"
)

// CleanOptions defines options for clean watcher.
type CleanOptions struct {
	MaxReserveDays int `json:"max-reserve-days" mapstructure:"max-reserve-days"`
}

// TaskOptions defines options for task watcher.
type TaskOptions struct {
	MaxInactiveDays int `json:"max-inactive-days" mapstructure:"max-inactive-days"`
}

// WatcherOptions defines options for watchers.
type WatcherOptions struct {
	Clean CleanOptions `json:"clean" mapstructure:"clean"`
	Task  TaskOptions  `json:"task"  mapstructure:"task"`
}

// Options runs a pumpserver.
type Options struct {
	HealthCheckPath    string                       `json:"health-check-path"    mapstructure:"health-check-path"`
	HealthCheckAddress string                       `json:"health-check-address" mapstructure:"health-check-address"`
	MySQLOptions       *genericoptions.MySQLOptions `json:"mysql"                mapstructure:"mysql"`
	RedisOptions       *genericoptions.RedisOptions `json:"redis"                mapstructure:"redis"`
	WatcherOptions     *WatcherOptions              `json:"watcher"              mapstructure:"watcher"`
	Log                *log.Options                 `json:"log"                  mapstructure:"log"`
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	s := Options{
		HealthCheckPath:    "healthz",
		HealthCheckAddress: "0.0.0.0:5050",
		MySQLOptions:       genericoptions.NewMySQLOptions(),
		RedisOptions:       genericoptions.NewRedisOptions(),
		WatcherOptions: &WatcherOptions{
			Clean: CleanOptions{
				MaxReserveDays: 180, // default half a year
			},
			Task: TaskOptions{
				MaxInactiveDays: 0, // not expire by default
			},
		},
		Log: log.NewOptions(),
	}

	return &s
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.Log.AddFlags(fss.FlagSet("logs"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.StringVar(&o.HealthCheckPath, "health-check-path", o.HealthCheckPath, ""+
		"Specifies liveness health check request path.")
	fs.StringVar(&o.HealthCheckAddress, "health-check-address", o.HealthCheckAddress, ""+
		"Specifies liveness health check bind address.")

	fs.IntVar(
		&o.WatcherOptions.Clean.MaxReserveDays,
		"watcher.counter.max-reserve-days",
		o.WatcherOptions.Clean.MaxReserveDays,
		"Policy audit log maximum retention days.",
	)
	fs.IntVar(
		&o.WatcherOptions.Task.MaxInactiveDays,
		"watcher.task.max-inactive-days",
		o.WatcherOptions.Task.MaxInactiveDays,
		"Maximum user inactivity time. Otherwise the account will be disabled.",
	)

	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
