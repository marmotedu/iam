// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package options contains flags and options for initializing an apiserver
package options

import (
	"encoding/json"

	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
	"github.com/marmotedu/log"

	genericoptions "github.com/marmotedu/iam/internal/pkg/options"
	"github.com/marmotedu/iam/internal/pkg/server"
)

// ServerRunOptions runs a iam api server.
type ServerRunOptions struct {
	APIConfig               string                                 `json:"apiconfig" mapstructure:"-"`
	GenericServerRunOptions *genericoptions.ServerRunOptions       `json:"server" mapstructure:"server"`
	GrpcOptions             *genericoptions.GrpcOptions            `json:"grpc" mapstructure:"grpc"`
	InsecureServing         *genericoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing           *genericoptions.SecureServingOptions   `json:"secure" mapstructure:"secure"`
	MySQLOptions            *genericoptions.MySQLOptions           `json:"mysql" mapstructure:"mysql"`
	RedisOptions            *genericoptions.RedisOptions           `json:"redis" mapstructure:"redis"`
	JwtOptions              *genericoptions.JwtOptions             `json:"jwt" mapstructure:"jwt"`
	Log                     *log.Options                           `json:"log" mapstructure:"log"`
	FeatureOptions          *genericoptions.FeatureOptions         `json:"feature" mapstructure:"feature"`
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters.
func NewServerRunOptions() *ServerRunOptions {
	s := ServerRunOptions{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		GrpcOptions:             genericoptions.NewGrpcOptions(),
		InsecureServing:         genericoptions.NewInsecureServingOptions(),
		SecureServing:           genericoptions.NewSecureServingOptions(),
		MySQLOptions:            genericoptions.NewMySQLOptions(),
		RedisOptions:            genericoptions.NewRedisOptions(),
		JwtOptions:              genericoptions.NewJwtOptions(),
		Log:                     log.NewOptions(),
		FeatureOptions:          genericoptions.NewFeatureOptions(),
	}

	return &s
}

// ApplyTo applies the run options to the method receiver and returns self.
func (s *ServerRunOptions) ApplyTo(c *server.Config) error {
	return nil
}

// Flags returns flags for a specific APIServer by section name.
func (s *ServerRunOptions) Flags() (fss cliflag.NamedFlagSets) {
	s.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	s.JwtOptions.AddFlags(fss.FlagSet("jwt"))
	s.GrpcOptions.AddFlags(fss.FlagSet("grpc"))
	s.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	s.RedisOptions.AddFlags(fss.FlagSet("redis"))
	s.FeatureOptions.AddFlags(fss.FlagSet("features"))
	s.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	s.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	s.Log.AddFlags(fss.FlagSet("logs"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.StringVar(&s.APIConfig, "apiconfig", s.APIConfig, "IAM APIServer config file.")

	return fss
}

func (s *ServerRunOptions) String() string {
	data, _ := json.Marshal(s)
	return string(data)
}
