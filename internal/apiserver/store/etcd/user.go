// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package etcd

import (
	"context"
	"github.com/marmotedu/component-base/pkg/json"

	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

type users struct {
	ds *datastore
}

func newUsers(ds *datastore) *users {
	return &users{ds: ds}
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
	resp, err := u.ds.Get(context.TODO(), username)
	if err != nil {
		return nil, err
	}

	var user v1.User
	if err := json.Unmarshal(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// List return all users.
func (u *users) List(opts metav1.ListOptions) (*v1.UserList, error) {
	return nil, nil
}
