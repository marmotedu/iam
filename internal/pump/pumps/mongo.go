// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pumps

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/marmotedu/errors"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2"

	"github.com/marmotedu/iam/internal/pump/analytics"
	"github.com/marmotedu/iam/pkg/log"
)

// Define unit constant.
const (
	_   = iota // ignore zero iota
	KiB = 1 << (10 * iota)
	MiB
	GiB
	TiB
)

// MongoPump defines a mongo pump with mongo specific options and common options.
type MongoPump struct {
	dbSession *mgo.Session
	dbConf    *MongoConf
	CommonPumpConfig
}

var mongoPumpPrefix = "PMP_MONGO"

// MongoType define a new mongo type.
type MongoType int

// Defines mongo type.
const (
	StandardMongo MongoType = iota
	AWSDocumentDB
)

// BaseMongoConf defines options needed when connnect to mongo db.
type BaseMongoConf struct {
	MongoURL                      string    `json:"mongo_url"                         mapstructure:"mongo_url"`
	MongoUseSSL                   bool      `json:"mongo_use_ssl"                     mapstructure:"mongo_use_ssl"`
	MongoSSLInsecureSkipVerify    bool      `json:"mongo_ssl_insecure_skip_verify"    mapstructure:"mongo_ssl_insecure_skip_verify"`
	MongoSSLAllowInvalidHostnames bool      `json:"mongo_ssl_allow_invalid_hostnames" mapstructure:"mongo_ssl_allow_invalid_hostnames"`
	MongoSSLCAFile                string    `json:"mongo_ssl_ca_file"                 mapstructure:"mongo_ssl_ca_file"`
	MongoSSLPEMKeyfile            string    `json:"mongo_ssl_pem_keyfile"             mapstructure:"mongo_ssl_pem_keyfile"`
	MongoDBType                   MongoType `json:"mongo_db_type"                     mapstructure:"mongo_db_type"`
}

// MongoConf defines mongo specific options.
type MongoConf struct {
	BaseMongoConf

	CollectionName            string `json:"collection_name"               mapstructure:"collection_name"`
	MaxInsertBatchSizeBytes   int    `json:"max_insert_batch_size_bytes"   mapstructure:"max_insert_batch_size_bytes"`
	MaxDocumentSizeBytes      int    `json:"max_document_size_bytes"       mapstructure:"max_document_size_bytes"`
	CollectionCapMaxSizeBytes int    `json:"collection_cap_max_size_bytes" mapstructure:"collection_cap_max_size_bytes"`
	CollectionCapEnable       bool   `json:"collection_cap_enable"         mapstructure:"collection_cap_enable"`
}

func loadCertificateAndKeyFromFile(path string) (*tls.Certificate, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	var cert tls.Certificate
	for {
		block, rest := pem.Decode(raw)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, block.Bytes)
		} else {
			cert.PrivateKey, err = parsePrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failure reading private key from \"%s\": %w", path, err)
			}
		}
		raw = rest
	}

	if len(cert.Certificate) == 0 {
		return nil, fmt.Errorf("no certificate found in \"%s\"", path)
	} else if cert.PrivateKey == nil {
		return nil, fmt.Errorf("no private key found in \"%s\"", path)
	}

	return &cert, nil
}

func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, fmt.Errorf("found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("failed to parse private key")
}

func mongoType(session *mgo.Session) MongoType {
	// Querying for the features which 100% not supported by AWS DocumentDB
	var result struct {
		Code int `bson:"code"`
	}

	_ = session.Run("features", &result)

	if result.Code == 303 {
		return AWSDocumentDB
	}

	return StandardMongo
}

