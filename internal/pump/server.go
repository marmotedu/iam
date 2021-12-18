// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pump

import (
	"context"
	"fmt"
	"sync"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/marmotedu/iam/internal/pump/analytics"
	"github.com/marmotedu/iam/internal/pump/config"
	"github.com/marmotedu/iam/internal/pump/options"
	"github.com/marmotedu/iam/internal/pump/pumps"
	"github.com/marmotedu/iam/internal/pump/storage"
	"github.com/marmotedu/iam/internal/pump/storage/redis"
	"github.com/marmotedu/iam/pkg/log"
)

var pmps []pumps.Pump

type pumpServer struct {
	secInterval    int
	omitDetails    bool
	mutex          *redsync.Mutex
	analyticsStore storage.AnalyticsStorage
	pumps          map[string]options.PumpConfig
}

// preparedGenericAPIServer is a private wrapper that enforces a call of PrepareRun() before Run can be invoked.
type preparedPumpServer struct {
	*pumpServer
}

func createPumpServer(cfg *config.Config) (*pumpServer, error) {
	// use the same redis database with authorization log history
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisOptions.Host, cfg.RedisOptions.Port),
		Username: cfg.RedisOptions.Username,
		Password: cfg.RedisOptions.Password,
	})

	rs := redsync.New(goredis.NewPool(client))

	server := &pumpServer{
		secInterval:    cfg.PurgeDelay,
		omitDetails:    cfg.OmitDetailedRecording,
		mutex:          rs.NewMutex("iam-pump", redsync.WithExpiry(10*time.Minute)),
		analyticsStore: &redis.RedisClusterStorageManager{},
		pumps:          cfg.Pumps,
	}

	if err := server.analyticsStore.Init(cfg.RedisOptions); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *pumpServer) PrepareRun() preparedPumpServer {
	s.initialize()

	return preparedPumpServer{s}
}

func (s preparedPumpServer) Run(stopCh <-chan struct{}) error {
	ticker := time.NewTicker(time.Duration(s.secInterval) * time.Second)
	defer ticker.Stop()

	log.Info("Now run loop to clean data from redis")
	for {
		select {
		case <-ticker.C:
			s.pump()
		// exit consumption cycle when receive SIGINT and SIGTERM signal
		case <-stopCh:
			log.Info("stop purge loop")

			return nil
		}
	}
}

// pump get authorization log from redis and write to pumps.
func (s *pumpServer) pump() {
	if err := s.mutex.Lock(); err != nil {
		log.Info("there is already an iam-pump instance running.")

		return
	}
	defer func() {
		if _, err := s.mutex.Unlock(); err != nil {
			log.Errorf("could not release iam-pump lock. err: %v", err)
		}
	}()

	analyticsValues := s.analyticsStore.GetAndDeleteSet(storage.AnalyticsKeyName)
	if len(analyticsValues) == 0 {
		return
	}

	// Convert to something clean
	keys := make([]interface{}, len(analyticsValues))

	for i, v := range analyticsValues {
		decoded := analytics.AnalyticsRecord{}
		err := msgpack.Unmarshal([]byte(v.(string)), &decoded)
		log.Debugf("Decoded Record: %v", decoded)
		if err != nil {
			log.Errorf("Couldn't unmarshal analytics data: %s", err.Error())
		} else {
			if s.omitDetails {
				decoded.Policies = ""
				decoded.Deciders = ""
			}
			keys[i] = interface{}(decoded)
		}
	}

	// Send to pumps
	writeToPumps(keys, s.secInterval)
}

func (s *pumpServer) initialize() {
	pmps = make([]pumps.Pump, len(s.pumps))
	i := 0
	for key, pmp := range s.pumps {
		pumpTypeName := pmp.Type
		if pumpTypeName == "" {
			pumpTypeName = key
		}

		pmpType, err := pumps.GetPumpByName(pumpTypeName)
		if err != nil {
			log.Errorf("Pump load error (skipping): %s", err.Error())
		} else {
			pmpIns := pmpType.New()
			initErr := pmpIns.Init(pmp.Meta)
			if initErr != nil {
				log.Errorf("Pump init error (skipping): %s", initErr.Error())
			} else {
				log.Infof("Init Pump: %s", pmpIns.GetName())
				pmpIns.SetFilters(pmp.Filters)
				pmpIns.SetTimeout(pmp.Timeout)
				pmpIns.SetOmitDetailedRecording(pmp.OmitDetailedRecording)
				pmps[i] = pmpIns
			}
		}
		i++
	}
}

func writeToPumps(keys []interface{}, purgeDelay int) {
	// Send to pumps
	if pmps != nil {
		var wg sync.WaitGroup
		wg.Add(len(pmps))
		for _, pmp := range pmps {
			go execPumpWriting(&wg, pmp, &keys, purgeDelay)
		}
		wg.Wait()
	} else {
		log.Warn("No pumps defined!")
	}
}

func filterData(pump pumps.Pump, keys []interface{}) []interface{} {
	filters := pump.GetFilters()
	if !filters.HasFilter() && !pump.GetOmitDetailedRecording() {
		return keys
	}
	filteredKeys := keys[:] // nolint: gocritic
	newLenght := 0

	for _, key := range filteredKeys {
		decoded, _ := key.(analytics.AnalyticsRecord)
		if pump.GetOmitDetailedRecording() {
			decoded.Policies = ""
			decoded.Deciders = ""
		}
		if filters.ShouldFilter(decoded) {
			continue
		}
		filteredKeys[newLenght] = decoded
		newLenght++
	}
	filteredKeys = filteredKeys[:newLenght]

	return filteredKeys
}

func execPumpWriting(wg *sync.WaitGroup, pmp pumps.Pump, keys *[]interface{}, purgeDelay int) {
	timer := time.AfterFunc(time.Duration(purgeDelay)*time.Second, func() {
		if pmp.GetTimeout() == 0 {
			log.Warnf(
				"Pump %s is taking more time than the value configured of purge_delay. You should try to set a timeout for this pump.",
				pmp.GetName(),
			)
		} else if pmp.GetTimeout() > purgeDelay {
			log.Warnf("Pump %s is taking more time than the value configured of purge_delay. You should try lowering the timeout configured for this pump.", pmp.GetName())
		}
	})
	defer timer.Stop()
	defer wg.Done()

	log.Debugf("Writing to: %s", pmp.GetName())

	ch := make(chan error, 1)
	var ctx context.Context
	var cancel context.CancelFunc
	// Initialize context depending if the pump has a configured timeout
	if tm := pmp.GetTimeout(); tm > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(tm)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}

	defer cancel()

	go func(ch chan error, ctx context.Context, pmp pumps.Pump, keys *[]interface{}) {
		filteredKeys := filterData(pmp, *keys)

		ch <- pmp.WriteData(ctx, filteredKeys)
	}(ch, ctx, pmp, keys)

	select {
	case err := <-ch:
		if err != nil {
			log.Warnf("Error Writing to: %s - Error: %s", pmp.GetName(), err.Error())
		}
	case <-ctx.Done():
		//nolint: errorlint
		switch ctx.Err() {
		case context.Canceled:
			log.Warnf("The writing to %s have got canceled.", pmp.GetName())
		case context.DeadlineExceeded:
			log.Warnf("Timeout Writing to: %s", pmp.GetName())
		}
	}
}
