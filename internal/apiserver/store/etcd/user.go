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

type users struct {
	cli             *clientv3.Client
	requestTimeout  time.Duration
	leaseTTLTimeout int
}

func newUsers(ds *datastore) *users {
	return &users{cli: ds.cli, requestTimeout: ds.requestTimeout, leaseTTLTimeout: ds.leaseTTLTimeout}
}

// Create creates a new user account.
func (u *users) Create(user *v1.User, opts metav1.CreateOptions) error {
	return nil
}

// Update updates an user account information.
func (u *users) Update(user *v1.User, opts metav1.UpdateOptions) error {
	return nil
}

// Delete deletes the user by the user identifier.
func (u *users) Delete(username string, opts metav1.DeleteOptions) error {
	return nil
}

// DeleteCollection batch deletes the users.
func (u *users) DeleteCollection(usernames []string, opts metav1.DeleteOptions) error {
	return nil
}

// Get return an user by the user identifier.
func (u *users) Get(username string, opts metav1.GetOptions) (*v1.User, error) {
	/*
		ctx, cancel := context.WithTimeout(context.Background(), u.requestTimeout)
		defer cancel()
		key := fmt.Sprintf("v2/%s", user.Name)
		rsp, err := f.cli.Get(ctx, key, clientv3.WithPrefix())
		if err != nil {
			return nil, fmt.Errorf("failed get function [key:%v] from etcd:%v ", key, err)
		}
		if len(resp.Kvs) == 0 {
			return nil, fmt.Errorf("failed to get valid function value[key:%v] from etcd:%v", key, err)
		}
		return resp.Kvs[0].Value, nil
	*/
	return nil, nil
}

// List return all users.
func (u *users) List(opts metav1.ListOptions) (*v1.UserList, error) {
	return nil, nil
}
