// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package gormutil is a util to convert offset and limit to default values.
package gormutil

import (
	"reflect"
	"testing"

	"github.com/AlekSi/pointer"
)

func TestUnpointer(t *testing.T) {
	type args struct {
		offset *int64
		limit  *int64
	}
	tests := []struct {
		name string
		args args
		want *LimitAndOffset
	}{
		{
			name: "both offset and limit are not zero",
			args: args{
				offset: pointer.ToInt64(0),
				limit:  pointer.ToInt64(10),
			},
			want: &LimitAndOffset{
				Offset: 0,
				Limit:  10,
			},
		},
		{
			name: "both offset and limit are zero",
			want: &LimitAndOffset{
				Offset: 0,
				Limit:  1000,
			},
		},
		{
			name: "offset not zero and limit zero",
			args: args{
				offset: pointer.ToInt64(2),
			},
			want: &LimitAndOffset{
				Offset: 2,
				Limit:  1000,
			},
		},
		{
			name: "offset zero and limit not zero",
			args: args{
				limit: pointer.ToInt64(10),
			},
			want: &LimitAndOffset{
				Offset: 0,
				Limit:  10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unpointer(tt.args.offset, tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unpointer() = %v, want %v", got, tt.want)
			}
		})
	}
}
