// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// iamctl is the command line tool for iam platform.
package main

import (
	"os"

	"github.com/marmotedu/iam/internal/iamctl/cmd"
)

func main() {
	command := cmd.NewDefaultIAMCtlCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
