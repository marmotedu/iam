// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/internal/apiserver/api/v1/policy"
	"github.com/marmotedu/iam/internal/apiserver/api/v1/secret"
	"github.com/marmotedu/iam/internal/apiserver/api/v1/user"
	"github.com/marmotedu/iam/internal/pkg/code"
	"github.com/marmotedu/iam/internal/pkg/middleware"

	// custom gin validators.
	_ "github.com/marmotedu/iam/pkg/validator"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installAPI(g)
}

func installMiddleware(g *gin.Engine) {
}

func installAPI(g *gin.Engine) *gin.Engine {
	// Middlewares.
	// the jwt middleware
	apiServerAuth := newAPIServerAuth(
		viper.GetString("jwt.Realm"),
		[]byte(viper.GetString("jwt.key")),
		viper.GetDuration("jwt.timeout"),
		viper.GetDuration("jwt.max-refresh"),
	)

	authMiddleware, _ := middleware.NewAuthMiddleware(apiServerAuth, nil)

	g.POST("/login", authMiddleware.JWT.LoginHandler)
	g.POST("/logout", authMiddleware.JWT.LogoutHandler)
	// Refresh time can be longer than token timeout
	g.POST("/refresh", authMiddleware.JWT.RefreshHandler)

	g.NoRoute(authMiddleware.AuthMiddlewareFunc(), func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})

	// v1 handlers, requiring authentication
	v1 := g.Group("/v1")
	{
		// user RESTful resource
		userv1 := v1.Group("/users")
		{
			userv1.POST("", user.Create)

			userv1.Use(authMiddleware.AuthMiddlewareFunc(), middleware.Validation())
			// v1.PUT("/find_password", user.FindPassword)
			userv1.DELETE("", user.DeleteCollection) // admin api
			userv1.DELETE(":name", user.Delete)      // admin api
			userv1.PUT(":name/change_password", user.ChangePassword)
			userv1.PUT(":name", user.Update)
			userv1.GET("", user.List)
			userv1.GET(":name", user.Get) // admin api
		}

		v1.Use(authMiddleware.AuthMiddlewareFunc())

		// policy RESTful resource
		policyv1 := v1.Group("/policies", middleware.Publish())
		{
			policyv1.POST("", policy.Create)
			policyv1.DELETE("", policy.DeleteCollection)
			policyv1.DELETE(":name", policy.Delete)
			policyv1.PUT(":name", policy.Update)
			policyv1.GET("", policy.List)
			policyv1.GET(":name", policy.Get)
		}

		// secret RESTful resource
		secretv1 := v1.Group("/secrets", middleware.Publish())
		{
			secretv1.POST("", secret.Create)
			secretv1.DELETE(":name", secret.Delete)
			secretv1.PUT(":name", secret.Update)
			secretv1.GET("", secret.List)
			secretv1.GET(":name", secret.Get)
		}
	}

	return g
}
