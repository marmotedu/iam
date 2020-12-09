// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package gormutil is a util to convert offset and limit to default values.
package gormutil

// DefaultLimit define the default number of records to be retrieved.
const DefaultLimit = 1000

// LimitAndOffset contains offset and limit fields.
type LimitAndOffset struct {
	Offset int
	Limit  int
}

// Unpointer fill LimitAndOffset with default values if offset/limit is nil
// or it will be filled with the passed value.
func Unpointer(offset *int64, limit *int64) *LimitAndOffset {
	var o, l int = 0, DefaultLimit

	if offset != nil {
		o = int(*offset)
	}

	if limit != nil {
		l = int(*limit)
	}

	return &LimitAndOffset{
		Offset: o,
		Limit:  l,
	}
}
