// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Usage:
//	go run main.go
//	curl http://127.0.0.1:7070/user/foo

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	"github.com/marmotedu/iam/internal/pkg/code"
)

func main() {
	r := gin.Default()

	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		if err := getUser(name); err != nil {
			core.WriteResponse(c, err, nil)
			return
		}

		core.WriteResponse(c, nil, map[string]string{"email": name + "@foxmail.com"})
	})

	r.Run(":7070")
}

func getUser(name string) error {
	if err := queryDatabase(name); err != nil {
		return errors.Wrap(err, "get user failed.")
	}

	return nil
}

func queryDatabase(name string) error {
	return errors.WithCode(code.ErrDatabase, "user '%s' not found.", name)
}
