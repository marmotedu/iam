// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"sync"

	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/pkg/code"
)

// ListUser returns user list in the storage. This function has a good performance.
func ListUser(opts metav1.ListOptions) (*v1.UserListV2, error) {
	users, err := store.Client().Users().List(opts)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	var m sync.Map

	// Improve query efficiency in parallel
	for _, u := range users.Items {
		wg.Add(1)

		go func(u *v1.User) {
			defer wg.Done()

			policies, err := store.Client().Policies().List(u.Name, metav1.ListOptions{})
			if err != nil {
				errChan <- errors.WithCode(code.ErrDatabase, err.Error())
				return
			}

			m.Store(u.ID, &v1.UserV2{
				User: &v1.User{
					ObjectMeta: metav1.ObjectMeta{
						ID:        u.ID,
						Name:      u.Name,
						CreatedAt: u.CreatedAt,
						UpdatedAt: u.UpdatedAt,
					},
					Nickname: u.Nickname,
					Email:    u.Email,
					Phone:    u.Phone,
				},
				TotalPolicy: policies.TotalCount,
			})
		}(u)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		return nil, err
	}

	//infos := make([]*v1.UserV2, 0)
	infos := make([]*v1.UserV2, 0, len(users.Items))
	for _, user := range users.Items {
		info, _ := m.Load(user.ID)
		infos = append(infos, info.(*v1.UserV2))
	}

	return &v1.UserListV2{ListMeta: users.ListMeta, Items: infos}, nil
}

// ListUserBadPerformance returns user list in the storage. This function has a bad performance.
func ListUserBadPerformance(opts metav1.ListOptions) (*v1.UserListV2, error) {
	users, err := store.Client().Users().List(opts)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	infos := make([]*v1.UserV2, 0)
	for _, u := range users.Items {
		policies, err := store.Client().Policies().List(u.Name, metav1.ListOptions{})
		if err != nil {
			return nil, errors.WithCode(code.ErrDatabase, err.Error())
		}

		infos = append(infos, &v1.UserV2{
			User: &v1.User{
				ObjectMeta: metav1.ObjectMeta{
					ID:        u.ID,
					Name:      u.Name,
					CreatedAt: u.CreatedAt,
					UpdatedAt: u.UpdatedAt,
				},
				Nickname: u.Nickname,
				Email:    u.Email,
				Phone:    u.Phone,
			},
			TotalPolicy: policies.TotalCount,
		})
	}

	return &v1.UserListV2{ListMeta: users.ListMeta, Items: infos}, nil
}
