// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const headerXRequestID = "X-Request-ID"

// RequestID add X-Request-ID to response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		rid := c.GetHeader(headerXRequestID)

		if rid == "" {
			rid = uuid.Must(uuid.NewV4()).String()
			c.Header(headerXRequestID, rid)
		}

		// Set headerXRequestID header
		c.Writer.Header().Set(headerXRequestID, rid)
		c.Next()
	}
}
