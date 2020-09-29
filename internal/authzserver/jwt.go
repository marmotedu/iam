// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authzserver

import (
	pb "github.com/marmotedu/api/proto/apiserver/v1"
	"github.com/marmotedu/iam/internal/authzserver/store"
	"github.com/marmotedu/iam/internal/pkg/middleware"
)

type authzAuth struct {
}

var _ middleware.CacheAuthInterface = &authzAuth{}

func newAuthzServerJwt() middleware.CacheAuthInterface {
	return &authzAuth{}
}

func (a *authzAuth) GetSecret(secretID string) (*pb.SecretInfo, error) {
	return store.GetSecret(secretID)
}