// nolint: gocognit
func mongoDialInfo(conf BaseMongoConf) (dialInfo *mgo.DialInfo, err error) {
	if dialInfo, err = mgo.ParseURL(conf.MongoURL); err != nil {
		return dialInfo, errors.Wrap(err, "failed to parse mongo url")
	}

	// nolint: nestif
	if conf.MongoUseSSL {
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			tlsConfig := &tls.Config{}
			if conf.MongoSSLInsecureSkipVerify {
				tlsConfig.InsecureSkipVerify = true
			}

			if conf.MongoSSLCAFile != "" {
				var caCert []byte
				caCert, err = ioutil.ReadFile(conf.MongoSSLCAFile)
				if err != nil {
					log.Fatalf("Can't load mongo CA certificates: %s", err.Error())
				}
				caCertPool := x509.NewCertPool()
				caCertPool.AppendCertsFromPEM(caCert)
				tlsConfig.RootCAs = caCertPool
			}

			if conf.MongoSSLAllowInvalidHostnames {
				tlsConfig.InsecureSkipVerify = true
				tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					// Code copy/pasted and adapted from
					// https://github.com/golang/go/blob/81555cb4f3521b53f9de4ce15f64b77cc9df61b9/src/crypto/tls/handshake_client.go#L327-L344,
					// but adapted to skip the hostname verification.
					// See https://github.com/golang/go/issues/21971#issuecomment-412836078.

					// If this is the first handshake on a connection, process and
					// (optionally) verify the server's certificates.
					certs := make([]*x509.Certificate, len(rawCerts))
					for i, asn1Data := range rawCerts {
						var cert *x509.Certificate
						cert, err = x509.ParseCertificate(asn1Data)
						if err != nil {
							return errors.Wrap(err, "failed to parse certificate")
						}
						certs[i] = cert
					}

					opts := x509.VerifyOptions{
						Roots:         tlsConfig.RootCAs,
						CurrentTime:   time.Now(),
						DNSName:       "", // <- skip hostname verification
						Intermediates: x509.NewCertPool(),
					}

					for i, cert := range certs {
						if i == 0 {
							continue
						}
						opts.Intermediates.AddCert(cert)
					}
					_, err = certs[0].Verify(opts)

					return errors.Wrap(err, "failed to verify certificate")
				}
			}

			if conf.MongoSSLPEMKeyfile != "" {
				var cert *tls.Certificate
				cert, err = loadCertificateAndKeyFromFile(conf.MongoSSLPEMKeyfile)
				if err != nil {
					log.Fatalf("Can't load mongo client certificate: %s", err.Error())
				}

				tlsConfig.Certificates = []tls.Certificate{*cert}
			}

			return tls.Dial("tcp", addr.String(), tlsConfig)
		}
	}

	return dialInfo, err
}

// New create a mongo pump instance.
func (m *MongoPump) New() Pump {
	newPump := MongoPump{}

	return &newPump
}

// GetName returns the mongo pump name.
func (m *MongoPump) GetName() string {
	return "MongoDB Pump"
}

// Init initialize the mongo pump instance.
func (m *MongoPump) Init(config interface{}) error {
	m.dbConf = &MongoConf{}
	err := mapstructure.Decode(config, &m.dbConf)
	if err == nil {
		err = mapstructure.Decode(config, &m.dbConf.BaseMongoConf)
		log.Info("Init", log.String("url", m.dbConf.MongoURL), log.String("collection_name", m.dbConf.CollectionName))
		if err != nil {
			panic(m.dbConf.BaseMongoConf)
		}
	}

	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	overrideErr := envconfig.Process(mongoPumpPrefix, m.dbConf)
	if overrideErr != nil {
		log.Errorf("Failed to process environment variables for mongo pump: %s", overrideErr.Error())
	}

	if m.dbConf.MaxInsertBatchSizeBytes == 0 {
		log.Info("-- No max batch size set, defaulting to 10MB")
		m.dbConf.MaxInsertBatchSizeBytes = 10 * MiB
	}

	if m.dbConf.MaxDocumentSizeBytes == 0 {
		log.Info("-- No max document size set, defaulting to 10MB")
		m.dbConf.MaxDocumentSizeBytes = 10 * MiB
	}

	m.connect()

	m.capCollection()

	indexCreateErr := m.ensureIndexes()
	if indexCreateErr != nil {
		log.Error(indexCreateErr.Error())
	}

	log.Debugf("MongoDB DB CS: %s", m.dbConf.MongoURL)
	log.Debugf("MongoDB Col: %s", m.dbConf.CollectionName)

	return nil
}

func (m *MongoPump) capCollection() (ok bool) {
	colName := m.dbConf.CollectionName
	colCapMaxSizeBytes := m.dbConf.CollectionCapMaxSizeBytes
	colCapEnable := m.dbConf.CollectionCapEnable

	if !colCapEnable {
		return false
	}

	exists, err := m.collectionExists(colName)
	if err != nil {
		log.Errorf("Unable to determine if collection (%s) exists. Not capping collection: %s", colName, err.Error())

		return false
	}

	if exists {
		log.Warnf("Collection (%s) already exists. Capping could result in data loss. Ignoring", colName)

		return false
	}

	if strconv.IntSize < 64 {
		log.Warn("Pump running < 64bit architecture. Not capping collection as max size would be 2gb")

		return false
	}

	if colCapMaxSizeBytes == 0 {
		defaultBytes := 5
		colCapMaxSizeBytes = defaultBytes * GiB

		log.Infof("-- No max collection size set for %s, defaulting to %d", colName, colCapMaxSizeBytes)
	}

	sess := m.dbSession.Copy()
	defer sess.Close()

	err = m.dbSession.DB("").C(colName).Create(&mgo.CollectionInfo{Capped: true, MaxBytes: colCapMaxSizeBytes})
	if err != nil {
		log.Errorf("Unable to create capped collection for (%s). %s", colName, err.Error())

		return false
	}

	log.Infof("Capped collection (%s) created. %d bytes", colName, colCapMaxSizeBytes)

	return true
}

// collectionExists checks to see if a collection name exists in the db.
func (m *MongoPump) collectionExists(name string) (bool, error) {
	sess := m.dbSession.Copy()
	defer sess.Close()

	colNames, err := sess.DB("").CollectionNames()
	if err != nil {
		log.Errorf("Unable to get column names: %s", err.Error())

		return false, errors.Wrap(err, "failed to get collection names")
	}

	for _, coll := range colNames {
		if coll == name {
			return true, nil
		}
	}

	return false, nil
}

