// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package secret

import (
	"github.com/gin-gonic/gin"

	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"

	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/pkg/code"
	"github.com/marmotedu/iam/pkg/log"
)

// Delete delete a secret by the secret identifier.
func Delete(c *gin.Context) {
	log.L(c).Info("delete secret function called.")
	if err := store.Client().Secrets().Delete(c.GetHeader("username"), c.Param("name"), metav1.DeleteOptions{Unscoped: true}); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}
