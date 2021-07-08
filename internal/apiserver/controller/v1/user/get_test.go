// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	v1 "github.com/marmotedu/api/apiserver/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"

	srvv1 "github.com/marmotedu/iam/internal/apiserver/service/v1"
)

func TestUserController_Get(t *testing.T) {
	user := &v1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: "admin",
			ID:   0,
		},
		Nickname: "admin",
		Password: "Admin@2020",
		Email:    "admin@foxmail.com",
	}

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/v1/users/colin", nil)
	c.Params = []gin.Param{{Key: "name", Value: "admin"}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvv1.NewMockService(ctrl)
	mockUserSrv := srvv1.NewMockUserSrv(ctrl)
	mockUserSrv.EXPECT().Get(gomock.Any(), gomock.Eq("admin"), gomock.Any()).Return(user, nil)
	mockService.EXPECT().Users().Return(mockUserSrv)

	type fields struct {
		srv srvv1.Service
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "default",
			fields: fields{
				srv: mockService,
			},
			args: args{
				c: c,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserController{
				srv: tt.fields.srv,
			}
			u.Get(tt.args.c)
		})
	}
}
