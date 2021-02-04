// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

import (
	"os"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	pb "github.com/marmotedu/api/proto/apiserver/v1"
	"github.com/ory/ladon"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockStoreClient := NewMockStoreClient(ctrl)
	secret := &pb.SecretInfo{
		Name:        "secret1",
		SecretId:    "zMsG7HklAmRORjGtfazMbnGCLaHkdfOxGdk4",
		Username:    "colin",
		SecretKey:   "kJYkcJNHLbaN9urj5dmZ37ZwWyg1p0pm",
		Expires:     0,
		Description: "secret test",
		CreatedAt:   "2020-08-27 13:55:16",
		UpdatedAt:   "2020-08-27 13:55:16",
	}

	mockStoreClient.EXPECT().GetSecrets().AnyTimes().Return(map[string]*pb.SecretInfo{
		"zMsG7HklAmRORjGtfazMbnGCLaHkdfOxGdk4": secret,
	}, nil)
	mockStoreClient.EXPECT().GetPolicies().AnyTimes().Return(map[string][]*ladon.DefaultPolicy{
		"colin": {{}},
	}, nil)
	s, _ := GetStoreInsOr(mockStoreClient)

	s.Reload()

	os.Exit(m.Run())
}

func TestGetStoreInsOr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStoreClient := NewMockStoreClient(ctrl)

	type args struct {
		cli StoreClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				cli: mockStoreClient,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetStoreInsOr(tt.args.cli)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStoreInsOr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStore_GetSecret(t *testing.T) {
	s, err := GetStoreInsOr(nil)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.SecretInfo
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				key: "zMsG7HklAmRORjGtfazMbnGCLaHkdfOxGdk4",
			},
			want: &pb.SecretInfo{
				Name:        "secret1",
				SecretId:    "zMsG7HklAmRORjGtfazMbnGCLaHkdfOxGdk4",
				Username:    "colin",
				SecretKey:   "kJYkcJNHLbaN9urj5dmZ37ZwWyg1p0pm",
				Expires:     0,
				Description: "secret test",
				CreatedAt:   "2020-08-27 13:55:16",
				UpdatedAt:   "2020-08-27 13:55:16",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetSecret(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GetSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Store.GetSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_GetPolicy(t *testing.T) {
	s, err := GetStoreInsOr(nil)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []*ladon.DefaultPolicy
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				key: "colin",
			},
			want:    []*ladon.DefaultPolicy{{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetPolicy(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GetPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Store.GetPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_Reload(t *testing.T) {
	s, err := GetStoreInsOr(nil)
	assert.Nil(t, err)
	assert.NotNil(t, s)
	assert.Nil(t, s.Reload())
}
