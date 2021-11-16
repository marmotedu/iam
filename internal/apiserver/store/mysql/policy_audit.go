// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysql

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type policyAudit struct {
	db *gorm.DB
}

func newPolicyAudits(ds *datastore) *policyAudit {
	return &policyAudit{ds.db}
}

// ClearOutdated clear data older than a given days.
func (p *policyAudit) ClearOutdated(ctx context.Context, maxReserveDays int) (int64, error) {
	date := time.Now().AddDate(0, 0, -maxReserveDays).Format("2006-01-02 15:04:05")

	d := p.db.Exec("delete from policy_audit where deletedAt < ?", date)

	return d.RowsAffected, d.Error
}
