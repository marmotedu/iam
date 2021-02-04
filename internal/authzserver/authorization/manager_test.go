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

func TestPolicyManager_FindPoliciesForSubject(t *testing.T) {
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
