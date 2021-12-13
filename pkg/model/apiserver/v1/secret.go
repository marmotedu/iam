// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"github.com/marmotedu/component-base/pkg/json"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/component-base/pkg/util/idutil"
	"gorm.io/gorm"
)

// Secret represents a secret restful resource.
// It is also used as gorm model.
type Secret struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `       json:"metadata,omitempty"`
	Username          string `json:"username"           gorm:"column:username"  validate:"omitempty"`
	SecretID          string `json:"secretID"           gorm:"column:secretID"  validate:"omitempty"`
	SecretKey         string `json:"secretKey"          gorm:"column:secretKey" validate:"omitempty"`

	// Required: true
	Expires     int64  `json:"expires"     gorm:"column:expires"     validate:"omitempty"`
	Description string `json:"description" gorm:"column:description" validate:"description"`
}

// SecretList is the whole list of all secrets which have been stored in stroage.
type SecretList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	metav1.ListMeta `json:",inline"`

	// List of secrets
	Items []*Secret `json:"items"`
}

// TableName maps to mysql table name.
func (s *Secret) TableName() string {
	return "secret"
}

// BeforeCreate run before create database record.
func (s *Secret) BeforeCreate(tx *gorm.DB) (err error) {
	s.SecretID = idutil.NewSecretID()
	s.SecretKey = idutil.NewSecretKey()
	s.ExtendShadow = s.Extend.String()

	return
}

// AfterCreate run after create database record.
func (s *Secret) AfterCreate(tx *gorm.DB) (err error) {
	s.InstanceID = idutil.GetInstanceID(s.ID, "secret-")

	return tx.Save(s).Error
}

// BeforeUpdate run before update database record.
func (s *Secret) BeforeUpdate(tx *gorm.DB) (err error) {
	s.ExtendShadow = s.Extend.String()

	return err
}

// AfterFind run after find to unmarshal a extend shadown string into metav1.Extend struct.
func (s *Secret) AfterFind(tx *gorm.DB) (err error) {
	if err := json.Unmarshal([]byte(s.ExtendShadow), &s.Extend); err != nil {
		return err
	}

	return nil
}
