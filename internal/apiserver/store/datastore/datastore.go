// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package datastore

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	v1 "github.com/marmotedu/api/apiserver/v1"
	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/pkg/logger"
	"github.com/marmotedu/iam/internal/pkg/options"
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
	dns := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		o.Username,
		o.Password,
		o.Host,
		o.Database,
		true,
		"Local")

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: logger.New(o.LogLevel),
	})
	if err != nil {
		return nil, err
	}

	if err := setupDatabase(db, o); err != nil {
		return nil, err
	}

	return &datastore{db}, nil
}

// setupDatabase initialize the database tables.
func setupDatabase(db *gorm.DB, o *options.MySQLOptions) error {
	// uncomment the following line if you need auto migration the given models
	// not suggested in production environment.
	// migrateDatabase(db)

	//db.LogMode(o.LogMode)
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(o.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(o.MaxConnectionLifeTime)
	sqlDB.SetMaxIdleConns(o.MaxIdleConnections)
	return nil
}

// cleanDatabase tear downs the database tables.
// nolint:unused // may be reused in the feature, or just show a migrate usage.
func cleanDatabase(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&v1.User{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&v1.Policy{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&v1.Secret{}); err != nil {
		return err
	}

	return nil
}

// migrateDatabase run auto migration for given models, will only add missing fields,
// won't delete/change current data.
// nolint:unused // may be reused in the feature, or just show a migrate usage.
func migrateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(&v1.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&v1.Policy{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&v1.Secret{}); err != nil {
		return err
	}

	return nil
}

// resetDatabase resets the database tables.
// nolint:unused,deadcode // may be reused in the feature, or just show a migrate usage.
func resetDatabase(db *gorm.DB, o *options.MySQLOptions) error {
	if err := cleanDatabase(db); err != nil {
		return err
	}
	if err := migrateDatabase(db); err != nil {
		return err
	}
	if err := setupDatabase(db, o); err != nil {
		return err
	}

	return nil
}
