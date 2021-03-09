// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// authzserver is the server for iam-authz-server.
// It is responsible for serving the ladon authorization request.
package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/marmotedu/iam/internal/authzserver"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	authzserver.NewApp("iam-authz-server").Run()
}
