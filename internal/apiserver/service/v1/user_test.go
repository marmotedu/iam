// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/AlekSi/pointer"
	gomock "github.com/golang/mock/gomock"
	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"

	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/apiserver/store/fake"
)

func TestMain(m *testing.M) {
	_, _ = fake.GetFakeFactoryOr()
	os.Exit(m.Run())
}

func BenchmarkListUser(b *testing.B) {
	opts := metav1.ListOptions{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(50),
	}
	storeIns, _ := fake.GetFakeFactoryOr()
	u := &userService{
		store: storeIns,
	}

	for i := 0; i < b.N; i++ {
		// _, _ = ListUserBadPerformance(opts)
		_, _ = u.List(context.TODO(), opts)
	}
}

func Test_newUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)

	type args struct {
		srv *service
	}
	tests := []struct {
		name string
		args args
		want *userService
	}{
		{
			name: "default",
			args: args{
				srv: &service{store: mockFactory},
			},
			want: &userService{
				store: mockFactory,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUsers(tt.args.srv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_List(t *testing.T) {
	storeIns, _ := fake.GetFakeFactoryOr()
	var limit int64 = 3

	opts := metav1.ListOptions{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(limit),
	}
	items := make([]*v1.UserV2, 0, limit)
	for _, user := range fake.FakeUsers(3) {
		userv2 := &v1.UserV2{
			User: &v1.User{
				ObjectMeta: metav1.ObjectMeta{
					ID:        user.ID,
					Name:      user.Name,
					CreatedAt: user.CreatedAt,
					UpdatedAt: user.UpdatedAt,
				},
				Nickname: user.Nickname,
				Email:    user.Email,
				Phone:    user.Phone,
			},
			TotalPolicy: fake.ResourceCount,
		}
		items = append(items, userv2)
	}
	wantUserList := &v1.UserListV2{
		ListMeta: metav1.ListMeta{
			TotalCount: fake.ResourceCount,
		},
		Items: items,
	}

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx  context.Context
		opts metav1.ListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.UserListV2
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: storeIns,
			},
			args: args{
				ctx:  context.TODO(),
				opts: opts,
			},
			want:    wantUserList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userService{
				store: tt.fields.store,
			}
			got, err := u.List(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