func (m *MongoPump) ensureIndexes() error {
	var err error

	sess := m.dbSession.Copy()
	defer sess.Close()

	c := sess.DB("").C(m.dbConf.CollectionName)

	orgIndex := mgo.Index{
		Key:        []string{"orgid"},
		Background: m.dbConf.MongoDBType == StandardMongo,
	}

	err = c.EnsureIndex(orgIndex)
	if err != nil {
		return errors.Wrap(err, "failed to ensures an index with the given key exists")
	}

	apiIndex := mgo.Index{
		Key:        []string{"apiid"},
		Background: m.dbConf.MongoDBType == StandardMongo,
	}

	err = c.EnsureIndex(apiIndex)
	if err != nil {
		return errors.Wrap(err, "failed to ensures an index with the given key exists")
	}

	logBrowserIndex := mgo.Index{
		Name:       "logBrowserIndex",
		Key:        []string{"-timestamp", "orgid", "apiid", "apikey", "responsecode"},
		Background: m.dbConf.MongoDBType == StandardMongo,
	}

	err = c.EnsureIndex(logBrowserIndex)
	if err != nil && !strings.Contains(err.Error(), "already exists with a different name") {
		return errors.Wrap(err, "failed to ensures an index with the given key exists")
	}

	return nil
}

func (m *MongoPump) connect() {
	var err error
	var dialInfo *mgo.DialInfo

	dialInfo, err = mongoDialInfo(m.dbConf.BaseMongoConf)
	if err != nil {
		log.Panicf("Mongo URL is invalid: %s", err.Error())
	}

	dialInfo.Timeout = time.Second * 5
	m.dbSession, err = mgo.DialWithInfo(dialInfo)

	for err != nil {
		log.Error("Mongo connection failed. Retrying.", log.String("error", err.Error()))
		time.Sleep(5 * time.Second)
		m.dbSession, err = mgo.DialWithInfo(dialInfo)
	}

	if err == nil && m.dbConf.MongoDBType == 0 {
		m.dbConf.MongoDBType = mongoType(m.dbSession)
	}
}

// WriteData write analyzed data to mongo persistent back-end storage.
func (m *MongoPump) WriteData(ctx context.Context, data []interface{}) error {
	collectionName := m.dbConf.CollectionName
	if collectionName == "" {
		log.Fatal("No collection name!")
	}

	log.Debugf("Writing %d records", len(data))

	for m.dbSession == nil {
		log.Debug("Connecting to analytics store")
		m.connect()
	}

	for _, dataSet := range m.AccumulateSet(data) {
		go func(dataSet []interface{}) {
			sess := m.dbSession.Copy()
			defer sess.Close()

			analyticsCollection := sess.DB("").C(collectionName)

			log.Infof("Purging %d records", len(dataSet))

			err := analyticsCollection.Insert(dataSet...)
			if err != nil {
				log.Errorf("Problem inserting to mongo collection: %s", err.Error())
				if strings.Contains(strings.ToLower(err.Error()), "closed explicitly") {
					log.Warn("--> Detected connection failure!")
				}
			}
		}(dataSet)
	}

	return nil
}

// AccumulateSet accumulate data.
func (m *MongoPump) AccumulateSet(data []interface{}) [][]interface{} {
	accumulatorTotal := 0
	returnArray := make([][]interface{}, 0)
	thisResultSet := make([]interface{}, 0)

	for i, item := range data {
		thisItem, _ := item.(analytics.AnalyticsRecord)

		// Add 1 KB for metadata as average
		sizeBytes := len(thisItem.Policies) + len(thisItem.Deciders) + 1024

		log.Debugf("Size is: %d", sizeBytes)

		if sizeBytes > m.dbConf.MaxDocumentSizeBytes {
			log.Warn("Document too large, not writing raw request and raw response!")

			thisItem.Policies = ""
			thisItem.Deciders = ""
		}

		if (accumulatorTotal + sizeBytes) <= m.dbConf.MaxInsertBatchSizeBytes {
			accumulatorTotal += sizeBytes
		} else {
			log.Debug("Created new chunk entry")
			if len(thisResultSet) > 0 {
				returnArray = append(returnArray, thisResultSet)
			}

			thisResultSet = make([]interface{}, 0)
			accumulatorTotal = sizeBytes
		}

		log.Debugf("Accumulator is: %d", accumulatorTotal)
		thisResultSet = append(thisResultSet, thisItem)

		log.Debugf("%d of %d", accumulatorTotal, m.dbConf.MaxInsertBatchSizeBytes)
		// Append the last element if the loop is about to end
		if i == (len(data) - 1) {
			log.Debug("Appending last entry")
			returnArray = append(returnArray, thisResultSet)
		}
	}

	return returnArray
}
