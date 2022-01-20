// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package apiserver does all the work necessary to create a iam APIServer.
package apiserver

import (
	"github.com/marmotedu/iam/internal/apiserver/config"
	"github.com/marmotedu/iam/internal/apiserver/options"
	"github.com/marmotedu/iam/pkg/app"
	"github.com/marmotedu/iam/pkg/log"
)

const commandDesc = `The IAM API server validates and configures data
for the api objects which include users, policies, secrets, and
others. The API Server services REST operations to do the api objects management.

Find more iam-apiserver information at:
    https://github.com/marmotedu/iam/blob/master/docs/guide/en-US/cmd/iam-apiserver.md`

// NewApp creates an App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("IAM API Server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
		app.WithRunFunc(run(opts)),
	)

	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		log.Init(opts.Log)
		defer log.Flush()

		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}

		return Run(cfg)
	}
}
