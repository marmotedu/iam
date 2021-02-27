// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package analytics

import (
	"github.com/spf13/pflag"
)

// AnalyticsOptions contains configuration items related to analytics.
type AnalyticsOptions struct {
	PoolSize                int    `json:"pool-size"                 mapstructure:"pool-size"`
	RecordsBufferSize       uint64 `json:"records-buffer-size"       mapstructure:"records-buffer-size"`
	StorageExpirationTime   int    `json:"storage-expiration-time"   mapstructure:"storage-expiration-time"`
	Enable                  bool   `json:"enable"                    mapstructure:"enable"`
	EnableDetailedRecording bool   `json:"enable-detailed-recording" mapstructure:"enable-detailed-recording"`
}

// NewAnalyticsOptions creates a AnalyticsOptions object with default parameters.
func NewAnalyticsOptions() *AnalyticsOptions {
	return &AnalyticsOptions{
		Enable:                  true,
		PoolSize:                50,
		RecordsBufferSize:       2000,
		EnableDetailedRecording: true,
		StorageExpirationTime:   60,
	}
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (o *AnalyticsOptions) Validate() []error {
	return []error{}
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (o *AnalyticsOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.BoolVar(&o.Enable, "analytics.enable", o.Enable,
		"Enable profiling via web interface host:port/debug/pprof/")

	fs.IntVar(&o.PoolSize, "analytics.pool-size", o.PoolSize,
		"Enable profiling via web interface host:port/debug/pprof/")

	fs.Uint64Var(&o.RecordsBufferSize, "analytics.records-buffer-size", o.RecordsBufferSize,
		"Enable profiling via web interface host:port/debug/pprof/")

	fs.BoolVar(&o.EnableDetailedRecording, "analytics.enable-detailed-recording", o.EnableDetailedRecording,
		"Enable profiling via web interface host:port/debug/pprof/")

	fs.IntVar(&o.StorageExpirationTime, "analytics.storage-expiration-time", o.StorageExpirationTime,
		"Enables metrics on the apiserver at /metrics")
}
