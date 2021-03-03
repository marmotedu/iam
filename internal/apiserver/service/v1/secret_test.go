// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"context"
	reflect "reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"

	"github.com/marmotedu/iam/internal/apiserver/store"
)

func (s *Suite) Test_secretService_Create() {
	s.mockSecretStore.EXPECT().Create(gomock.Any(), gomock.Eq(s.secrets[0]), gomock.Any()).Return(nil)
	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx    context.Context
		secret *v1.Secret
		opts   metav1.CreateOptions
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
				ctx:    context.TODO(),
				secret: s.secrets[0],
				opts:   metav1.CreateOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ss := &secretService{
				store: tt.fields.store,
			}
			if err := ss.Create(tt.args.ctx, tt.args.secret, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("secretService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_secretService_Update() {
	s.mockSecretStore.EXPECT().Update(gomock.Any(), gomock.Eq(s.secrets[0]), gomock.Any()).Return(nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx    context.Context
		secret *v1.Secret
		opts   metav1.UpdateOptions
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
				ctx:    context.TODO(),
				secret: s.secrets[0],
				opts:   metav1.UpdateOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &secretService{
				store: tt.fields.store,
			}
			if err := s.Update(tt.args.ctx, tt.args.secret, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("secretService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_secretService_Delete() {
	s.mockSecretStore.EXPECT().Delete(gomock.Any(), gomock.Eq("admin"), "secret0", gomock.Any()).Return(nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx      context.Context
		username string
		name     string
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
				username: "admin",
				name:     "secret0",
				opts:     metav1.DeleteOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &secretService{
				store: tt.fields.store,
			}
			if err := s.Delete(tt.args.ctx, tt.args.username, tt.args.name, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("secretService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_secretService_DeleteCollection() {
	s.mockSecretStore.EXPECT().DeleteCollection(
		gomock.Any(),
		gomock.Eq("admin"),
		[]string{"secret1", "secret2"},
		gomock.Any(),
	).Return(
		nil,
	)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx      context.Context
		username string
		names    []string
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
				username: "admin",
				names:    []string{"secret1", "secret2"},
				opts:     metav1.DeleteOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &secretService{
				store: tt.fields.store,
			}
			if err := s.DeleteCollection(tt.args.ctx, tt.args.username, tt.args.names, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("secretService.DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_secretService_Get() {
	s.mockSecretStore.EXPECT().Get(gomock.Any(), gomock.Eq("admin"), "secret0", gomock.Any()).Return(s.secrets[0], nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx      context.Context
		username string
		name     string
		opts     metav1.GetOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.Secret
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:      context.TODO(),
				username: "admin",
				name:     "secret0",
				opts:     metav1.GetOptions{},
			},
			want:    s.secrets[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &secretService{
				store: tt.fields.store,
			}
			got, err := s.Get(tt.args.ctx, tt.args.username, tt.args.name, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("secretService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("secretService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_secretService_List() {
	secrets := &v1.SecretList{
		ListMeta: metav1.ListMeta{
			TotalCount: 10,
		},
		Items: s.secrets,
	}
	s.mockSecretStore.EXPECT().List(gomock.Any(), gomock.Eq("admin"), gomock.Any()).Return(secrets, nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx      context.Context
		username string
		opts     metav1.ListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.SecretList
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				store: s.mockFactory,
			},
			args: args{
				ctx:      context.TODO(),
				username: "admin",
				opts:     metav1.ListOptions{},
			},
			want:    secrets,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &secretService{
				store: tt.fields.store,
			}
			got, err := s.List(tt.args.ctx, tt.args.username, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("secretService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("secretService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newSecrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)

	type args struct {
		srv *service
	}
	tests := []struct {
		name string
		args args
		want *secretService
	}{
		{
			name: "default",
			args: args{
				srv: &service{store: mockFactory},
			},
			want: &secretService{
				store: mockFactory,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSecrets(tt.args.srv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}
