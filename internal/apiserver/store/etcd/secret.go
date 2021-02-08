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

type secrets struct {
	ds *datastore
}

func newSecrets(ds *datastore) *secrets {
	return &secrets{ds: ds}
}

var keySecret = "/secrets/%v/%v"

func (u *secrets) getKey(username string, secretID string) string {
	return fmt.Sprintf(keySecret, username, secretID)
}

// Create creates a new secret.
func (u *secrets) Create(ctx context.Context, secret *v1.Secret, opts metav1.CreateOptions) error {
	if err := u.ds.Put(ctx, u.getKey(secret.Username, secret.SecretID), jsonutil.ToString(secret)); err != nil {
		return err
	}

	return nil
}

// Update updates an secret information.
func (u *secrets) Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOptions) error {
	if err := u.ds.Put(ctx, u.getKey(secret.Username, secret.SecretID), jsonutil.ToString(secret)); err != nil {
		return err
	}

	return nil
}

// Delete deletes the secret by the secret identifier.
func (u *secrets) Delete(ctx context.Context, username, secretID string, opts metav1.DeleteOptions) error {
	if _, err := u.ds.Delete(ctx, u.getKey(username, secretID)); err != nil {
		return err
	}

	return nil
}

// DeleteCollection batch deletes the secrets.
func (u *secrets) DeleteCollection(ctx context.Context, username string, secretIDs []string, opts metav1.DeleteOptions) error {
	return nil
}

// Get return an secret by the secret identifier.
func (u *secrets) Get(ctx context.Context, username, secretID string, opts metav1.GetOptions) (*v1.Secret, error) {
	resp, err := u.ds.Get(ctx, u.getKey(username, secretID))
	if err != nil {
		return nil, err
	}

	var secret v1.Secret
	if err := json.Unmarshal(resp, &secret); err != nil {
		return nil, err
	}
	return &secret, nil
}

// List return all secrets.
func (u *secrets) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error) {
	kvs, err := u.ds.List(ctx, u.getKey(username))
	if err != nil {
		return nil, err
	}

	ret := &v1.SecretList{
		TotalCount: len(kvs),
	}

	for k, v := range kvs {
		var secret v1.Secret
		if err := json.Unmarshal(v.Value, &secret); err != nil {
			return err
		}

		ret.Items = append(ret.Items, &secret)
	}

	return ret, nil
}
