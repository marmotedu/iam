// Copyright 2020 Lingfei Kong <marmotedu@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package watcher

import (
	"net/http"

	"github.com/marmotedu/iam/internal/watcher/config"
	"github.com/marmotedu/iam/pkg/log"
)

// Run runs the specified pump server. This should never exit.
func Run(cfg *config.Config) error {
	go serveHealthCheck(cfg.HealthCheckPath, cfg.HealthCheckAddress)

	return createWatcherServer(cfg).PrepareRun().Run()
}

// serveHealthCheck runs a http server used to provide a api to check pump health status.
func serveHealthCheck(healthPath string, healthAddress string) {
	http.HandleFunc("/"+healthPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	if err := http.ListenAndServe(healthAddress, nil); err != nil {
		log.Fatalf("Error serving health check endpoint: %s", err.Error())
	}
}
