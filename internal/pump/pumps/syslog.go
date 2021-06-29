// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pumps

import (
	"context"
	"fmt"
	"log/syslog"

	"github.com/mitchellh/mapstructure"

	"github.com/marmotedu/iam/internal/pump/analytics"
	"github.com/marmotedu/iam/pkg/log"
)

// SyslogPump defines a syslog pump with syslog specific options and common options.
type SyslogPump struct {
	syslogConf *SyslogConf
	writer     *syslog.Writer
	filters    analytics.AnalyticsFilters
	timeout    int
	CommonPumpConfig
}

var logPrefix = "syslog-pump"

// SyslogConf defines syslog specific options.
type SyslogConf struct {
	Transport   string `mapstructure:"transport"`
	NetworkAddr string `mapstructure:"network_addr"`
	LogLevel    int    `mapstructure:"log_level"`
	Tag         string `mapstructure:"tag"`
}

// New create a syslog pump instance.
func (s *SyslogPump) New() Pump {
	newPump := SyslogPump{}

	return &newPump
}

// GetName returns the syslog pump name.
func (s *SyslogPump) GetName() string {
	return "Syslog Pump"
}

// Init initialize the syslog pump instance.
func (s *SyslogPump) Init(config interface{}) error {
	// Read configuration file
	s.syslogConf = &SyslogConf{}
	err := mapstructure.Decode(config, &s.syslogConf)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	// Init the configs
	initConfigs(s)

	// Init the Syslog writer
	initWriter(s)

	log.Debug("Syslog Pump active")

	return nil
}

func initWriter(s *SyslogPump) {
	tag := logPrefix
	if s.syslogConf.Tag != "" {
		tag = s.syslogConf.Tag
	}
	syslogWriter, err := syslog.Dial(
		s.syslogConf.Transport,
		s.syslogConf.NetworkAddr,
		syslog.Priority(s.syslogConf.LogLevel),
		tag)
	if err != nil {
		log.Fatalf("failed to connect to Syslog Daemon: %s", err.Error())
	}

	s.writer = syslogWriter
}

// Set default values if they are not explicitly given and perform validation.
func initConfigs(pump *SyslogPump) {
	if pump.syslogConf.Transport == "" {
		pump.syslogConf.Transport = "udp"
		log.Info("No Transport given, using 'udp'")
	}

	if pump.syslogConf.Transport != "udp" &&
		pump.syslogConf.Transport != "tcp" &&
		pump.syslogConf.Transport != "tls" {
		log.Fatal("Chosen invalid Transport type.  Please use a supported Transport type for Syslog")
	}

	if pump.syslogConf.NetworkAddr == "" {
		pump.syslogConf.NetworkAddr = "localhost:5140"
		log.Info("No host given, using 'localhost:5140'")
	}

	if pump.syslogConf.LogLevel == 0 {
		log.Warn("Using Log Level 0 (KERNEL) for Syslog pump")
	}
}

// WriteData write analyzed data to syslog persistent back-end storage.
func (s *SyslogPump) WriteData(ctx context.Context, data []interface{}) error {
	// Data is all the analytics being written
	for _, v := range data {
		select {
		case <-ctx.Done():
			return nil
		default:
			// Decode the raw analytics into Form
			decoded, _ := v.(analytics.AnalyticsRecord)
			message := Message{
				"timestamp":  decoded.TimeStamp,
				"username":   decoded.Username,
				"effect":     decoded.Effect,
				"conclusion": decoded.Conclusion,
				"request":    decoded.Request,
				"policies":   decoded.Policies,
				"deciders":   decoded.Deciders,
				"expireAt":   decoded.ExpireAt,
			}

			// Print to Syslog
			_, _ = fmt.Fprintf(s.writer, "%s", message)
		}
	}

	return nil
}

// SetTimeout set attributes `timeout` for SyslogPump.
func (s *SyslogPump) SetTimeout(timeout int) {
	s.timeout = timeout
}

// GetTimeout get attributes `timeout` for SyslogPump.
func (s *SyslogPump) GetTimeout() int {
	return s.timeout
}

// SetFilters set attributes `filters` for SyslogPump.
func (s *SyslogPump) SetFilters(filters analytics.AnalyticsFilters) {
	s.filters = filters
}

// GetFilters get attributes `filters` for SyslogPump.
func (s *SyslogPump) GetFilters() analytics.AnalyticsFilters {
	return s.filters
}
