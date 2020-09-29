// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

// Validate checks PumpOptions and return a slice of found errs.
func (s *PumpOptions) Validate() []error {
	var errs []error

	errs = append(errs, s.RedisOptions.Validate()...)
	errs = append(errs, s.Log.Validate()...)

	return errs
}
