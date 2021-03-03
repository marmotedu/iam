// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authorization

import (
	authzv1 "github.com/marmotedu/api/authz/v1"
	"github.com/ory/ladon"

	"github.com/marmotedu/iam/pkg/log"
)

// Authorizer implement the authorize interface that use local repository to
// authorize the subject access review.
type Authorizer struct {
	warden ladon.Warden
}

// NewAuthorizer creates a local repository authorizer and returns it.
func NewAuthorizer(authorizationClient AuthorizationInterface) *Authorizer {
	return &Authorizer{
		warden: &ladon.Ladon{
			Manager:     NewPolicyManager(authorizationClient),
			AuditLogger: NewAuditLogger(authorizationClient),
		},
	}
}

// Authorize to determine the subject access.
func (a *Authorizer) Authorize(request *ladon.Request) *authzv1.Response {
	log.Debug("authorize request", log.Any("request", request))

	if err := a.warden.IsAllowed(request); err != nil {
		return &authzv1.Response{
			Denied: true,
			Reason: err.Error(),
		}
	}

	return &authzv1.Response{
		Allowed: true,
	}
}
