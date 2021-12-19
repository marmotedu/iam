// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package authorizer defines authorization interface.
package authorizer

import (
	"fmt"
	"strings"
	"time"

	"github.com/marmotedu/component-base/pkg/json"
	"github.com/ory/ladon"

	"github.com/marmotedu/iam/internal/authzserver/analytics"
	"github.com/marmotedu/iam/internal/authzserver/authorization"
)

// PolicyGetter defines function to get policy for a given user.
type PolicyGetter interface {
	GetPolicy(key string) ([]*ladon.DefaultPolicy, error)
}

// Authorization implements authorization.AuthorizationInterface interface.
type Authorization struct {
	getter PolicyGetter
}

// NewAuthorization create a new Authorization instance.
func NewAuthorization(getter PolicyGetter) authorization.AuthorizationInterface {
	return &Authorization{getter}
}

// Create create a policy.
// Return nil because we use mysql storage to store the policy.
func (auth *Authorization) Create(policy *ladon.DefaultPolicy) error {
	return nil
}

// Update update a policy.
// Return nil because we use mysql storage to store the policy.
func (auth *Authorization) Update(policy *ladon.DefaultPolicy) error {
	return nil
}

// Delete delete a policy by the given identifier.
// Return nil because we use mysql storage to store the policy.
func (auth *Authorization) Delete(id string) error {
	return nil
}

// DeleteCollection batch delete policies by the given identifiers.
// Return nil because we use mysql storage to store the policy.
func (auth *Authorization) DeleteCollection(idList []string) error {
	return nil
}

// Get returns the policy detail by the given identifier.
// Return nil because we use mysql storage to store the policy.
func (auth *Authorization) Get(id string) (*ladon.DefaultPolicy, error) {
	return &ladon.DefaultPolicy{}, nil
}

// List returns all the policies under the username.
func (auth *Authorization) List(username string) ([]*ladon.DefaultPolicy, error) {
	return auth.getter.GetPolicy(username)
}

// LogRejectedAccessRequest write rejected subject access to redis.
func (auth *Authorization) LogRejectedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	var conclusion string
	if len(d) > 1 {
		allowed := joinPoliciesNames(d[0 : len(d)-1])
		denied := d[len(d)-1].GetID()
		conclusion = fmt.Sprintf("policies %s allow access, but policy %s forcefully denied it", allowed, denied)
	} else if len(d) == 1 {
		denied := d[len(d)-1].GetID()
		conclusion = fmt.Sprintf("policy %s forcefully denied the access", denied)
	} else {
		conclusion = "no policy allowed access"
	}

	rstring, pstring, dstring := convertToString(r, p, d)
	record := analytics.AnalyticsRecord{
		TimeStamp:  time.Now().Unix(),
		Username:   r.Context["username"].(string),
		Effect:     ladon.DenyAccess,
		Conclusion: conclusion,
		Request:    rstring,
		Policies:   pstring,
		Deciders:   dstring,
	}

	record.SetExpiry(0)
	_ = analytics.GetAnalytics().RecordHit(&record)
}

// LogGrantedAccessRequest write granted subject access to redis.
func (auth *Authorization) LogGrantedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	conclusion := fmt.Sprintf("policies %s allow access", joinPoliciesNames(d))
	rstring, pstring, dstring := convertToString(r, p, d)
	record := analytics.AnalyticsRecord{
		TimeStamp:  time.Now().Unix(),
		Username:   r.Context["username"].(string),
		Effect:     ladon.AllowAccess,
		Conclusion: conclusion,
		Request:    rstring,
		Policies:   pstring,
		Deciders:   dstring,
	}

	record.SetExpiry(0)
	_ = analytics.GetAnalytics().RecordHit(&record)
}

func joinPoliciesNames(policies ladon.Policies) string {
	names := []string{}
	for _, policy := range policies {
		names = append(names, policy.GetID())
	}

	return strings.Join(names, ", ")
}

func convertToString(r *ladon.Request, p ladon.Policies, d ladon.Policies) (string, string, string) {
	rbytes, _ := json.Marshal(r)
	pbytes, _ := json.Marshal(p)
	dbytes, _ := json.Marshal(d)

	return string(rbytes), string(pbytes), string(dbytes)
}
