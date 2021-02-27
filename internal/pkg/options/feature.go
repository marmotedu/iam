// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"github.com/spf13/pflag"

	"github.com/marmotedu/iam/internal/pkg/server"
)

// FeatureOptions contains configuration items related to API server features.
type FeatureOptions struct {
	EnableProfiling bool `json:"profiling"      mapstructure:"profiling"`
	EnableMetrics   bool `json:"enable-metrics" mapstructure:"enable-metrics"`
}

// NewFeatureOptions creates a FeatureOptions object with default parameters.
func NewFeatureOptions() *FeatureOptions {
	defaults := server.NewConfig()

	return &FeatureOptions{
		EnableMetrics:   defaults.EnableMetrics,
		EnableProfiling: defaults.EnableProfiling,
	}
}

// ApplyTo applies the run options to the method receiver and returns self.
func (o *FeatureOptions) ApplyTo(c *server.Config) error {
	c.EnableProfiling = o.EnableProfiling
	c.EnableMetrics = o.EnableMetrics

	return nil
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (o *FeatureOptions) Validate() []error {
	return []error{}
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (o *FeatureOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.BoolVar(&o.EnableProfiling, "feature.profiling", o.EnableProfiling,
		"Enable profiling via web interface host:port/debug/pprof/")

	fs.BoolVar(&o.EnableMetrics, "feature.enable-metrics", o.EnableMetrics,
		"Enables metrics on the apiserver at /metrics")
}
