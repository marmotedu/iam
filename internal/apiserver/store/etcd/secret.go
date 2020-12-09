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

type secrets struct {
	cli             *clientv3.Client
	requestTimeout  time.Duration
	leaseTTLTimeout int
}

func newSecrets(ds *datastore) *secrets {
	return &secrets{cli: ds.cli, requestTimeout: ds.requestTimeout, leaseTTLTimeout: ds.leaseTTLTimeout}
}

// Create creates a new secret account.
func (s *secrets) Create(secret *v1.Secret, opts metav1.CreateOptions) error {
	return nil
}

// Update updates an secret information by the secret identifier.
func (s *secrets) Update(secret *v1.Secret, opts metav1.UpdateOptions) error {
	return nil
}

// Delete deletes the secret by the secret identifier.
func (s *secrets) Delete(username, name string, opts metav1.DeleteOptions) error {
	return nil
}

// DeleteCollection batch deletes the secrets.
func (s *secrets) DeleteCollection(username string, names []string, opts metav1.DeleteOptions) error {
	return nil
}

// Get return an secret by the secret identifier.
func (s *secrets) Get(username, name string, opts metav1.GetOptions) (*v1.Secret, error) {
	return nil, nil
}

// List return all secrets.
func (s *secrets) List(username string, opts metav1.ListOptions) (*v1.SecretList, error) {
	return nil, nil
}
