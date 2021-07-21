// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

import pb "github.com/marmotedu/api/proto/apiserver/v1"

// SecretStore defines the secret storage interface.
type SecretStore interface {
	// List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error)
	List() (map[string]*pb.SecretInfo, error)
}
