// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

import (
	"encoding/json"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	pb "github.com/marmotedu/api/proto/apiserver/v1"
	"github.com/ory/ladon"
	"github.com/stretchr/testify/assert"
)

func TestGetGRPCClientOrDie(t *testing.T) {
	client := GetGRPCClientOrDie("", "")
	assert.Nil(t, client)
}

func TestGRPCClient_GetSecrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCacheClient := NewMockCacheClient(ctrl)

	secret1 := &pb.SecretInfo{
		Name:        "secret1",
		SecretId:    "zMsG7HklAmRORjGtfazMbnGCLaHkdfOxGdk4",
		Username:    "colin",
		SecretKey:   "kJYkcJNHLbaN9urj5dmZ37ZwWyg1p0pm",
		Expires:     0,
		Description: "secret test",
		CreatedAt:   "2020-08-27 13:55:16",
		UpdatedAt:   "2020-08-27 13:55:16",
	}
	secret2 := &pb.SecretInfo{
		Name:        "secret2",
		SecretId:    "zMsG7HklDmRORjGtfazMbnGCLaHkdfOxGdk4",
		Username:    "colin",
		SecretKey:   "kJYkcJNHDbaN9urj5dmZ37ZwWyg1p0pm",
		Expires:     0,
		Description: "secret test",
		CreatedAt:   "2020-08-27 16:53:16",
		UpdatedAt:   "2020-08-27 16:55:16",
	}
	mockCacheClient.EXPECT().ListSecrets(gomock.Any(), gomock.Any()).Return(&pb.ListSecretsResponse{
		TotalCount: 2,
		Items:      []*pb.SecretInfo{secret1, secret2},
	}, nil)

	tests := []struct {
		name    string
		want    map[string]*pb.SecretInfo
		wantErr bool
	}{
		{
			name: "default",
			want: map[string]*pb.SecretInfo{
				secret1.SecretId: secret1,
				secret2.SecretId: secret2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GRPCClient{
				cli: mockCacheClient,
			}
			got, err := c.GetSecrets()
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCClient.GetSecrets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCClient.GetSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGRPCClient_GetPolicies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCacheClient := NewMockCacheClient(ctrl)
	policyStr := "{}"
	var policy ladon.DefaultPolicy
	_ = json.Unmarshal([]byte(policyStr), &policy)

	policy1 := &pb.PolicyInfo{
		Name:      "policy1",
		Username:  "colin",
		PolicyStr: policyStr,
		CreatedAt: "2020-08-27 13:55:16",
	}
	policy2 := &pb.PolicyInfo{
		Name:      "policy2",
		Username:  "peter",
		PolicyStr: policyStr,
		CreatedAt: "2020-08-27 13:55:16",
	}
	mockCacheClient.EXPECT().ListPolicies(gomock.Any(), gomock.Any()).Return(&pb.ListPoliciesResponse{
		TotalCount: 2,
		Items:      []*pb.PolicyInfo{policy1, policy2},
	}, nil)

	tests := []struct {
		name    string
		want    map[string][]*ladon.DefaultPolicy
		wantErr bool
	}{
		{
			name: "default",
			want: map[string][]*ladon.DefaultPolicy{
				policy1.Username: {
					&policy,
				},
				policy2.Username: {
					&policy,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &GRPCClient{
				cli: mockCacheClient,
			}
			got, err := c.GetPolicies()
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCClient.GetPolicies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCClient.GetPolicies() = %v, want %v", got, tt.want)
			}
		})
	}
}
