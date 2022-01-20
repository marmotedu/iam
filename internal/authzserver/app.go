// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package authzserver does all of the work necessary to create a authzserver
package authzserver

import (
	"github.com/marmotedu/iam/internal/authzserver/config"
	"github.com/marmotedu/iam/internal/authzserver/options"
	"github.com/marmotedu/iam/pkg/app"
	"github.com/marmotedu/iam/pkg/log"
)

const commandDesc = `Authorization server to run ladon policies which can protecting your resources.
It is written inspired by AWS IAM policiis.

Find more iam-authz-server information at:
    https://github.com/marmotedu/iam/blob/master/docs/guide/en-US/cmd/iam-authz-server.md,

Find more ladon information at:
    https://github.com/ory/ladon`

// NewApp creates an App object with default parameters.
func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("IAM Authorization Server",
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
