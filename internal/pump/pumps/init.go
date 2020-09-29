// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pumps

var availablePumps map[string]Pump

// nolint: gochecknoinits
func init() {
	availablePumps = make(map[string]Pump)

	// Register all the storage handlers here
	availablePumps["csv"] = &CSVPump{}
	availablePumps["mongo"] = &MongoPump{}
	availablePumps["dummy"] = &DummyPump{}
	availablePumps["elasticsearch"] = &ElasticsearchPump{}
	availablePumps["influx"] = &InfluxPump{}
	availablePumps["prometheus"] = &PrometheusPump{}
	availablePumps["kafka"] = &KafkaPump{}
	availablePumps["syslog"] = &SyslogPump{}
}
