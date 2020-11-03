// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package options contains flags and options for initializing an apiserver
package options

import (
	"encoding/json"

	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
	genericoptions "github.com/marmotedu/iam/internal/pkg/options"
	"github.com/marmotedu/iam/internal/pump/analytics"
	"github.com/marmotedu/log"
)

// PumpConfig defines options for pump back-end.
type PumpConfig struct {
	Type                  string                     `json:"type" mapstructure:"type"`
	Filters               analytics.AnalyticsFilters `json:"filters" mapstructure:"filters"`
	Timeout               int                        `json:"timeout" mapstructure:"timeout"`
	OmitDetailedRecording bool                       `json:"omit-detailed-recording" mapstructure:"omit-detailed-recording"`
	Meta                  map[string]interface{}     `json:"meta" mapstructure:"meta"`
}

// PumpOptions runs a pumpserver.
type PumpOptions struct {
	PumpConfig            string                       `json:"pumpconfig" mapstructure:"-"`
	PurgeDelay            int                          `json:"purge-delay" mapstructure:"purge-delay"`
	Pumps                 map[string]PumpConfig        `json:"pumps" mapstructure:"pumps"`
	HealthCheckPath       string                       `json:"health-check-path" mapstructure:"health-check-path"`
	HealthCheckAddress    string                       `json:"health-check-address" mapstructure:"health-check-address"`
	OmitDetailedRecording bool                         `json:"omit-detailed-recording" mapstructure:"omit-detailed-recording"`
	RedisOptions          *genericoptions.RedisOptions `json:"redis" mapstructure:"redis"`
	Log                   *log.Options                 `json:"log" mapstructure:"log"`
}

// NewPumpOptions creates a new PumpOptions object with default parameters.
func NewPumpOptions() *PumpOptions {
	s := PumpOptions{
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
func (s *PumpOptions) Flags() (fss cliflag.NamedFlagSets) {
	s.RedisOptions.AddFlags(fss.FlagSet("redis"))
	s.Log.AddFlags(fss.FlagSet("logs"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.StringVar(&s.PumpConfig, "pumpconfig", s.PumpConfig, "IAM pump config file.")
	fs.IntVar(&s.PurgeDelay, "purge-delay", s.PurgeDelay, ""+
		"This setting the purge delay (in seconds) when purge the data from Redis to MongoDB or other data stores.")
	fs.StringVar(&s.HealthCheckPath, "health-check-path", s.HealthCheckPath, ""+
		"Specifies liveness health check request path.")
	fs.StringVar(&s.HealthCheckAddress, "health-check-address", s.HealthCheckAddress, ""+
		"Specifies liveness health check bind address.")
	fs.BoolVar(&s.OmitDetailedRecording, "omit-detailed-recording", s.OmitDetailedRecording, ""+
		"Setting this to true will avoid writing policy fields for each authorization request in pumps.")

	return fss
}

func (s *PumpOptions) String() string {
	data, _ := json.Marshal(s)
	return string(data)
}
