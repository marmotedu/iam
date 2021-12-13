// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package redis provides an redis implementation of the AnalyticsStorage storage interface.
package redis

import (
	"crypto/tls"
	"strconv"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/marmotedu/errors"
	"github.com/mitchellh/mapstructure"

	genericoptions "github.com/marmotedu/iam/internal/pkg/options"
	"github.com/marmotedu/iam/pkg/log"
)

// ------------------- REDIS CLUSTER STORAGE MANAGER -------------------------------

// RedisKeyPrefix defines prefix for iam analytics key.
const (
	RedisKeyPrefix      = "analytics-"
	defaultRedisAddress = "127.0.0.1:6379"
)

var redisClusterSingleton redis.UniversalClient

// RedisClusterStorageManager is a storage manager that uses the redis database.
type RedisClusterStorageManager struct {
	db        redis.UniversalClient
	KeyPrefix string
	HashKeys  bool
	Config    genericoptions.RedisOptions
}

// NewRedisClusterPool returns a redis cluster client.
func NewRedisClusterPool(forceReconnect bool, config genericoptions.RedisOptions) redis.UniversalClient {
	if !forceReconnect {
		if redisClusterSingleton != nil {
			log.Debug("Redis pool already INITIALIZED")

			return redisClusterSingleton
		}
	} else {
		if redisClusterSingleton != nil {
			redisClusterSingleton.Close()
		}
	}

	log.Debug("Creating new Redis connection pool")

	maxActive := 500
	if config.MaxActive > 0 {
		maxActive = config.MaxActive
	}

	timeout := 5 * time.Second

	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	var tlsConfig *tls.Config
	if config.UseSSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: config.SSLInsecureSkipVerify,
		}
	}

	var client redis.UniversalClient
	opts := &RedisOpts{
		MasterName:   config.MasterName,
		Addrs:        getRedisAddrs(config),
		DB:           config.Database,
		Password:     config.Password,
		PoolSize:     maxActive,
		IdleTimeout:  240 * time.Second,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		DialTimeout:  timeout,
		TLSConfig:    tlsConfig,
	}

	if opts.MasterName != "" {
		log.Info("--> [REDIS] Creating sentinel-backed failover client")
		client = redis.NewFailoverClient(opts.failover())
	} else if config.EnableCluster {
		log.Info("--> [REDIS] Creating cluster client")
		client = redis.NewClusterClient(opts.cluster())
	} else {
		log.Info("--> [REDIS] Creating single-node client")
		client = redis.NewClient(opts.simple())
	}

	redisClusterSingleton = client

	return client
}

func getRedisAddrs(config genericoptions.RedisOptions) (addrs []string) {
	if len(config.Addrs) != 0 {
		addrs = config.Addrs
	}

	if len(addrs) == 0 && config.Port != 0 {
		addr := config.Host + ":" + strconv.Itoa(config.Port)
		addrs = append(addrs, addr)
	}

	return addrs
}

// RedisOpts is the overridden type of redis.UniversalOptions. simple() and cluster() functions are not public
// in redis library. Therefore, they are redefined in here to use in creation of new redis cluster logic.
// We don't want to use redis.NewUniversalClient() logic.
type RedisOpts redis.UniversalOptions

func (o *RedisOpts) cluster() *redis.ClusterOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{defaultRedisAddress}
	}

	return &redis.ClusterOptions{
		Addrs:     o.Addrs,
		OnConnect: o.OnConnect,

		Password: o.Password,

		MaxRedirects:   o.MaxRedirects,
		ReadOnly:       o.ReadOnly,
		RouteByLatency: o.RouteByLatency,
		RouteRandomly:  o.RouteRandomly,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:        o.DialTimeout,
		ReadTimeout:        o.ReadTimeout,
		WriteTimeout:       o.WriteTimeout,
		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,

		TLSConfig: o.TLSConfig,
	}
}

func (o *RedisOpts) simple() *redis.Options {
	addr := defaultRedisAddress
	if len(o.Addrs) > 0 {
		addr = o.Addrs[0]
	}

	return &redis.Options{
		Addr:      addr,
		OnConnect: o.OnConnect,

		DB:       o.DB,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,

		TLSConfig: o.TLSConfig,
	}
}

