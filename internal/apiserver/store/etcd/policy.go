// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package etcd

import (
	"context"
	"fmt"

	v1 "github.com/marmotedu/api/apiserver/v1"
	"github.com/marmotedu/component-base/pkg/json"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/component-base/pkg/util/jsonutil"
)

type policies struct {
	ds *datastore
}

func newPolicies(ds *datastore) *policies {
	return &policies{ds: ds}
}

var keyPolicy = "/policies/%v/%v"

func (u *policies) getKey(username string, policyID string) string {
	return fmt.Sprintf(keyPolicy, username, policyID)
}

// Create creates a new policy.
func (u *policies) Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) error {
	if err := u.ds.Put(ctx, u.getKey(policy.Username, policy.PolicyID), jsonutil.ToString(policy)); err != nil {
		return err
	}

	return nil
}

// Update updates an policy information.
func (u *policies) Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) error {
	if err := u.ds.Put(ctx, u.getKey(policy.Username, policy.PolicyID), jsonutil.ToString(policy)); err != nil {
		return err
	}

	return nil
}

// Delete deletes the policy by the policy identifier.
func (u *policies) Delete(ctx context.Context, username, policyID string, opts metav1.DeleteOptions) error {
	if _, err := u.ds.Delete(ctx, u.getKey(username, policyID)); err != nil {
		return err
	}

	return nil
}

// DeleteByUser deletes policies by username.
func (p *policies) DeleteByUser(username string, opts metav1.DeleteOptions) error {
	if _, err := u.ds.Delete(ctx, u.getKey(username, "")); err != nil {
		return err
	}

	return nil
}

// DeleteCollection batch deletes the policies.
func (u *policies) DeleteCollection(ctx context.Context, username string, policyIDs []string, opts metav1.DeleteOptions) error {
	return nil
}

// DeleteCollectionByUser batch deletes policies usernames.
func (p *policies) DeleteCollectionByUser(usernames []string, opts metav1.DeleteOptions) error {
	return nil
}

// Get return an policy by the policy identifier.
func (u *policies) Get(ctx context.Context, username, policyID string, opts metav1.GetOptions) (*v1.Policy, error) {
	resp, err := u.ds.Get(ctx, u.getKey(username, policyID))
	if err != nil {
		return nil, err
	}

	var policy v1.Policy
	if err := json.Unmarshal(resp, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}

// List return all policies.
func (u *policies) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error) {
	kvs, err := u.ds.List(ctx, u.getKey(username))
	if err != nil {
		return nil, err
	}

	ret := &v1.PolicyList{
		TotalCount: len(kvs),
	}

	for k, v := range kvs {
		var policy v1.Policy
		if err := json.Unmarshal(v.Value, &policy); err != nil {
			return err
		}

		ret.Items = append(ret.Items, &policy)
	}

	return ret, nil
}
