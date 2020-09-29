// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package server runs a http server used to provide a api to check pump health status.
package server

import (
	"fmt"
	"net/http"

	"github.com/marmotedu/log"
)

var defaultHealthEndpoint = "healthz"
var defaultHealthPort = 7070

// ServeHealthCheck runs a http server used to provide a api to check pump health status.
func ServeHealthCheck(configHealthEndpoint string, configHealthPort int) {
	healthEndpoint := configHealthEndpoint
	if healthEndpoint == "" {
		healthEndpoint = defaultHealthEndpoint
	}
	healthPort := configHealthPort
	if healthPort == 0 {
		healthPort = defaultHealthPort
	}

	http.HandleFunc("/"+healthEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	log.Infof("Serving health check endpoint at http://localhost:%d/%s ...", healthPort, healthEndpoint)
	if err := http.ListenAndServe(":"+fmt.Sprint(healthPort), nil); err != nil {
		log.Fatalf("Error serving health check endpoint: %s", err.Error())
	}
}
