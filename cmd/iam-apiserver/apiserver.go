// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// apiserver is the api server for iam-apiserver service.
// it is responsible for serving the platform RESTful resource management.
package main

import (
	"os"

	"github.com/marmotedu/iam/internal/apiserver"
)

func main() {
	command := apiserver.NewAPIServerCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
