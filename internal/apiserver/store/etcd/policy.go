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
	"github.com/marmotedu/errors"
)

type policies struct {
	ds *datastore
}

func newPolicies(ds *datastore) *policies {
	return &policies{ds: ds}
}

var keyPolicy = "/policies/%v/%v"

func (p *policies) getKey(username string, name string) string {
	return fmt.Sprintf(keyPolicy, username, name)
}

// Create creates a new policy.
func (p *policies) Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) error {
	return p.ds.Put(ctx, p.getKey(policy.Username, policy.Name), jsonutil.ToString(policy))
}

// Update updates an policy information.
func (p *policies) Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) error {
	return p.ds.Put(ctx, p.getKey(policy.Username, policy.Name), jsonutil.ToString(policy))
}

// Delete deletes the policy by the policy identifier.
func (p *policies) Delete(ctx context.Context, username, name string, opts metav1.DeleteOptions) error {
	if _, err := p.ds.Delete(ctx, p.getKey(username, name)); err != nil {
		return err
	}

	return nil
}

// DeleteByUser deletes policies by username.
func (p *policies) DeleteByUser(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	if _, err := p.ds.Delete(ctx, p.getKey(username, "")); err != nil {
		return err
	}

	return nil
}

// DeleteCollection batch deletes the policies.
func (p *policies) DeleteCollection(
	ctx context.Context,
	username string,
	names []string,
	opts metav1.DeleteOptions,
) error {
	return nil
}

// DeleteCollectionByUser batch deletes policies usernames.
func (p *policies) DeleteCollectionByUser(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	return nil
}

// Get return an policy by the policy identifier.
func (p *policies) Get(ctx context.Context, username, name string, opts metav1.GetOptions) (*v1.Policy, error) {
	resp, err := p.ds.Get(ctx, p.getKey(username, name))
	if err != nil {
		return nil, err
	}

	var policy v1.Policy
	if err := json.Unmarshal(resp, &policy); err != nil {
		return nil, errors.Wrap(err, "unmarshal to Policy struct failed")
	}

	return &policy, nil
}

// List return all policies.
func (p *policies) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.PolicyList, error) {
	kvs, err := p.ds.List(ctx, p.getKey(username, ""))
	if err != nil {
		return nil, err
	}

	ret := &v1.PolicyList{
		ListMeta: metav1.ListMeta{
			TotalCount: int64(len(kvs)),
		},
	}

	for _, v := range kvs {
		var policy v1.Policy
		if err := json.Unmarshal(v.Value, &policy); err != nil {
			return nil, errors.Wrap(err, "unmarshal to Policy struct failed")
		}

		ret.Items = append(ret.Items, &policy)
	}

	return ret, nil
}
