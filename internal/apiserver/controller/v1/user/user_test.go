// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	srvv1 "github.com/marmotedu/iam/internal/apiserver/service/v1"
	"github.com/marmotedu/iam/internal/apiserver/store"
)

func TestNewUserController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFactory := store.NewMockFactory(ctrl)

	type args struct {
		store store.Factory
	}
	tests := []struct {
		name string
		args args
		want *UserController
	}{
		{
			name: "default",
			args: args{
				store: mockFactory,
			},
			want: &UserController{
				srv: srvv1.NewService(mockFactory),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserController(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserController() = %v, want %v", got, tt.want)
			}
		})
	}
}
