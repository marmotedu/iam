// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package secret

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/pkg/code"
	"github.com/marmotedu/log"
)

// DeleteCollection delete secrets by secret names.
func DeleteCollection(c *gin.Context) {
	log.Info("batch delete policy function called.", log.String("X-Request-Id", requestid.Get(c)))

	if err := store.Client().Policies().DeleteCollection(
		c.GetHeader("username"),
		c.QueryArray("name"),
		metav1.DeleteOptions{},
	); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}
