// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysql

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	v1 "github.com/marmotedu/api/apiserver/v1"

	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/pkg/logger"
	genericoptions "github.com/marmotedu/iam/internal/pkg/options"
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

var mysqlFactory store.Factory
var once sync.Once

// GetMySQLFactoryOr create mysql factory with the given config.
func GetMySQLFactoryOr(opt *genericoptions.MySQLOptions) (store.Factory, error) {
	if opt == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	once.Do(func() {
		var db *gorm.DB
		dns := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
			opt.Username,
			opt.Password,
			opt.Host,
			opt.Database,
			true,
			"Local")

		db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
			Logger: logger.New(opt.LogLevel),
		})
		if err != nil {
			return
		}

		err = setupDatabase(opt, db)
		if err != nil {
			return
		}

		mysqlFactory = &datastore{db}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %v", mysqlFactory, err)
	}

	return mysqlFactory, nil
}

// setupDatabase initialize the database tables.
func setupDatabase(opt *genericoptions.MySQLOptions, db *gorm.DB) error {
	// uncomment the following line if you need auto migration the given models
	// not suggested in production environment.
	// migrateDatabase(db)

	// db.LogMode(opt.LogMode)
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(opt.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(opt.MaxConnectionLifeTime)
	sqlDB.SetMaxIdleConns(opt.MaxIdleConnections)
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
func resetDatabase(db *gorm.DB, opt *genericoptions.MySQLOptions) error {
	if err := cleanDatabase(db); err != nil {
		return err
	}
	if err := migrateDatabase(db); err != nil {
		return err
	}
	if err := setupDatabase(opt, db); err != nil {
		return err
	}

	return nil
}
