// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/AlekSi/pointer"
	pb "github.com/marmotedu/api/proto/apiserver/v1"
	"github.com/marmotedu/component-base/pkg/json"
	"github.com/ory/ladon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/marmotedu/iam/pkg/log"
)

var (
	rpcConnectMu      sync.Mutex
	clientIsConnected bool
	client            atomic.Value
)

// GRPCClient is a storage manager that uses the redis database.
type GRPCClient struct {
	Addr     string
	ClientCA string
}

// Client returns grpc client.
func (c *GRPCClient) Client() pb.CacheClient {
	v := client.Load()
	if v != nil {
		return v.(pb.CacheClient)
	}

	return nil
}

// Connect will establish a connection to the RPC.
func (c *GRPCClient) Connect() bool {
	log.Debug("connecting to grpc server in block mode")
	rpcConnectMu.Lock()
	defer rpcConnectMu.Unlock()

	if clientIsConnected {
		return true
	}

	creds, err := credentials.NewClientTLSFromFile(c.ClientCA, "")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}

	conn, err := grpc.Dial(c.Addr, grpc.WithBlock(), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Connect to grpc server failed, error: %s", err.Error())
	}

	log.Infof("Connected to grpc server, address: %s", c.Addr)

	client.Store(pb.NewCacheClient(conn))

	clientIsConnected = true

	return true
}

// GetSecrets returns all the authorization secrets.
func (c *GRPCClient) GetSecrets() (map[string]*pb.SecretInfo, error) {
	secrets := make(map[string]*pb.SecretInfo)

	log.Info("Loading secrets")

	req := &pb.ListSecretsRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	}

	resp, err := c.Client().ListSecrets(context.Background(), req)
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

	resp, err := c.Client().ListPolicies(context.Background(), req)
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