func (o *RedisOpts) failover() *redis.FailoverOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{"127.0.0.1:26379"}
	}

	return &redis.FailoverOptions{
		SentinelAddrs: o.Addrs,
		MasterName:    o.MasterName,
		OnConnect:     o.OnConnect,

		DB:       o.DB,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,

		TLSConfig: o.TLSConfig,
	}
}

// GetName returns the redis cluster storage manager name.
func (r *RedisClusterStorageManager) GetName() string {
	return "redis"
}

// Init initialize the redis cluster storage manager.
func (r *RedisClusterStorageManager) Init(config interface{}) error {
	r.Config = genericoptions.RedisOptions{}
	err := mapstructure.Decode(config, &r.Config)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	r.KeyPrefix = RedisKeyPrefix

	return nil
}

// Connect will establish a connection to the r.db.
func (r *RedisClusterStorageManager) Connect() bool {
	if r.db == nil {
		log.Debug("Connecting to redis cluster")
		r.db = NewRedisClusterPool(false, r.Config)

		return true
	}

	log.Debug("Storage Engine already initialized...")

	// Reset it just in case
	r.db = redisClusterSingleton

	return true
}

func (r *RedisClusterStorageManager) hashKey(in string) string {
	return in
}

func (r *RedisClusterStorageManager) fixKey(keyName string) string {
	setKeyName := r.KeyPrefix + r.hashKey(keyName)

	log.Debugf("Input key was: %s", setKeyName)

	return setKeyName
}

// GetAndDeleteSet get and delete key from redis.
func (r *RedisClusterStorageManager) GetAndDeleteSet(keyName string) []interface{} {
	log.Debugf("Getting raw key set: %s", keyName)

	if r.db == nil {
		log.Warn("Connection dropped, connecting..")
		r.Connect()

		return r.GetAndDeleteSet(keyName)
	}

	log.Debugf("keyName is: %s", keyName)

	fixedKey := r.fixKey(keyName)

	log.Debugf("Fixed keyname is: %s", fixedKey)

	var lrange *redis.StringSliceCmd
	_, err := r.db.TxPipelined(func(pipe redis.Pipeliner) error {
		lrange = pipe.LRange(fixedKey, 0, -1)
		pipe.Del(fixedKey)

		return nil
	})
	if err != nil {
		log.Errorf("Multi command failed: %s", err)
		r.Connect()
	}

	vals := lrange.Val()

	result := make([]interface{}, len(vals))
	for i, v := range vals {
		result[i] = v
	}

	log.Debugf("Unpacked vals: %d", len(result))

	return result
}

// SetKey will create (or update) a key value in the store.
func (r *RedisClusterStorageManager) SetKey(keyName, session string, timeout int64) error {
	log.Debugf("[STORE] SET Raw key is: %s", keyName)
	log.Debugf("[STORE] Setting key: %s", r.fixKey(keyName))

	r.ensureConnection()
	err := r.db.Set(r.fixKey(keyName), session, 0).Err()
	if timeout > 0 {
		if expErr := r.SetExp(keyName, timeout); expErr != nil {
			return expErr
		}
	}
	if err != nil {
		log.Errorf("Error trying to set value: %s", err.Error())

		return errors.Wrap(err, "failed to set key")
	}

	return nil
}

// SetExp is used to set the expiry of a key.
func (r *RedisClusterStorageManager) SetExp(keyName string, timeout int64) error {
	err := r.db.Expire(r.fixKey(keyName), time.Duration(timeout)*time.Second).Err()
	if err != nil {
		log.Errorf("Could not EXPIRE key: %s", err.Error())
	}

	return errors.Wrap(err, "failed to set expire time for key")
}

func (r *RedisClusterStorageManager) ensureConnection() {
	if r.db != nil {
		// already connected
		return
	}
	log.Info("Connection dropped, reconnecting...")
	for {
		r.Connect()
		if r.db != nil {
			// reconnection worked
			return
		}
		log.Info("Reconnecting again...")
	}
}
