// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package etcd

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/marmotedu/errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	"github.com/marmotedu/iam/internal/apiserver/store"
	genericoptions "github.com/marmotedu/iam/internal/pkg/options"
	"github.com/marmotedu/iam/pkg/log"
)

// EtcdCreateEventFunc defines etcd create event function handler.
type EtcdCreateEventFunc func(ctx context.Context, key, value []byte)

// EtcdModifyEventFunc defines etcd update event function handler.
type EtcdModifyEventFunc func(ctx context.Context, key, oldvalue, value []byte)

// EtcdDeleteEventFunc defines etcd delete event function handler.
type EtcdDeleteEventFunc func(ctx context.Context, key []byte)

// EtcdWatcher defines a etcd watcher.
type EtcdWatcher struct {
	watcher clientv3.Watcher
	cancel  context.CancelFunc
}

type datastore struct {
	cli             *clientv3.Client
	requestTimeout  time.Duration
	leaseTTLTimeout int

	leaseID            clientv3.LeaseID
	onKeepaliveFailure func()
	leaseLiving        bool

	watchers  map[string]*EtcdWatcher
	namespace string
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

func (ds *datastore) PolicyAudits() store.PolicyAuditStore {
	return newPolicyAudits(ds)
}

// Close clsoe the etcdStore clinet.
func (ds *datastore) Close() error {
	if ds.cli != nil {
		return ds.cli.Close()
	}

	return nil
}

func defaultOnKeepAliveFailed() {
	log.Warn("etcdStore keepalive failed")
}

var (
	etcdFactory store.Factory
	once        sync.Once
)

// GetEtcdFactoryOr create a etcdFactory store with given options.
func GetEtcdFactoryOr(opt *genericoptions.EtcdOptions, onKeepaliveFailure func()) (store.Factory, error) {
	if opt == nil && etcdFactory == nil {
		return nil, fmt.Errorf("failed to get etcd store fatory")
	}

	var err error
	once.Do(func() {
		var (
			tlsConfig *tls.Config
			cli       *clientv3.Client
		)
		tlsConfig, err = opt.GetEtcdTLSConfig()
		if err != nil {
			return
		}

		if opt.UseTLS && tlsConfig == nil {
			err = fmt.Errorf("enable etcdFactory tls but tls config is empty")

			return
		}

		ds := &datastore{}
		if onKeepaliveFailure == nil {
			onKeepaliveFailure = defaultOnKeepAliveFailed
		}
		ds.onKeepaliveFailure = onKeepaliveFailure

		cli, err = clientv3.New(clientv3.Config{
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
			return
		}

		ds.cli = cli
		ds.requestTimeout = time.Duration(opt.RequestTimeout) * time.Second
		ds.leaseTTLTimeout = opt.LeaseExpire
		ds.watchers = make(map[string]*EtcdWatcher)
		ds.namespace = opt.Namespace

		err = ds.startSession()
		if err != nil {
			if e := ds.Close(); e != nil {
				log.Errorf("etcdStore client close failed %s", e)
			}

			return
		}
		etcdFactory = ds
	})

	if etcdFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get etcd store fatory, etcdFactory: %+v, error: %w", etcdFactory, err)
	}

	return etcdFactory, nil
}

func (ds *datastore) startSession() error {
	ctx := context.TODO()

	resp, err := ds.cli.Grant(ctx, int64(ds.leaseTTLTimeout))
	if err != nil {
		return errors.Wrap(err, "creates new lease failed")
	}
	ds.leaseID = resp.ID

	ch, err := ds.cli.KeepAlive(ctx, ds.leaseID)
	if err != nil {
		return errors.Wrap(err, "keep alive failed")
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

func (ds *datastore) Client() *clientv3.Client {
	return ds.cli
}

func (ds *datastore) SessionLiving() bool {
	return ds.leaseLiving
}

func (ds *datastore) RestartSession() error {
	if ds.leaseLiving {
		return fmt.Errorf("session is living, can't restart")
	}

	return ds.startSession()
}

func (ds *datastore) getKey(key string) string {
	if len(ds.namespace) > 0 {
		return fmt.Sprintf("%s%s", ds.namespace, key)
	}

	return key
}

func (ds *datastore) Put(ctx context.Context, key string, val string) error {
	return ds.put(ctx, key, val, false)
}

func (ds *datastore) PutSession(ctx context.Context, key string, val string) error {
	return ds.put(ctx, key, val, true)
}

func (ds *datastore) grantLease(ctx context.Context, ttlSeconds int64) (*clientv3.LeaseGrantResponse, error) {
	nctx, cancel := context.WithTimeout(ctx, ds.requestTimeout)
	defer cancel()
	resp, err := ds.cli.Grant(nctx, ttlSeconds)
	if err != nil {
		return nil, fmt.Errorf("grant lease: %w", err)
	}

	return resp, nil
}

func (ds *datastore) PutWithLease(ctx context.Context, key string, val string, ttlSeconds int64) error {
	resp, err := ds.grantLease(ctx, ttlSeconds)
	if err != nil {
		return fmt.Errorf("put with grant lease: %w", err)
	}

	nctx, cancel := context.WithTimeout(ctx, ds.requestTimeout)
	defer cancel()

	key = ds.getKey(key)
	leaseID := resp.ID
	opts := []clientv3.OpOption{
		clientv3.WithLease(leaseID),
	}
	if _, err = ds.cli.Put(nctx, key, val, opts...); err != nil {
		return errors.Wrap(err, "put key-value pair to etcd failed")
	}

	return nil
}

func (ds *datastore) put(ctx context.Context, key string, val string, session bool) error {
	nctx, cancel := context.WithTimeout(ctx, ds.requestTimeout)
	defer cancel()

	key = ds.getKey(key)
	if session {
		if _, err := ds.cli.Put(nctx, key, val, clientv3.WithLease(ds.leaseID)); err != nil {
			return errors.Wrap(err, "put key-value pair to etcd failed")
		}

		return nil
	}

	if _, err := ds.cli.Put(nctx, key, val); err != nil {
		return errors.Wrap(err, "put key-value pair to etcd failed")
	}

	return nil
}

func (ds *datastore) Get(ctx context.Context, key string) ([]byte, error) {
	nctx, cancel := context.WithTimeout(ctx, ds.requestTimeout)
	defer cancel()

	key = ds.getKey(key)

	resp, err := ds.cli.Get(nctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "get key from etcd failed")
	}
	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("no such key")
	}

	return resp.Kvs[0].Value, nil
}

// EtcdKeyValue defines etcd returned key-value pairs.
type EtcdKeyValue struct {
	Key   string
	Value []byte
}

func (ds *datastore) List(ctx context.Context, prefix string) ([]EtcdKeyValue, error) {
	nctx, cancel := context.WithTimeout(ctx, ds.requestTimeout)
	defer cancel()

	prefix = ds.getKey(prefix)

	resp, err := ds.cli.Get(nctx, prefix, clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return nil, errors.Wrap(err, "get key from etcd failed")
	}
	ret := make([]EtcdKeyValue, len(resp.Kvs))
	for i := 0; i < len(resp.Kvs); i++ {
		ret[i] = EtcdKeyValue{
			Key:   string(resp.Kvs[i].Key[len(ds.namespace):]),
			Value: resp.Kvs[i].Value,
		}
	}

	return ret, nil
}

// Cancel cancel etcd client.
func (w *EtcdWatcher) Cancel() {
	w.watcher.Close()
	w.cancel()
}

// Watch watch etcd.
func (ds *datastore) Watch(
	ctx context.Context,
	prefix string,
	onCreate EtcdCreateEventFunc,
	onModify EtcdModifyEventFunc,
	onDelete EtcdDeleteEventFunc,
) error {
	if _, ok := ds.watchers[prefix]; ok {
		return fmt.Errorf("watch prefix %s already registered", prefix)
	}

	watcher := clientv3.NewWatcher(ds.cli)
	nctx, cancel := context.WithCancel(ctx)

	ds.watchers[prefix] = &EtcdWatcher{
		watcher: watcher,
		cancel:  cancel,
	}

	prefix = ds.getKey(prefix)

	rch := watcher.Watch(nctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	go func() {
		for wresp := range rch {
			for _, ev := range wresp.Events {
				key := ev.Kv.Key[len(ds.namespace):]
				if ev.PrevKv == nil {
					onCreate(nctx, key, ev.Kv.Value)
				} else {
					switch ev.Type {
					case mvccpb.PUT:
						onModify(nctx, key, ev.PrevKv.Value, ev.Kv.Value)
					case mvccpb.DELETE:
						if onDelete != nil {
							onDelete(nctx, key)
						}
					}
				}
			}
		}
		log.Infof("stop watching %s", prefix)
	}()

	return nil
}

func (ds *datastore) Unwatch(prefix string) {
	watcher, ok := ds.watchers[prefix]
	if ok {
		log.Debugf("unwatch %s", prefix)
		watcher.Cancel()
		delete(ds.watchers, prefix)
	} else {
		log.Debugf("prefix %s not watched!!", prefix)
	}
}

func (ds *datastore) Delete(ctx context.Context, key string) ([]byte, error) {
	nctx, cancel := context.WithTimeout(ctx, ds.requestTimeout)
	defer cancel()

	key = ds.getKey(key)

	dresp, err := ds.cli.Delete(nctx, key, clientv3.WithPrevKV())
	if err != nil {
		return nil, errors.Wrap(err, "delete key from etcd failed")
	}

	if dresp.Deleted == 1 {
		return dresp.PrevKvs[0].Value, nil
	}

	return nil, nil
}
