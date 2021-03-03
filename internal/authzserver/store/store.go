// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

//go:generate mockgen -destination mock_store.go -package store github.com/marmotedu/iam/internal/authzserver/store StoreClient

import (
	"errors"
	"sync"

	"github.com/dgraph-io/ristretto"
	pb "github.com/marmotedu/api/proto/apiserver/v1"
	"github.com/ory/ladon"
)

// StoreClient defines functions used to get all secrets and policies.
type StoreClient interface {
	GetSecrets() (map[string]*pb.SecretInfo, error)
	GetPolicies() (map[string][]*ladon.DefaultPolicy, error)
}

// Store is used to store secrets and policies.
type Store struct {
	lock     *sync.RWMutex
	cli      StoreClient
	secrets  *ristretto.Cache
	policies *ristretto.Cache
}

var (
	// ErrSecretNotFound defines secret not found error.
	ErrSecretNotFound = errors.New("secret not found")
	// ErrPolicyNotFound defines policy not found error.
	ErrPolicyNotFound = errors.New("policy not found")
)

var (
	onceStore sync.Once
	storeIns  *Store
)

// GetStoreInsOr return store instance.
func GetStoreInsOr(cli StoreClient) (*Store, error) {
	var err error
	if cli != nil {
		var (
			secretCache *ristretto.Cache
			policyCache *ristretto.Cache
		)

		onceStore.Do(func() {
			c := &ristretto.Config{
				NumCounters: 1e7,     // number of keys to track frequency of (10M).
				MaxCost:     1 << 30, // maximum cost of cache (1GB).
				BufferItems: 64,      // number of keys per Get buffer.
				Cost:        nil,
			}

			secretCache, err = ristretto.NewCache(c)
			if err != nil {
				return
			}
			policyCache, err = ristretto.NewCache(c)
			if err != nil {
				return
			}

			storeIns = &Store{
				cli:      cli,
				lock:     new(sync.RWMutex),
				secrets:  secretCache,
				policies: policyCache,
			}
		})
	}

	return storeIns, err
}

// GetSecret return secret detail for the given key.
func (s *Store) GetSecret(key string) (*pb.SecretInfo, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	value, ok := s.secrets.Get(key)
	if !ok {
		return nil, ErrSecretNotFound
	}

	return value.(*pb.SecretInfo), nil
}

// GetPolicy return user's ladon policies for the given user.
func (s *Store) GetPolicy(key string) ([]*ladon.DefaultPolicy, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	value, ok := s.policies.Get(key)
	if !ok {
		return nil, ErrPolicyNotFound
	}

	return value.([]*ladon.DefaultPolicy), nil
}

// Reload reload secrets and policies.
func (s *Store) Reload() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	// reload secrets
	secrets, err := s.cli.GetSecrets()
	if err != nil {
		return err
	}

	s.secrets.Clear()
	for key, val := range secrets {
		s.secrets.Set(key, val, 1)
	}

	// reload policies
	policies, err := s.cli.GetPolicies()
	if err != nil {
		return err
	}

	s.policies.Clear()
	for key, val := range policies {
		s.policies.Set(key, val, 1)
	}

	return nil
}
