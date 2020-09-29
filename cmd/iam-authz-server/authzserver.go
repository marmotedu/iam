// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// authzserver is the server for iam-authz-server.
// It is responsible for serving the ladon authorization request.
package main

import (
	"os"

	"github.com/marmotedu/iam/internal/authzserver"
)

func main() {
	command := authzserver.NewAuthzServerCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
