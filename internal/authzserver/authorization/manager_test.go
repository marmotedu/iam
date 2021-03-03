// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authorization

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/ory/ladon"
)

func TestNewPolicyManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)

	type args struct {
		client AuthorizationInterface
	}
	tests := []struct {
		name string
		args args
		want ladon.Manager
	}{
		{
			name: "default",
			args: args{
				client: mockAuthz,
			},
			want: &PolicyManager{mockAuthz},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPolicyManager(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPolicyManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolicyManager_FindRequestCandidates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)
	policy := ladon.DefaultPolicy{
		ID:          "68819e5a-738b-41ec-b03c-b58a1b19d043",
		Description: "One policy to rule them all.",
		Subjects:    []string{"users:<peter|ken>", "users:maria", "groups:admins"},
		Resources:   []string{"resources:articles:<.*>", "resources:printer"},
		Actions:     []string{"delete", "<create|update>"},
		Effect:      ladon.AllowAccess,
		Conditions:  ladon.Conditions{"remoteIPAddress": &ladon.CIDRCondition{CIDR: "192.168.0.1/16"}},
	}
	mockAuthz.EXPECT().List(gomock.Any()).Return([]*ladon.DefaultPolicy{&policy}, nil)

	type args struct {
		r *ladon.Request
	}
	tests := []struct {
		name    string
		args    args
		want    ladon.Policies
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				r: &ladon.Request{},
			},
			want:    []ladon.Policy{&policy},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewPolicyManager(mockAuthz)
			got, err := m.FindRequestCandidates(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("PolicyManager.FindRequestCandidates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PolicyManager.FindRequestCandidates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolicyManager_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)

	type fields struct {
		client AuthorizationInterface
	}
	type args struct {
		policy ladon.Policy
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
				client: mockAuthz,
			},
			args: args{
				policy: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PolicyManager{
				client: tt.fields.client,
			}
			if err := p.Create(tt.args.policy); (err != nil) != tt.wantErr {
				t.Errorf("PolicyManager.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPolicyManager_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)

	type fields struct {
		client AuthorizationInterface
	}
	type args struct {
		policy ladon.Policy
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
				client: mockAuthz,
			},
			args: args{
				policy: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PolicyManager{
				client: tt.fields.client,
			}
			if err := p.Update(tt.args.policy); (err != nil) != tt.wantErr {
				t.Errorf("PolicyManager.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPolicyManager_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)

	type fields struct {
		client AuthorizationInterface
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ladon.Policy
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				client: mockAuthz,
			},
			args: args{
				id: "test",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PolicyManager{
				client: tt.fields.client,
			}
			got, err := p.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PolicyManager.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PolicyManager.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolicyManager_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)

	type fields struct {
		client AuthorizationInterface
	}
	type args struct {
		id string
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
				client: mockAuthz,
			},
			args: args{
				id: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PolicyManager{
				client: tt.fields.client,
			}
			if err := p.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("PolicyManager.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPolicyManager_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)

	type fields struct {
		client AuthorizationInterface
	}
	type args struct {
		limit  int64
		offset int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ladon.Policies
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				client: mockAuthz,
			},
			args:    args{},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PolicyManager{
				client: tt.fields.client,
			}
			got, err := p.GetAll(tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("PolicyManager.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PolicyManager.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolicyManager_FindPoliciesForSubject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)

	type fields struct {
		client AuthorizationInterface
	}
	type args struct {
		subject string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ladon.Policies
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				client: mockAuthz,
			},
			args: args{
				subject: "",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PolicyManager{
				client: tt.fields.client,
			}
			got, err := m.FindPoliciesForSubject(tt.args.subject)
			if (err != nil) != tt.wantErr {
				t.Errorf("PolicyManager.FindPoliciesForSubject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PolicyManager.FindPoliciesForSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolicyManager_FindPoliciesForResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)

	type fields struct {
		client AuthorizationInterface
	}
	type args struct {
		resource string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ladon.Policies
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				client: mockAuthz,
			},
			args: args{
				resource: "",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PolicyManager{
				client: tt.fields.client,
			}
			got, err := m.FindPoliciesForResource(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("PolicyManager.FindPoliciesForResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PolicyManager.FindPoliciesForResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
