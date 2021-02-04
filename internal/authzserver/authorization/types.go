// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authorization

//go:generate mockgen -destination mock_authorization.go -package authorization github.com/marmotedu/iam/internal/authzserver/authorization AuthorizationInterface

import (
	"github.com/ory/ladon"
)

// AuthorizationInterface defiens the CURD method for lady policy.
type AuthorizationInterface interface {
	Create(*ladon.DefaultPolicy) error
	Update(*ladon.DefaultPolicy) error
	Delete(id string) error
	DeleteCollection(idList []string) error
	Get(id string) (*ladon.DefaultPolicy, error)
	List(username string) ([]*ladon.DefaultPolicy, error)

	// The following two functions tracks denied and granted authorizations.
	LogRejectedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies)
	LogGrantedAccessRequest(request *ladon.Request, pool ladon.Policies, deciders ladon.Policies)
}
