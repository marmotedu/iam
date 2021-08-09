// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authzserver

import (
	"github.com/marmotedu/errors"

	"github.com/marmotedu/iam/internal/authzserver/load/cache"
	"github.com/marmotedu/iam/internal/pkg/middleware"
	"github.com/marmotedu/iam/internal/pkg/middleware/auth"
)

func newCacheAuth() middleware.AuthStrategy {
	return auth.NewCacheStrategy(getSecretFunc())
}

func getSecretFunc() func(string) (auth.Secret, error) {
	return func(kid string) (auth.Secret, error) {
		cli, err := cache.GetCacheInsOr(nil)
		if err != nil || cli == nil {
			return auth.Secret{}, errors.Wrap(err, "get cache instance failed")
		}

		secret, err := cli.GetSecret(kid)
		if err != nil {
			return auth.Secret{}, err
		}

		return auth.Secret{
			Username: secret.Username,
			ID:       secret.SecretId,
			Key:      secret.SecretKey,
			Expires:  secret.Expires,
		}, nil
	}
}
