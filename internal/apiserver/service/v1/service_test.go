// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"reflect"
	"testing"

	"github.com/marmotedu/iam/internal/apiserver/store"
)

func TestNewService(t *testing.T) {
	type args struct {
		store store.Factory
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Users(t *testing.T) {
	type fields struct {
		store store.Factory
	}
	tests := []struct {
		name   string
		fields fields
		want   UserSrv
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				store: tt.fields.store,
			}
			if got := s.Users(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Users() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Secrets(t *testing.T) {
	type fields struct {
		store store.Factory
	}
	tests := []struct {
		name   string
		fields fields
		want   SecretSrv
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				store: tt.fields.store,
			}
			if got := s.Secrets(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Secrets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Policies(t *testing.T) {
	type fields struct {
		store store.Factory
	}
	tests := []struct {
		name   string
		fields fields
		want   PolicySrv
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				store: tt.fields.store,
			}
			if got := s.Policies(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Policies() = %v, want %v", got, tt.want)
			}
		})
	}
}
