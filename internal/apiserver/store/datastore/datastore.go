// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package datastore

import (
	"fmt"

	"github.com/jinzhu/gorm"

	v1 "github.com/marmotedu/api/apiserver/v1"
	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/pkg/options"

	// MySQL driver.
	_ "github.com/go-sql-driver/mysql"
)

type datastore struct {
	*gorm.DB

	// can include two database instance if needed
	// docker *grom.DB
	// db *gorm.DB
}

func (ds *datastore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *datastore) Secrets() store.SecretStore {
	return newSecrets(ds)
}

func (ds *datastore) Policies() store.PolicyStore {
	return newPolicies(ds)
}

// NewMySQLStore create mysql store with the given config.
func NewMySQLStore(o *options.MySQLOptions) (store.Store, error) {
	config := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		o.Username,
		o.Password,
		o.Host,
		o.Database,
		true,
		"Local")

	db, err := gorm.Open("mysql", config)
	if err != nil {
		return nil, err
	}

	setupDatabase(db, o)

	return &datastore{db}, nil
}

// setupDatabase initialize the database tables.
func setupDatabase(db *gorm.DB, o *options.MySQLOptions) {
	// uncomment the following line if you need auto migration the given models
	// not suggested in production environment.
	// migrateDatabase(db)

	db.LogMode(o.LogMode)
	db.DB().SetMaxOpenConns(o.MaxOpenConnections)
	db.DB().SetConnMaxLifetime(o.MaxConnectionLifeTime)
	db.DB().SetMaxIdleConns(o.MaxIdleConnections)
}

// cleanDatabase tear downs the database tables.
// nolint:unused // may be reused in the feature, or just show a migrate usage.
func cleanDatabase(db *gorm.DB) {
	db.DropTable(&v1.User{})
	db.DropTable(&v1.Policy{})
	db.DropTable(&v1.Secret{})
}

// migrateDatabase run auto migration for given models, will only add missing fields,
// won't delete/change current data.
// nolint:unused // may be reused in the feature, or just show a migrate usage.
func migrateDatabase(db *gorm.DB) {
	db.AutoMigrate(&v1.User{})
	db.AutoMigrate(&v1.Policy{})
	db.AutoMigrate(&v1.Secret{})
}

// resetDatabase resets the database tables.
// nolint:unused,deadcode // may be reused in the feature, or just show a migrate usage.
func resetDatabase(db *gorm.DB, o *options.MySQLOptions) {
	cleanDatabase(db)
	migrateDatabase(db)
	setupDatabase(db, o)
}
