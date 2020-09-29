// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pumps

import "github.com/marmotedu/iam/internal/pump/analytics"

// CommonPumpConfig defines common options used by all persistent store, like elasticsearch, kafka, mongo and etc.
type CommonPumpConfig struct {
	filters               analytics.AnalyticsFilters
	timeout               int
	OmitDetailedRecording bool
}

// SetFilters set attributes `filters` for CommonPumpConfig.
func (p *CommonPumpConfig) SetFilters(filters analytics.AnalyticsFilters) {
	p.filters = filters
}

// GetFilters get attributes `filters` for CommonPumpConfig.
func (p *CommonPumpConfig) GetFilters() analytics.AnalyticsFilters {
	return p.filters
}

// SetTimeout set attributes `timeout` for CommonPumpConfig.
func (p *CommonPumpConfig) SetTimeout(timeout int) {
	p.timeout = timeout
}

// GetTimeout get attributes `timeout` for CommonPumpConfig.
func (p *CommonPumpConfig) GetTimeout() int {
	return p.timeout
}

// SetOmitDetailedRecording set attributes `OmitDetailedRecording` for CommonPumpConfig.
func (p *CommonPumpConfig) SetOmitDetailedRecording(omitDetailedRecording bool) {
	p.OmitDetailedRecording = omitDetailedRecording
}

// GetOmitDetailedRecording get attributes `OmitDetailedRecording` for CommonPumpConfig.
func (p *CommonPumpConfig) GetOmitDetailedRecording() bool {
	return p.OmitDetailedRecording
}
