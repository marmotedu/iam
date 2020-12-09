// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package etcd

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

type policies struct {
	cli             *clientv3.Client
	requestTimeout  time.Duration
	leaseTTLTimeout int
}

func newPolicies(ds *datastore) *policies {
	return &policies{cli: ds.cli, requestTimeout: ds.requestTimeout, leaseTTLTimeout: ds.leaseTTLTimeout}
}

// Create creates a new ladon policy.
func (p *policies) Create(policy *v1.Policy, opts metav1.CreateOptions) error {
	return nil
}

// Update updates policy by the policy identifier.
func (p *policies) Update(policy *v1.Policy, opts metav1.UpdateOptions) error {
	return nil
}

// Delete deletes the policy by the policy identifier.
func (p *policies) Delete(username, name string, opts metav1.DeleteOptions) error {
	return nil
}

// DeleteCollection batch deletes policies by policies ids.
func (p *policies) DeleteCollection(username string, names []string, opts metav1.DeleteOptions) error {
	return nil
}

// DeleteByUser deletes policies by username.
func (p *policies) DeleteByUser(username string, opts metav1.DeleteOptions) error {
	return nil
}

// DeleteCollectionByUser batch deletes policies usernames.
func (p *policies) DeleteCollectionByUser(usernames []string, opts metav1.DeleteOptions) error {
	return nil
}

// Get return policy by the policy identifier.
func (p *policies) Get(username, name string, opts metav1.GetOptions) (*v1.Policy, error) {
	return nil, nil
}

// List return all policies.
func (p *policies) List(username string, opts metav1.ListOptions) (*v1.PolicyList, error) {
	return nil, nil
}
