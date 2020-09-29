// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

// Validate checks ServerRunOptions and return a slice of found errs.
func (s *ServerRunOptions) Validate() []error {
	var errs []error

	errs = append(errs, s.GenericServerRunOptions.Validate()...)
	errs = append(errs, s.InsecureServing.Validate()...)
	errs = append(errs, s.SecureServing.Validate()...)
	errs = append(errs, s.RedisOptions.Validate()...)
	errs = append(errs, s.FeatureOptions.Validate()...)
	errs = append(errs, s.Log.Validate()...)
	errs = append(errs, s.AnalyticsOptions.Validate()...)

	return errs
}
