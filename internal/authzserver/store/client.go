// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/AlekSi/pointer"
	jsoniter "github.com/json-iterator/go"
	"github.com/ory/ladon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/marmotedu/api/proto/apiserver/v1"
	"github.com/marmotedu/log"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	rpcConnectMu      sync.Mutex
	clientIsConnected bool
	client            atomic.Value
)

// GrpcClient is a storage manager that uses the redis database.
type GrpcClient struct {
	Addr     string
	ClientCA string
}

func getClient() pb.CacheClient {
	v := client.Load()
	if v != nil {
		return v.(pb.CacheClient)
	}

	return nil
}

// Connect will establish a connection to the RPC.
func (c *GrpcClient) Connect() bool {
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
func (c *GrpcClient) GetSecrets() (map[string]*pb.SecretInfo, error) {
	secrets := make(map[string]*pb.SecretInfo)

	log.Info("Loading secrets")

	body := &pb.ListSecretsRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	}

	resp, err := getClient().ListSecrets(context.Background(), body)
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
func (c *GrpcClient) GetPolicies() (map[string][]*ladon.DefaultPolicy, error) {
	pols := make(map[string][]*ladon.DefaultPolicy)

	log.Info("Loading policies")

	body := &pb.ListPoliciesRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	}

	resp, err := getClient().ListPolicies(context.Background(), body)
	if err != nil {
		return nil, err
	}

	log.Infof("Policies found (%d total)[username - name]:", len(resp.Items))

	for _, v := range resp.Items {
		log.Infof("%s - %s", v.Username, v.Name)

		var policy ladon.DefaultPolicy

		if err := json.Unmarshal([]byte(v.PolicyStr), &policy); err != nil {
			log.Warnf("failed to load policy for %s, error: %s", v.Name, err.Error())
			continue
		}

		pols[v.Username] = append(pols[v.Username], &policy)
	}

	return pols, nil
}
