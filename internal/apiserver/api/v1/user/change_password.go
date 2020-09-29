// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

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

// ChangePasswordRequest defines the ChangePasswordRequest data format.
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"omitempty"`
	NewPassword string `json:"newPassword" binding:"password"`
}

// ChangePassword change the user's password by the user identifier.
func ChangePassword(c *gin.Context) {
	log.Info("change password function called.", log.String("X-Request-Id", requestid.Get(c)))

	var r ChangePasswordRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	u, err := store.Client().Users().Get(c.Param("name"), metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}

	if err := u.Compare(r.OldPassword); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrPasswordIncorrect, err.Error()), nil)
		return
	}

	u.Password = r.NewPassword

	// Save changed fields.
	if err := store.Client().Users().Update(u, metav1.UpdateOptions{}); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrDatabase, err.Error()), nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}
