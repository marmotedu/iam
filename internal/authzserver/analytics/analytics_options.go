// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package analytics

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

// AnalyticsOptions contains configuration items related to analytics.
type AnalyticsOptions struct {
	PoolSize                int           `json:"pool-size"                 mapstructure:"pool-size"`
	RecordsBufferSize       uint64        `json:"records-buffer-size"       mapstructure:"records-buffer-size"`
	FlushInterval           uint64        `json:"flush-interval"            mapstructure:"flush-interval"`
	StorageExpirationTime   time.Duration `json:"storage-expiration-time"   mapstructure:"storage-expiration-time"`
	Enable                  bool          `json:"enable"                    mapstructure:"enable"`
	EnableDetailedRecording bool          `json:"enable-detailed-recording" mapstructure:"enable-detailed-recording"`
}

// NewAnalyticsOptions creates a AnalyticsOptions object with default parameters.
func NewAnalyticsOptions() *AnalyticsOptions {
	return &AnalyticsOptions{
		Enable:                  true,
		PoolSize:                50,
		RecordsBufferSize:       1000,
		FlushInterval:           200,
		EnableDetailedRecording: true,
		StorageExpirationTime:   time.Duration(24) * time.Hour,
	}
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (o *AnalyticsOptions) Validate() []error {
	if o == nil {
		return nil
	}
	errors := []error{}

	if o.Enable && (o.FlushInterval < 1 || o.FlushInterval > 1000) {
		errors = append(errors, fmt.Errorf("--analytics.flush-interval %v must be between 1 and 1000", o.FlushInterval))
	}

	return errors
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (o *AnalyticsOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.BoolVar(&o.Enable, "analytics.enable", o.Enable, ""+
		"This sets the iam-authz-server to record analytics data.")

	fs.IntVar(&o.PoolSize, "analytics.pool-size", o.PoolSize,
		"Specify number of pool workers.")

	fs.Uint64Var(&o.RecordsBufferSize, "analytics.records-buffer-size", o.RecordsBufferSize,
		"Specifies buffer size for pool workers (size of each pipeline operation).")

	fs.BoolVar(&o.EnableDetailedRecording, "analytics.enable-detailed-recording", o.EnableDetailedRecording,
		"Enable detailed analytics at the key level.")

	fs.DurationVar(&o.StorageExpirationTime, "analytics.storage-expiration-time", o.StorageExpirationTime, ""+
		"Set to a value larger than the Pump's purge_delay. "+
		"This allows the analytics data to exist long enough in Redis to be processed by the Pump.")
}
