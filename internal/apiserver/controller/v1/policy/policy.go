// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package policy

import (
	srvv1 "github.com/marmotedu/iam/internal/apiserver/service/v1"
	"github.com/marmotedu/iam/internal/apiserver/store"
)

// PolicyHandler create a policy handler used to handle request for policy resource.
type PolicyHandler struct {
	srv   srvv1.Service
	store store.Factory
}

// NewPolicyHandler creates a policy handler.
func NewPolicyHandler(store store.Factory) *PolicyHandler {
	return &PolicyHandler{
		srv:   srvv1.NewService(store),
		store: store,
	}
}
