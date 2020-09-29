// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

var client Store

// Store defines the iam platform storage interface.
type Store interface {
	Users() UserStore
	Secrets() SecretStore
	Policies() PolicyStore
}

// Client return the store client instance.
func Client() Store {
	return client
}

// SetClient set the iam store client.
func SetClient(store Store) {
	client = store
}
