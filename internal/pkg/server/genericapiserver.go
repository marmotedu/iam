// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"fmt"
	"net/http"

	"strings"
	"sync"
	"time"

	// limits "github.com/gin-contrib/size".

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/component-base/pkg/version"

	"github.com/marmotedu/iam/pkg/log"

	"github.com/marmotedu/iam/internal/pkg/middleware"
)

// GenericAPIServer contains state for a iam api server.
// type GenericAPIServer gin.Engine.
type GenericAPIServer struct {
	middlewares []string
	mode        string
	// SecureServingInfo holds configuration of the TLS server.
	SecureServingInfo *SecureServingInfo

	// InsecureServingInfo holds configuration of the insecure HTTP server.
	InsecureServingInfo *InsecureServingInfo

	maxPingCount int
	// ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
	// gracefully shutdown returns.
	ShutdownTimeout time.Duration

	*gin.Engine
	healthz         bool
	enableMetrics   bool
	enableProfiling bool
	// wrapper for gin.Engine
}

// InstallAPIs install generic apis.
func (s *GenericAPIServer) InstallAPIs() {
	// install healthz handler
	if s.healthz {
		s.GET("/healthz", func(c *gin.Context) {
			core.WriteResponse(c, nil, map[string]string{"status": "ok"})
		})
	}

	// install metric handler
	if s.enableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.Engine)
	}

	// install pprof handler
	if s.enableProfiling {
		pprof.Register(s.Engine)
	}

	s.GET("/version", func(c *gin.Context) {
		core.WriteResponse(c, nil, version.Get())
	})
}

// Setup do some setup work for gin engine.
func (s *GenericAPIServer) Setup() {
	gin.SetMode(s.mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

// InstallMiddlewares install generic middlewares.
func (s *GenericAPIServer) InstallMiddlewares() {
	// necessary middlewares
	s.Use(middleware.RequestID())

	// install custom middlewares
	for _, m := range s.middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware: %s", m)
			continue
		}

		log.Infof("install middleware: %s", m)
		s.Use(mw)
	}

	s.Use(middleware.Context())
	// s.Use(gin.Logger())
	// s.Use(limits.RequestSizeLimiter(10))
	// s.GET("/debug/vars", expvar.Handler())
}

// Run spawns the http server. It only returns if stopCh is closed
// or the port cannot be listened on initially.
func (s *GenericAPIServer) Run(stopCh <-chan struct{}) error {
	insecureServer := &http.Server{
		Addr:    s.InsecureServingInfo.Address,
		Handler: s,
	}

	secureServer := &http.Server{
		Addr:    s.SecureServingInfo.Address(),
		Handler: s,
	}

	wg := sync.WaitGroup{}

	wg.Add(2)
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		defer wg.Done()
		log.Infof("Start to listening the incoming requests on http address: %s", s.InsecureServingInfo.Address)

		if err := insecureServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}

		log.Infof("Server on %s stopped", s.InsecureServingInfo.Address)
	}()

	go func() {
		defer wg.Done()

		key, cert := s.SecureServingInfo.CertKey.KeyFile, s.SecureServingInfo.CertKey.CertFile
		if cert == "" || key == "" || s.SecureServingInfo.BindPort == 0 {
			return
		}

		log.Infof("Start to listening the incoming requests on https address: %s", s.SecureServingInfo.Address())

		if err := secureServer.ListenAndServeTLS(cert, key); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}

		log.Infof("Server on %s stopped", s.SecureServingInfo.Address())
	}()

	// Ping the server to make sure the router is working.
	if s.healthz {
		if err := s.pingGenericAPIServer(stopCh); err != nil {
			return err
		}
	}

	<-stopCh

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	// nolint: gomnd
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := secureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown secure server failed: %s", err.Error())
	}

	if err := insecureServer.Shutdown(ctx); err != nil {
		log.Warnf("Shutdown insecure server failed: %s", err.Error())
	}

	wg.Wait()
	// log.Info("All servers stopped. Exiting.")

	return nil
}

// pingGenericAPIServer pings the http server to make sure the router is working.
func (s *GenericAPIServer) pingGenericAPIServer(stopCh <-chan struct{}) error {
	url := fmt.Sprintf("http://%s/healthz", s.InsecureServingInfo.Address)
	if strings.Contains(s.InsecureServingInfo.Address, "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(s.InsecureServingInfo.Address, ":")[1])
	}

	for i := 0; i < s.maxPingCount; i++ {
		// Ping the server by sending a GET request to `/healthz`.
		// nolint: gosec
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Info("The router has been deployed successfully.")

			resp.Body.Close()

			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)

		select {
		case <-stopCh:
			log.Warn("Ping server stoped.")
			return nil
		default:
		}
	}

	return fmt.Errorf("the router has no response, or it might took too long to start up")
}
