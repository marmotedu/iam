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
		//_, _ = u.ListWithBadPerformance(context.TODO(), opts)
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

func (s *Suite) Test_userService_List() {
	storeIns, _ := fake.GetFakeFactoryOr()
	var limit int64 = 3

	opts := metav1.ListOptions{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(limit),
	}
	items := make([]*v1.User, 0, limit)
	for _, u := range fake.FakeUsers(3) {
		user := &v1.User{
			ObjectMeta: metav1.ObjectMeta{
				ID:        u.ID,
				Name:      u.Name,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			},
			Nickname:    u.Nickname,
			Email:       u.Email,
			Phone:       u.Phone,
			TotalPolicy: fake.ResourceCount,
		}
		items = append(items, user)
	}
	wantUserList := &v1.UserList{
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
		want    *v1.UserList
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
		s.T().Run(tt.name, func(t *testing.T) {
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

func (s *Suite) Test_userService_ListWithBadPerformance() {
	storeIns, _ := fake.GetFakeFactoryOr()
	var limit int64 = 3

	opts := metav1.ListOptions{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(limit),
	}
	items := make([]*v1.User, 0, limit)
	for _, u := range fake.FakeUsers(3) {
		user := &v1.User{
			ObjectMeta: metav1.ObjectMeta{
				ID:        u.ID,
				Name:      u.Name,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			},
			Nickname:    u.Nickname,
			Email:       u.Email,
			Phone:       u.Phone,
			TotalPolicy: fake.ResourceCount,
		}
		items = append(items, user)
	}
	wantUserList := &v1.UserList{
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
		want    *v1.UserList
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
		s.T().Run(tt.name, func(t *testing.T) {
			u := &userService{
				store: tt.fields.store,
			}
			got, err := u.ListWithBadPerformance(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUserBadPerformance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListUserBadPerformance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_userService_Create() {
	s.mockUserStore.EXPECT().Create(gomock.Any(), gomock.Eq(s.users[0]), gomock.Any()).Return(nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx  context.Context
		user *v1.User
		opts metav1.CreateOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:  context.TODO(),
				user: s.users[0],
				opts: metav1.CreateOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			u := &userService{
				store: tt.fields.store,
			}
			if err := u.Create(tt.args.ctx, tt.args.user, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("userService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_userService_DeleteCollection() {
	s.mockUserStore.EXPECT().DeleteCollection(gomock.Any(), []string{"colin", "john"}, gomock.Any()).Return(nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx       context.Context
		usernames []string
		opts      metav1.DeleteOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:       context.TODO(),
				usernames: []string{"colin", "john"},
				opts:      metav1.DeleteOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			u := &userService{
				store: tt.fields.store,
			}
			if err := u.DeleteCollection(tt.args.ctx, tt.args.usernames, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("userService.DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_userService_Delete() {
	s.mockUserStore.EXPECT().Delete(gomock.Any(), "colin", gomock.Any()).Return(nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx      context.Context
		username string
		opts     metav1.DeleteOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:      context.TODO(),
				username: "colin",
				opts:     metav1.DeleteOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			u := &userService{
				store: tt.fields.store,
			}
			if err := u.Delete(tt.args.ctx, tt.args.username, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("userService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_userService_Get() {
	s.mockUserStore.EXPECT().Get(gomock.Any(), "colin", gomock.Any()).Return(s.users[0], nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx      context.Context
		username string
		opts     metav1.GetOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.User
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:      context.TODO(),
				username: "colin",
				opts:     metav1.GetOptions{},
			},
			want:    s.users[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			u := &userService{
				store: tt.fields.store,
			}
			got, err := u.Get(tt.args.ctx, tt.args.username, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_userService_Update() {
	s.mockUserStore.EXPECT().Update(gomock.Any(), s.users[0], gomock.Any()).Return(nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx  context.Context
		user *v1.User
		opts metav1.UpdateOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:  context.TODO(),
				user: s.users[0],
				opts: metav1.UpdateOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			u := &userService{
				store: tt.fields.store,
			}
			if err := u.Update(tt.args.ctx, tt.args.user, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("userService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_userService_ChangePassword() {
	s.mockUserStore.EXPECT().Update(gomock.Any(), s.users[0], gomock.Any()).Return(nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx  context.Context
		user *v1.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:  context.TODO(),
				user: s.users[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			u := &userService{
				store: tt.fields.store,
			}
			if err := u.ChangePassword(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userService.ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
