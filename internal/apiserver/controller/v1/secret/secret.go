// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package secret

import (
	srvv1 "github.com/marmotedu/iam/internal/apiserver/service/v1"
	"github.com/marmotedu/iam/internal/apiserver/store"
)

// SecretHandler create a secret handler used to handle request for secret resource.
type SecretHandler struct {
	srv   srvv1.Service
	store store.Factory
}

// NewSecretHandler creates a secret handler.
func NewSecretHandler(store store.Factory) *SecretHandler {
	return &SecretHandler{
		srv:   srvv1.NewService(store),
		store: store,
	}
}
