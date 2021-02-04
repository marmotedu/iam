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

func TestNewAuditLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)

	type args struct {
		client AuthorizationInterface
	}
	tests := []struct {
		name string
		args args
		want *AuditLogger
	}{
		{
			name: "default",
			args: args{
				client: mockAuthz,
			},
			want: &AuditLogger{mockAuthz},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuditLogger(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuditLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuditLogger_LogRejectedAccessRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)
	mockAuthz.EXPECT().LogRejectedAccessRequest(gomock.Any(), gomock.Any(), gomock.Any())

	type args struct {
		r *ladon.Request
		p ladon.Policies
		d ladon.Policies
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "default",
			args: args{
				r: &ladon.Request{},
				p: ladon.Policies{},
				d: ladon.Policies{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAuditLogger(mockAuthz)
			a.LogRejectedAccessRequest(tt.args.r, tt.args.p, tt.args.d)
		})
	}
}

func TestAuditLogger_LogGrantedAccessRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthz := NewMockAuthorizationInterface(ctrl)
	mockAuthz.EXPECT().LogGrantedAccessRequest(gomock.Any(), gomock.Any(), gomock.Any())

	type args struct {
		r *ladon.Request
		p ladon.Policies
		d ladon.Policies
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "default",
			args: args{
				r: &ladon.Request{},
				p: ladon.Policies{},
				d: ladon.Policies{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAuditLogger(mockAuthz)
			a.LogGrantedAccessRequest(tt.args.r, tt.args.p, tt.args.d)
		})
	}
}
