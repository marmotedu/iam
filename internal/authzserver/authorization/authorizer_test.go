// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authorization

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	authzv1 "github.com/marmotedu/api/authz/v1"
	"github.com/ory/ladon"
)

func TestNewAuthorizer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)

	type args struct {
		authorizationClient AuthorizationInterface
	}
	tests := []struct {
		name string
		args args
		want *Authorizer
	}{
		{
			name: "default",
			args: args{
				authorizationClient: mockAuthz,
			},
			want: &Authorizer{
				warden: &ladon.Ladon{
					Manager:     NewPolicyManager(mockAuthz),
					AuditLogger: NewAuditLogger(mockAuthz),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthorizer(tt.args.authorizationClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthorizer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorizer_Authorize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)

	mockAuthz.EXPECT().LogRejectedAccessRequest(gomock.Any(), gomock.Any(), gomock.Any()).Times(3)
	mockAuthz.EXPECT().LogGrantedAccessRequest(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)
	gomock.InOrder(
		mockAuthz.EXPECT().List(gomock.Any()).Return([]*ladon.DefaultPolicy{}, nil),
		mockAuthz.EXPECT().List(gomock.Any()).Times(2).Return([]*ladon.DefaultPolicy{{
			ID:          "68819e5a-738b-41ec-b03c-b58a1b19d043",
			Description: "One policy to rule them all.",
			Subjects:    []string{"users:<peter|ken>", "users:maria", "groups:admins"},
			Resources:   []string{"resources:articles:<.*>", "resources:printer"},
			Actions:     []string{"delete", "<create|update>"},
			Effect:      ladon.AllowAccess,
			Conditions:  ladon.Conditions{"remoteIPAddress": &ladon.CIDRCondition{CIDR: "192.168.0.1/16"}},
		}}, nil),
		mockAuthz.EXPECT().List(gomock.Eq("colin")).Return([]*ladon.DefaultPolicy{}, nil),
	)

	type args struct {
		request *ladon.Request
	}
	tests := []struct {
		name string
		args args
		want *authzv1.Response
	}{
		{
			name: "deny",
			args: args{
				request: &ladon.Request{
					Subject:  "users:peter",
					Action:   "delete",
					Resource: "resources:articles:ladon-introduction",
					Context: ladon.Context{
						"remoteIPAddress": "192.168.0.5",
					},
				},
			},
			want: &authzv1.Response{
				Denied: true,
				Reason: "Request was denied by default",
			},
		},
		{
			name: "allow",
			args: args{
				request: &ladon.Request{
					Subject:  "users:peter",
					Action:   "delete",
					Resource: "resources:articles:ladon-introduction",
					Context: ladon.Context{
						"remoteIPAddress": "192.168.0.5",
					},
				},
			},
			want: &authzv1.Response{
				Allowed: true,
			},
		},
		{
			name: "deny_with_policy",
			args: args{
				request: &ladon.Request{
					Subject:  "users:colin",
					Action:   "delete",
					Resource: "resources:articles:ladon-introduction",
					Context: ladon.Context{
						"remoteIPAddress": "192.168.0.5",
					},
				},
			},
			want: &authzv1.Response{
				Denied: true,
				Reason: "Request was denied by default",
			},
		},
		{
			name: "deny_with_username",
			args: args{
				request: &ladon.Request{
					Subject:  "users:colin",
					Action:   "delete",
					Resource: "resources:articles:ladon-introduction",
					Context: ladon.Context{
						"remoteIPAddress": "192.168.0.5",
						"username":        "colin",
					},
				},
			},
			want: &authzv1.Response{
				Denied: true,
				Reason: "Request was denied by default",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAuthorizer(mockAuthz)
			if got := a.Authorize(tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Authorizer.Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}
