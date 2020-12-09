// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"

	"github.com/marmotedu/iam/internal/apiserver/store"
	genericoptions "github.com/marmotedu/iam/internal/pkg/options"
	"github.com/marmotedu/iam/pkg/log"
)

type datastore struct {
	cli *clientv3.Client

	requestTimeout  time.Duration
	leaseTTLTimeout int

	leaseID            clientv3.LeaseID
	onKeepaliveFailure func()
	leaseLiving        bool

	// watchers map[string]*SEtcdWatcher
	// namespace string
}

func (ds *datastore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *datastore) Secrets() store.SecretStore {
	return newSecrets(ds)
}

func (ds *datastore) Policies() store.PolicyStore {
	return newPolicies(ds)
}

func defaultOnKeepAliveFailed() {
	log.Warn("etcdStore keepalive failed")
}

// NewEtcdStore create a etcd store with given options.
func NewEtcdStore(opt *genericoptions.EtcdOptions, onKeepaliveFailure func()) (store.Store, error) {
	tlsConfig, err := opt.GetEtcdTLSConfig()
	if err != nil {
		return nil, err
	}

	if opt.UseTLS && tlsConfig == nil {
		return nil, fmt.Errorf("enable etcdStore tls but tls config is empty")
	}

	etcdClient := &datastore{}
	if onKeepaliveFailure == nil {
		onKeepaliveFailure = defaultOnKeepAliveFailed
	}
	etcdClient.onKeepaliveFailure = onKeepaliveFailure

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   opt.Endpoints,
		DialTimeout: time.Duration(opt.Timeout) * time.Second,
		Username:    opt.Username,
		Password:    opt.Password,
		TLS:         tlsConfig,

		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
		},
	})
	if err != nil {
		return nil, err
	}

	etcdClient.cli = cli
	etcdClient.requestTimeout = time.Duration(opt.RequestTimeout) * time.Second
	etcdClient.leaseTTLTimeout = opt.LeaseExpire
	// etcdClient.watchers = make(map[string]*SEtcdWatcher)

	if err := etcdClient.startSession(); err != nil {
		if e := etcdClient.Close(); e != nil {
			log.Errorf("etcdStore client close failed %s", e)
		}
		return nil, err
	}

	return etcdClient, nil
}

func (ds *datastore) startSession() error {
	ctx := context.TODO()

	resp, err := ds.cli.Grant(ctx, int64(ds.leaseTTLTimeout))
	if err != nil {
		return err
	}
	ds.leaseID = resp.ID

	ch, err := ds.cli.KeepAlive(ctx, ds.leaseID)
	if err != nil {
		return err
	}
	ds.leaseLiving = true

	go func() {
		for {
			if _, ok := <-ch; !ok {
				ds.leaseLiving = false
				log.Errorf("fail to keepalive session")
				if ds.onKeepaliveFailure != nil {
					ds.onKeepaliveFailure()
				}
				break
			}
		}
	}()

	return nil
}

// Close clsoe the etcdStore clinet.
func (ds *datastore) Close() error {
	if ds.cli != nil {
		err := ds.cli.Close()
		if err != nil {
			return err
		}
		ds.cli = nil
	}

	return nil
}
