// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fake

import (
	"fmt"
	"sync"

	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/component-base/pkg/util/idutil"
	"github.com/ory/ladon"

	"github.com/marmotedu/iam/internal/apiserver/store"
)

// ResourceCount defines the number of fake resources.
const ResourceCount = 1000

type datastore struct {
	sync.RWMutex
	users    []*v1.User
	secrets  []*v1.Secret
	policies []*v1.Policy
}

func (ds *datastore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *datastore) Secrets() store.SecretStore {
	return newSecrets(ds)
}

func (ds *datastore) Policies() store.PolicyStore {
	return newPolicies(ds)
}

// NewFakeStore create fake store.
func NewFakeStore() (store.Store, error) {
	// init some user records
	users := make([]*v1.User, 0)
	for i := 1; i <= ResourceCount; i++ {
		users = append(users, &v1.User{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("user%d", i),
				ID:   uint64(i),
			},
			Nickname: fmt.Sprintf("user%d", i),
			Password: fmt.Sprintf("User%d@2020", i),
			Email:    fmt.Sprintf("user%d@qq.com", i),
		})
	}

	// init some secrets records
	secrets := make([]*v1.Secret, 0)
	for i := 1; i <= ResourceCount; i++ {
		secrets = append(secrets, &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("secret%d", i),
				ID:   uint64(i),
			},
			Username:  fmt.Sprintf("user%d", i),
			SecretID:  idutil.NewSecretID(),
			SecretKey: idutil.NewSecretKey(),
		})
	}

	// init some policy records
	policies := make([]*v1.Policy, 0)
	for i := 1; i <= ResourceCount; i++ {
		policies = append(policies, &v1.Policy{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("policy%d", i),
				ID:   uint64(i),
			},
			Username: fmt.Sprintf("user%d", i),
			Policy:   ladon.DefaultPolicy{},
		})
	}

	return &datastore{
		users:    users,
		secrets:  secrets,
		policies: policies,
	}, nil
}
