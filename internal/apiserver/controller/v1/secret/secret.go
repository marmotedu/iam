// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package secret

import (
	srvv1 "github.com/marmotedu/iam/internal/apiserver/service/v1"
	"github.com/marmotedu/iam/internal/apiserver/store"
)

// SecretController create a secret handler used to handle request for secret resource.
type SecretController struct {
	srv srvv1.Service
}

// NewSecretController creates a secret handler.
func NewSecretController(store store.Factory) *SecretController {
	return &SecretController{
		srv: srvv1.NewService(store),
	}
}
