// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package options contains flags and options for initializing an apiserver
package options

import (
	"encoding/json"

	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
	"github.com/marmotedu/log"

	"github.com/marmotedu/iam/internal/authzserver/analytics"
	genericoptions "github.com/marmotedu/iam/internal/pkg/options"
	"github.com/marmotedu/iam/internal/pkg/server"
)

// ServerRunOptions runs a authzserver.
type ServerRunOptions struct {
	AuthzConfig             string                                 `json:"authzconfig" mapstructure:"-"`
	RPCServer               string                                 `json:"rpcserver" mapstructure:"rpcserver"`
	ClientCA                string                                 `json:"client-ca-file" mapstructure:"client-ca-file"`
	GenericServerRunOptions *genericoptions.ServerRunOptions       `json:"server" mapstructure:"server"`
	InsecureServing         *genericoptions.InsecureServingOptions `json:"insecure" mapstructure:"insecure"`
	SecureServing           *genericoptions.SecureServingOptions   `json:"secure" mapstructure:"secure"`
	RedisOptions            *genericoptions.RedisOptions           `json:"redis" mapstructure:"redis"`
	FeatureOptions          *genericoptions.FeatureOptions         `json:"feature" mapstructure:"feature"`
	Log                     *log.Options                           `json:"log" mapstructure:"log"`
	AnalyticsOptions        *analytics.AnalyticsOptions            `json:"analytics" mapstructure:"analytics"`
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters.
func NewServerRunOptions() *ServerRunOptions {
	s := ServerRunOptions{
		RPCServer:               "127.0.0.1:8081",
		ClientCA:                "",
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		InsecureServing:         genericoptions.NewInsecureServingOptions(),
		SecureServing:           genericoptions.NewSecureServingOptions(),
		RedisOptions:            genericoptions.NewRedisOptions(),
		FeatureOptions:          genericoptions.NewFeatureOptions(),
		Log:                     log.NewOptions(),
		AnalyticsOptions:        analytics.NewAnalyticsOptions(),
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
	s.AnalyticsOptions.AddFlags(fss.FlagSet("analytics"))
	s.RedisOptions.AddFlags(fss.FlagSet("redis"))
	s.FeatureOptions.AddFlags(fss.FlagSet("features"))
	s.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	s.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	s.Log.AddFlags(fss.FlagSet("logs"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.StringVar(&s.AuthzConfig, "authzconfig", s.AuthzConfig, "IAM AuthzServer config file.")
	fs.StringVar(&s.RPCServer, "rpcserver", s.RPCServer, "The address of iam rpc server. "+
		"The rpc server can provide all the secrets and policies to use.")
	fs.StringVar(&s.ClientCA, "client-ca-file", s.ClientCA, ""+
		"If set, any request presenting a client certificate signed by one of "+
		"the authorities in the client-ca-file is authenticated with an identity "+
		"corresponding to the CommonName of the client certificate.")

	return fss
}

func (s *ServerRunOptions) String() string {
	data, _ := json.Marshal(s)
	return string(data)
}
