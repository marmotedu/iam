// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"context"
	"reflect"
	"testing"

	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"

	"github.com/marmotedu/iam/internal/apiserver/store"
)

func Test_newPolicies(t *testing.T) {
	type args struct {
		srv *service
	}
	tests := []struct {
		name string
		args args
		want *policyService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPolicies(tt.args.srv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPolicies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_policyService_Create(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &policyService{
				store: tt.fields.store,
			}
			if err := s.Create(tt.args.ctx, tt.args.policy, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("policyService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_policyService_Update(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &policyService{
				store: tt.fields.store,
			}
			if err := s.Update(tt.args.ctx, tt.args.policy, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("policyService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_policyService_Delete(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &policyService{
				store: tt.fields.store,
			}
			if err := s.Delete(tt.args.ctx, tt.args.username, tt.args.name, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("policyService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_policyService_DeleteCollection(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &policyService{
				store: tt.fields.store,
			}
			if err := s.DeleteCollection(tt.args.ctx, tt.args.username, tt.args.names, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("policyService.DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_policyService_Get(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func Test_policyService_List(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
