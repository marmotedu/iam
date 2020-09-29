// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"net"

	"google.golang.org/grpc"

	"github.com/marmotedu/log"
)

type grpcAPIServer struct {
	*grpc.Server
	address string
}

func (s *grpcAPIServer) Run(stopCh <-chan struct{}) {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}

	log.Infof("Start grpc server at %s", s.address)

	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatalf("failed to start grpc server: %s", err.Error())
		}
	}()

	<-stopCh

	log.Infof("Grpc server on %s stopped", s.address)
	s.GracefulStop()
}
