// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pumps

import (
	"context"
	"errors"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/ory/ladon"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/marmotedu/iam/internal/pump/analytics"
	"github.com/marmotedu/iam/pkg/log"
)

// PrometheusPump defines a prometheus pump with prometheus specific options and common options.
type PrometheusPump struct {
	conf *PrometheusConf
	// Per service
	TotalStatusMetrics *prometheus.CounterVec

	CommonPumpConfig
}

// PrometheusConf defines prometheus specific options.
type PrometheusConf struct {
	Addr string `mapstructure:"listen_address"`
	Path string `mapstructure:"path"`
}

// New create a prometheus pump instance.
func (p *PrometheusPump) New() Pump {
	newPump := PrometheusPump{}
	newPump.TotalStatusMetrics = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "iam_user_authorization_status_total",
			Help: "authorization effect per user",
		},
		[]string{"code", "username"},
	)

	prometheus.MustRegister(newPump.TotalStatusMetrics)

	return &newPump
}

// GetName returns the prometheus pump name.
func (p *PrometheusPump) GetName() string {
	return "Prometheus Pump"
}

// Init initialize the prometheus pump instance.
func (p *PrometheusPump) Init(conf interface{}) error {
	p.conf = &PrometheusConf{}
	err := mapstructure.Decode(conf, &p.conf)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	if p.conf.Path == "" {
		p.conf.Path = "/metrics"
	}

	if p.conf.Addr == "" {
		return errors.New("prometheus listen_addr not set")
	}

	log.Infof("Starting prometheus listener on: %s", p.conf.Addr)

	http.Handle(p.conf.Path, promhttp.Handler())

	go func() {
		log.Fatal(http.ListenAndServe(p.conf.Addr, nil).Error())
	}()

	return nil
}

// WriteData write analyzed data to prometheus persistent back-end storage.
func (p *PrometheusPump) WriteData(ctx context.Context, data []interface{}) error {
	log.Debugf("Writing %d records", len(data))

	for _, item := range data {
		record, _ := item.(analytics.AnalyticsRecord)
		code := "0"
		if record.Effect != ladon.AllowAccess {
			code = "1"
		}

		p.TotalStatusMetrics.WithLabelValues(code, record.Username).Inc()
	}

	return nil
}
