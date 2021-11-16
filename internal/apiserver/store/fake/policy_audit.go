// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fake

import (
	"context"
)

type policyAudit struct {
	ds *datastore
}

func newPolicyAudits(ds *datastore) *policyAudit {
	return &policyAudit{ds}
}

// ClearOutdated clear data older than a given days.
func (p *policyAudit) ClearOutdated(ctx context.Context, maxReserveDays int) (int64, error) {
	return 0, nil
}
