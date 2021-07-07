// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package auth

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	"github.com/marmotedu/iam/internal/pkg/code"
	"github.com/marmotedu/iam/internal/pkg/middleware"
)

// BasicStrategy defines Basic authentication strategy.
type BasicStrategy struct {
	compare func(username string, password string) bool
}

var _ middleware.AuthStrategy = &BasicStrategy{}

// NewBasicStrategy create basic strategy with compare function.
func NewBasicStrategy(compare func(username string, password string) bool) BasicStrategy {
	return BasicStrategy{
		compare: compare,
	}
}

// AuthFunc defines basic strategy as the gin authentication middleware.
func (b BasicStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			core.WriteResponse(
				c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()

			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !b.compare(pair[0], pair[1]) {
			core.WriteResponse(
				c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()

			return
		}

		c.Set(middleware.UsernameKey, pair[0])

		c.Next()
	}
}
