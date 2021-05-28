// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysql

import (
	"context"

	v1 "github.com/marmotedu/api/apiserver/v1"
	"github.com/marmotedu/component-base/pkg/fields"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"
	"gorm.io/gorm"

	"github.com/marmotedu/iam/internal/pkg/code"
	"github.com/marmotedu/iam/internal/pkg/util/gormutil"
)

type secrets struct {
	db *gorm.DB
}

func newSecrets(ds *datastore) *secrets {
	return &secrets{ds.db}
}

// Create creates a new secret.
func (s *secrets) Create(ctx context.Context, secret *v1.Secret, opts metav1.CreateOptions) error {
	return s.db.Create(&secret).Error
}

// Update updates an secret information by the secret identifier.
func (s *secrets) Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOptions) error {
	return s.db.Save(secret).Error
}

// Delete deletes the secret by the secret identifier.
func (s *secrets) Delete(ctx context.Context, username, name string, opts metav1.DeleteOptions) error {
	if opts.Unscoped {
		s.db = s.db.Unscoped()
	}

	err := s.db.Where("username = ? and name = ?", username, name).Delete(&v1.Secret{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

// DeleteCollection batch deletes the secrets.
func (s *secrets) DeleteCollection(
	ctx context.Context,
	username string,
	names []string,
	opts metav1.DeleteOptions,
) error {
	if opts.Unscoped {
		s.db = s.db.Unscoped()
	}

	return s.db.Where("username = ? and name in (?)", username, names).Delete(&v1.Secret{}).Error
}

// Get return an secret by the secret identifier.
func (s *secrets) Get(ctx context.Context, username, name string, opts metav1.GetOptions) (*v1.Secret, error) {
	secret := &v1.Secret{}
	err := s.db.Where("username = ? and name= ?", username, name).First(&secret).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrSecretNotFound, err.Error())
		}

		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return secret, nil
}

// List return all secrets.
func (s *secrets) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error) {
	ret := &v1.SecretList{}
	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	if username != "" {
		s.db = s.db.Where("username = ?", username)
	}

	selector, _ := fields.ParseSelector(opts.FieldSelector)
	name, _ := selector.RequiresExactMatch("name")

	d := s.db.Where(" name like ?", "%"+name+"%").
		Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&ret.Items).
		Offset(-1).
		Limit(-1).
		Count(&ret.TotalCount)

	return ret, d.Error
}
