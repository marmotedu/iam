// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

//go:generate mockgen -destination mock_cacheclient.go -package store github.com/marmotedu/api/proto/apiserver/v1 CacheClient

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/AlekSi/pointer"
	pb "github.com/marmotedu/api/proto/apiserver/v1"
	"github.com/ory/ladon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/marmotedu/iam/pkg/log"
)

// GRPCClient defines a grpc client used to get all secrets and policies.
type GRPCClient struct {
	cli pb.CacheClient
}

var _ StoreClient = &GRPCClient{}

var once sync.Once
var client *GRPCClient

// GetGRPCClientOrDie return cache instance and panics on any error.
func GetGRPCClientOrDie(address string, clientCA string) StoreClient {
	if address != "" && clientCA != "" {
		once.Do(func() {
			var (
				err   error
				conn  *grpc.ClientConn
				creds credentials.TransportCredentials
			)

			creds, err = credentials.NewClientTLSFromFile(clientCA, "")
			if err != nil {
				log.Panicf("credentials.NewClientTLSFromFile err: %v", err)
			}

			conn, err = grpc.Dial(address, grpc.WithBlock(), grpc.WithTransportCredentials(creds))
			if err != nil {
				log.Panicf("Connect to grpc server failed, error: %s", err.Error())
			}

			client = &GRPCClient{pb.NewCacheClient(conn)}
			log.Infof("Connected to grpc server, address: %s", address)
		})
	}

	return client
}

// GetSecrets returns all the authorization secrets.
func (c *GRPCClient) GetSecrets() (map[string]*pb.SecretInfo, error) {
	secrets := make(map[string]*pb.SecretInfo)

	log.Info("Loading secrets")

	req := &pb.ListSecretsRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	}

	resp, err := c.cli.ListSecrets(context.Background(), req)
	if err != nil {
		return nil, err
	}

	log.Infof("Secrets found (%d total):", len(resp.Items))

	for _, v := range resp.Items {
		log.Infof(" - %s", v.SecretId)
		secrets[v.SecretId] = v
	}

	return secrets, err
}

// GetPolicies returns all the authorization policies.
func (c *GRPCClient) GetPolicies() (map[string][]*ladon.DefaultPolicy, error) {
	pols := make(map[string][]*ladon.DefaultPolicy)

	log.Info("Loading policies")

	req := &pb.ListPoliciesRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	}

	resp, err := c.cli.ListPolicies(context.Background(), req)
	if err != nil {
		return nil, err
	}

	log.Infof("Policies found (%d total)[username:name]:", len(resp.Items))

	for _, v := range resp.Items {
		log.Infof(" - %s:%s", v.Username, v.Name)

		var policy ladon.DefaultPolicy

		if err := json.Unmarshal([]byte(v.PolicyStr), &policy); err != nil {
			log.Warnf("failed to load policy for %s, error: %s", v.Name, err.Error())
			continue
		}

		pols[v.Username] = append(pols[v.Username], &policy)
	}

	return pols, nil
}
