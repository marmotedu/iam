// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package options contains flags and options for initializing an apiserver
package options

import (
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
	"github.com/marmotedu/component-base/pkg/json"

	genericoptions "github.com/marmotedu/iam/internal/pkg/options"
	"github.com/marmotedu/iam/internal/pump/analytics"
	"github.com/marmotedu/iam/pkg/log"
)

// PumpConfig defines options for pump back-end.
type PumpConfig struct {
	Type                  string                     `json:"type"                    mapstructure:"type"`
	Filters               analytics.AnalyticsFilters `json:"filters"                 mapstructure:"filters"`
	Timeout               int                        `json:"timeout"                 mapstructure:"timeout"`
	OmitDetailedRecording bool                       `json:"omit-detailed-recording" mapstructure:"omit-detailed-recording"`
	Meta                  map[string]interface{}     `json:"meta"                    mapstructure:"meta"`
}

// Options runs a pumpserver.
type Options struct {
	PurgeDelay            int                          `json:"purge-delay"             mapstructure:"purge-delay"`
	Pumps                 map[string]PumpConfig        `json:"pumps"                   mapstructure:"pumps"`
	HealthCheckPath       string                       `json:"health-check-path"       mapstructure:"health-check-path"`
	HealthCheckAddress    string                       `json:"health-check-address"    mapstructure:"health-check-address"`
	OmitDetailedRecording bool                         `json:"omit-detailed-recording" mapstructure:"omit-detailed-recording"`
	RedisOptions          *genericoptions.RedisOptions `json:"redis"                   mapstructure:"redis"`
	Log                   *log.Options                 `json:"log"                     mapstructure:"log"`
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	s := Options{
		PurgeDelay: 10,
		Pumps: map[string]PumpConfig{
			"csv": {
				Type: "csv",
				Meta: map[string]interface{}{
					"csv_dir": "./analytics-data",
				},
			},
		},
		HealthCheckPath:    "healthz",
		HealthCheckAddress: "0.0.0.0:7070",
		RedisOptions:       genericoptions.NewRedisOptions(),
		Log:                log.NewOptions(),
	}

	return &s
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.Log.AddFlags(fss.FlagSet("logs"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.IntVar(&o.PurgeDelay, "purge-delay", o.PurgeDelay, ""+
		"This setting the purge delay (in seconds) when purge the data from Redis to MongoDB or other data stores.")
	fs.StringVar(&o.HealthCheckPath, "health-check-path", o.HealthCheckPath, ""+
		"Specifies liveness health check request path.")
	fs.StringVar(&o.HealthCheckAddress, "health-check-address", o.HealthCheckAddress, ""+
		"Specifies liveness health check bind address.")
	fs.BoolVar(&o.OmitDetailedRecording, "omit-detailed-recording", o.OmitDetailedRecording, ""+
		"Setting this to true will avoid writing policy fields for each authorization request in pumps.")

	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}
