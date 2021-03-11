// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"reflect"
	"testing"

	"github.com/marmotedu/iam/internal/pump/options"
)

func TestCreateConfigFromOptions(t *testing.T) {
	opts := options.NewOptions()
	type args struct {
		opts *options.Options
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				opts: opts,
			},
			want:    &Config{opts},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateConfigFromOptions(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateConfigFromOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateConfigFromOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
