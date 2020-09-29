// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	v1 "github.com/marmotedu/api/apiserver/v1"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/pkg/code"
	"github.com/marmotedu/iam/internal/pkg/middleware"
	"github.com/marmotedu/log"
)

type auth struct {
	realm      string
	key        []byte
	timeout    time.Duration
	maxRefresh time.Duration
}

type jwtAuth struct {
	realm      string
	key        []byte
	timeout    time.Duration
	maxRefresh time.Duration
}

func newAPIServerAuth(realm string, key []byte, timeout, maxRefresh time.Duration) middleware.AuthInterface {
	return &auth{
		realm:      realm,
		key:        key,
		timeout:    timeout,
		maxRefresh: maxRefresh,
	}
}

func (a *auth) JWTAuth() middleware.JWTAuthInterface {
	return newAuthMiddleware(a.realm, a.key, a.timeout, a.maxRefresh)
}

func (a *auth) BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			core.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."), nil)
			c.Abort()

			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !authenticateUser(pair[0], pair[1]) {
			core.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."), nil)
			c.Abort()

			return
		}

		c.Request.Header.Add("username", pair[0])

		c.Next()
	}
}

// Use gin binding feature to validate the input.
type basicAuth struct {
	Username string `form:"username" json:"username" binding:"required,username"`
	Password string `form:"password" json:"password" binding:"required,password"`
}

func newAuthMiddleware(realm string, key []byte, timeout, maxRefresh time.Duration) middleware.JWTAuthInterface {
	return &jwtAuth{
		realm:      realm,
		key:        key,
		timeout:    timeout,
		maxRefresh: maxRefresh,
	}
}

func (auth *jwtAuth) Realm() string {
	return auth.realm
}

func (auth *jwtAuth) Key() []byte {
	return auth.key
}

func (auth *jwtAuth) Timeout() time.Duration {
	return auth.timeout
}

func (auth *jwtAuth) MaxRefresh() time.Duration {
	return auth.maxRefresh
}

func (auth *jwtAuth) Authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var login basicAuth
		if err := c.ShouldBindJSON(&login); err != nil {
			log.Errorf("parse login parameters: %s", err.Error())
			return "", jwt.ErrMissingLoginValues
		}

		username := login.Username
		password := login.Password

		// Get the user information by the login username.
		user, err := store.Client().Users().Get(username, metav1.GetOptions{})
		if err != nil {
			log.Errorf("get user information failed: %s", err.Error())
			return "", jwt.ErrFailedAuthentication
		}

		// Compare the login password with the user password.
		if err := user.Compare(password); err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		return user, nil
	}
}

func (auth *jwtAuth) LoginResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func (auth *jwtAuth) LogoutResponse() func(c *gin.Context, code int) {
	return func(c *gin.Context, code int) {
		c.JSON(http.StatusOK, nil)
	}
}

func (auth *jwtAuth) PayloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		claims := jwt.MapClaims{
			"iss": "iam-apiserver",
			"sub": "user of iam-apiserver",
		}
		if u, ok := data.(*v1.User); ok {
			claims[jwt.IdentityKey] = u.Name
		}

		return claims
	}
}

func (auth *jwtAuth) IdentityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return claims[jwt.IdentityKey]
	}
}

func (auth *jwtAuth) Authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		// add username to header
		if v, ok := data.(string); ok {
			c.Request.Header.Add("username", v)
			return true
		}

		return false
	}
}

func (auth *jwtAuth) Unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"message": message,
		})
	}
}

func authenticateUser(username, password string) bool {
	// fetch user from database
	user, err := store.Client().Users().Get(username, metav1.GetOptions{})
	if err != nil {
		return false
	}

	// Compare the login password with the user password.
	if err := user.Compare(password); err != nil {
		return false
	}

	return true
}
