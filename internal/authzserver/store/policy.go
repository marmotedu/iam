// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

import "github.com/ory/ladon"

// PolicyStore defines the policy storage interface.
type PolicyStore interface {
	List() (map[string][]*ladon.DefaultPolicy, error)
}
