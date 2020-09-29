// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

import (
	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

// PolicyStore defines the policy storage interface.
type PolicyStore interface {
	Create(policy *v1.Policy, opts metav1.CreateOptions) error
	Update(policy *v1.Policy, opts metav1.UpdateOptions) error
	Delete(username string, policyID string, opts metav1.DeleteOptions) error
	DeleteCollection(username string, policyIDs []string, opts metav1.DeleteOptions) error
	DeleteByUser(username string, opts metav1.DeleteOptions) error
	DeleteCollectionByUser(usernames []string, opts metav1.DeleteOptions) error
	Get(username string, policyID string, opts metav1.GetOptions) (*v1.Policy, error)
	List(username string, opts metav1.ListOptions) (*v1.PolicyList, error)
}
