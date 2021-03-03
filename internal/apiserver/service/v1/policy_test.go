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
	"github.com/stretchr/testify/suite"

	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/apiserver/store/fake"
)

type Suite struct {
	suite.Suite
	mockFactory *store.MockFactory

	mockPolicyStore *store.MockPolicyStore
	policies        []*v1.Policy

	mockSecretStore *store.MockSecretStore
	secrets         []*v1.Secret

	mockUserStore *store.MockUserStore
	users         []*v1.User
}

func (s *Suite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.policies = fake.FakePolicies(10)
	s.secrets = fake.FakeSecrets(10)
	s.users = fake.FakeUsers(10)
	s.mockFactory = store.NewMockFactory(ctrl)
	s.mockPolicyStore = store.NewMockPolicyStore(ctrl)
	s.mockFactory.EXPECT().Policies().AnyTimes().Return(s.mockPolicyStore)

	s.mockSecretStore = store.NewMockSecretStore(ctrl)
	s.mockFactory.EXPECT().Secrets().AnyTimes().Return(s.mockSecretStore)

	s.mockUserStore = store.NewMockUserStore(ctrl)
	s.mockFactory.EXPECT().Users().AnyTimes().Return(s.mockUserStore)
}

func TestPolicy(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_policyService_Create() {
	s.mockPolicyStore.EXPECT().Create(gomock.Any(), gomock.Eq(s.policies[0]), gomock.Any()).Return(nil)
	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx    context.Context
		policy *v1.Policy
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
				policy: s.policies[0],
				opts:   metav1.CreateOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ss := &policyService{
				store: tt.fields.store,
			}
			if err := ss.Create(tt.args.ctx, tt.args.policy, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("policyService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_policyService_Update() {
	s.mockPolicyStore.EXPECT().Update(gomock.Any(), gomock.Eq(s.policies[0]), gomock.Any()).Return(nil)

	type fields struct {
		store store.Factory
	}
	type args struct {
		ctx    context.Context
		policy *v1.Policy
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
				policy: s.policies[0],
				opts:   metav1.UpdateOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &policyService{
				store: tt.fields.store,
			}
			if err := s.Update(tt.args.ctx, tt.args.policy, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("policyService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_policyService_Delete() {
	s.mockPolicyStore.EXPECT().Delete(gomock.Any(), gomock.Eq("admin"), "policy0", gomock.Any()).Return(nil)

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
				name:     "policy0",
				opts:     metav1.DeleteOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &policyService{
				store: tt.fields.store,
			}
			if err := s.Delete(tt.args.ctx, tt.args.username, tt.args.name, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("policyService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_policyService_DeleteCollection() {
	s.mockPolicyStore.EXPECT().DeleteCollection(
		gomock.Any(),
		gomock.Eq("admin"),
		[]string{"policy1", "policy2"},
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
				names:    []string{"policy1", "policy2"},
				opts:     metav1.DeleteOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &policyService{
				store: tt.fields.store,
			}
			if err := s.DeleteCollection(tt.args.ctx, tt.args.username, tt.args.names, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("policyService.DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (s *Suite) Test_policyService_Get() {
	s.mockPolicyStore.EXPECT().Get(gomock.Any(), gomock.Eq("admin"), "policy0", gomock.Any()).Return(s.policies[0], nil)

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
		want    *v1.Policy
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
				name:     "policy0",
				opts:     metav1.GetOptions{},
			},
			want:    s.policies[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &policyService{
				store: tt.fields.store,
			}
			got, err := s.Get(tt.args.ctx, tt.args.username, tt.args.name, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("policyService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("policyService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *Suite) Test_policyService_List() {
	policies := &v1.PolicyList{
		ListMeta: metav1.ListMeta{
			TotalCount: 10,
		},
		Items: s.policies,
	}
	s.mockPolicyStore.EXPECT().List(gomock.Any(), gomock.Eq("admin"), gomock.Any()).Return(policies, nil)

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
		want    *v1.PolicyList
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
			want:    policies,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s := &policyService{
				store: tt.fields.store,
			}
			got, err := s.List(tt.args.ctx, tt.args.username, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("policyService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("policyService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newPolicies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)

	type args struct {
		srv *service
	}
	tests := []struct {
		name string
		args args
		want *policyService
	}{
		{
			name: "default",
			args: args{
				srv: &service{store: mockFactory},
			},
			want: &policyService{
				store: mockFactory,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPolicies(tt.args.srv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPolicies() = %v, want %v", got, tt.want)
			}
		})
	}
}
