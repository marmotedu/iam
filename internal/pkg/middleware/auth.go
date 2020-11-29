// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"fmt"
	"strings"
	"time"

	ginjwt "github.com/appleboy/gin-jwt/v2"
	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"

	pb "github.com/marmotedu/api/proto/apiserver/v1"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	"github.com/marmotedu/iam/internal/pkg/code"
)

const (
	authHeaderCount = 2

	// AuthzAudience defines the value of jwt audience field.
	AuthzAudience = "iam.authz.marmotedu.com"
)

// Defined errors.
var (
	ErrMissingKID    = errors.New("Invalid token format: missing kid field in claims")
	ErrMissingSecret = errors.New("Can not obtain secret information from cache")
)

// JWTAuthInterface defines jwt authentication interface.
type JWTAuthInterface interface {
	Realm() string
	Key() []byte
	Timeout() time.Duration
	MaxRefresh() time.Duration
	Authenticator() func(c *gin.Context) (interface{}, error)
	LoginResponse() func(*gin.Context, int, string, time.Time)
	LogoutResponse() func(c *gin.Context, code int)
	PayloadFunc() func(data interface{}) ginjwt.MapClaims
	IdentityHandler() func(*gin.Context) interface{}
	Authorizator() func(data interface{}, c *gin.Context) bool
	Unauthorized() func(*gin.Context, int, string)
}

// AuthInterface defines interface with basic and jwt authentication method.
type AuthInterface interface {
	BasicAuth() gin.HandlerFunc
	JWTAuth() JWTAuthInterface
}

// CacheAuthInterface authentication interface for authzserver.
type CacheAuthInterface interface {
	GetSecret(secretID string) (*pb.SecretInfo, error)
	// GetKeyExpires(secretID string) (int64, error)
	// GetUsername(secretID string) (string, error)
}

// AuthMiddleware defines authentication middleware struct.
type AuthMiddleware struct {
	basic       gin.HandlerFunc
	JWT         *ginjwt.GinJWTMiddleware
	cacheClient CacheAuthInterface
}

// NewAuthMiddleware returns a new authentication middleware.
func NewAuthMiddleware(auth AuthInterface, cacheClient CacheAuthInterface) (*AuthMiddleware, error) {
	authMiddleware := &AuthMiddleware{
		cacheClient: cacheClient,
	}

	if auth != nil {
		jwtAuth, err := ginjwt.New(&ginjwt.GinJWTMiddleware{
			Realm:            auth.JWTAuth().Realm(),
			SigningAlgorithm: "HS256",
			Key:              auth.JWTAuth().Key(),
			Timeout:          auth.JWTAuth().Timeout(),
			MaxRefresh:       auth.JWTAuth().MaxRefresh(),
			Authenticator:    auth.JWTAuth().Authenticator(),   // 登陆1: 登陆认证
			LoginResponse:    auth.JWTAuth().LoginResponse(),   // 登陆3：返回
			LogoutResponse:   auth.JWTAuth().LogoutResponse(),  // 登出1
			PayloadFunc:      auth.JWTAuth().PayloadFunc(),     // 登陆2：增加payload信息
			IdentityHandler:  auth.JWTAuth().IdentityHandler(), // 解析claims
			Authorizator:     auth.JWTAuth().Authorizator(),    // 登陆成功处理，header中增加username
			Unauthorized:     auth.JWTAuth().Unauthorized(),    // 授权失败
			TokenLookup:      "header: Authorization, query: token, cookie: jwt",
			TokenHeadName:    "Bearer",
			SendCookie:       true,
			TimeFunc:         time.Now,
			// TODO: HTTPStatusMessageFunc:
		})

		if err != nil {
			return nil, err
		}

		authMiddleware.JWT = jwtAuth
		authMiddleware.basic = auth.BasicAuth()
	}

	return authMiddleware, nil
}

// AuthMiddlewareFunc defines authentication middleware which can deal
// username/password and jwt at the same time.
func (a *AuthMiddleware) AuthMiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(authHeader) != authHeaderCount {
			core.WriteResponse(c, errors.WithCode(code.ErrInvalidAuthHeader, "Authorization header format is wrong."), nil)
			c.Abort()

			return
		}

		switch authHeader[0] {
		case "Basic":
			a.basic(c)
		case "Bearer":
			a.JWT.MiddlewareFunc()(c)
		default:
			core.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "unrecognized Authorization header."), nil)
			c.Abort()

			return
		}

		c.Next()
	}
}

// AuthCacheMiddlewareFunc defines authentication middleware form authzserver.
func (a *AuthMiddleware) AuthCacheMiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if len(header) == 0 {
			core.WriteResponse(c, errors.WithCode(code.ErrMissingHeader, "Authorization header cannot be empty."), nil)
			c.Abort()

			return
		}

		var rawJWT string
		// Parse the header to get the token part.
		fmt.Sscanf(header, "Bearer %s", &rawJWT)

		// Use own validation logic, see below
		var secret *pb.SecretInfo

		claims := &jwt.MapClaims{}
		// Verify the token
		parsedT, err := jwt.ParseWithClaims(rawJWT, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is HMAC signature
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, ErrMissingKID
			}

			var err error
			secret, err = a.cacheClient.GetSecret(kid)
			if err != nil {
				return nil, ErrMissingSecret
			}

			return []byte(secret.SecretKey), nil
		}, jwt.WithAudience(AuthzAudience))
		if err != nil || !parsedT.Valid {
			core.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, err.Error()), nil)
			c.Abort()
			return
		}

		if KeyExpired(secret.Expires) {
			tm := time.Unix(secret.Expires, 0).Format("2006-01-02 15:04:05")
			core.WriteResponse(c, errors.WithCode(code.ErrExpired, "expired at: %s", tm), nil)
			c.Abort()

			return
		}

		c.Request.Header.Add("username", secret.Username)
		c.Next()
	}
}

// KeyExpired checks if a key has expired, if the value of user.SessionState.Expires is 0, it will be ignored.
func KeyExpired(expires int64) bool {
	if expires >= 1 {
		return time.Now().After(time.Unix(expires, 0))
	}

	return false
}
