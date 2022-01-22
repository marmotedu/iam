// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pumps

import (
	"context"
	"encoding/base64"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/marmotedu/errors"
	"github.com/mitchellh/mapstructure"
	elastic "github.com/olivere/elastic/v7"

	"github.com/marmotedu/iam/internal/pump/analytics"
	"github.com/marmotedu/iam/pkg/log"
)

// ElasticsearchPump defines an elasticsearch pump with elasticsearch specific options and common options.
type ElasticsearchPump struct {
	operator ElasticsearchOperator
	esConf   *ElasticsearchConf
	CommonPumpConfig
}

// ElasticsearchConf defines elasticsearch specific options.
type ElasticsearchConf struct {
	BulkConfig       ElasticsearchBulkConfig `mapstructure:"bulk_config"`
	IndexName        string                  `mapstructure:"index_name"`
	ElasticsearchURL string                  `mapstructure:"elasticsearch_url"`
	DocumentType     string                  `mapstructure:"document_type"`
	AuthAPIKeyID     string                  `mapstructure:"auth_api_key_id"`
	AuthAPIKey       string                  `mapstructure:"auth_api_key"`
	Username         string                  `mapstructure:"auth_basic_username"`
	Password         string                  `mapstructure:"auth_basic_password"`
	EnableSniffing   bool                    `mapstructure:"use_sniffing"`
	RollingIndex     bool                    `mapstructure:"rolling_index"`
	DisableBulk      bool                    `mapstructure:"disable_bulk"`
}

// ElasticsearchBulkConfig defines elasticsearch bulk config.
type ElasticsearchBulkConfig struct {
	Workers       int `mapstructure:"workers"`
	FlushInterval int `mapstructure:"flush_interval"`
	BulkActions   int `mapstructure:"bulk_actions"`
	BulkSize      int `mapstructure:"bulk_size"`
}

// ElasticsearchOperator defines interface for all elasticsearch operator.
type ElasticsearchOperator interface {
	processData(ctx context.Context, data []interface{}, esConf *ElasticsearchConf) error
}

// Elasticsearch7Operator defines elasticsearch6 operator.
type Elasticsearch7Operator struct {
	esClient      *elastic.Client
	bulkProcessor *elastic.BulkProcessor
}

// APIKeyTransport defiens elasticsearch api key.
type APIKeyTransport struct {
	APIKey   string
	APIKeyID string
}

// RoundTrip for APIKeyTransport auth.
func (t *APIKeyTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	auth := t.APIKeyID + ":" + t.APIKey
	key := base64.StdEncoding.EncodeToString([]byte(auth))

	r.Header.Set("Authorization", "ApiKey "+key)

	return http.DefaultTransport.RoundTrip(r)
}

func getOperator(ctx context.Context, conf ElasticsearchConf) (ElasticsearchOperator, error) {
	var err error
	urls := strings.Split(conf.ElasticsearchURL, ",")
	httpClient := http.DefaultClient
	if conf.AuthAPIKey != "" && conf.AuthAPIKeyID != "" {
		conf.Username = ""
		conf.Password = ""
		httpClient = &http.Client{Transport: &APIKeyTransport{APIKey: conf.AuthAPIKey, APIKeyID: conf.AuthAPIKeyID}}
	}

	e := new(Elasticsearch7Operator)

	e.esClient, err = elastic.NewClient(
		elastic.SetURL(urls...),
		elastic.SetSniff(conf.EnableSniffing),
		elastic.SetBasicAuth(conf.Username, conf.Password),
		elastic.SetHttpClient(httpClient),
	)

	if err != nil {
		return e, errors.Wrap(err, "failed to new es client")
	}
	// Setup a bulk processor
	p := e.esClient.BulkProcessor().Name("IAMPumpESv6BackgroundProcessor")
	if conf.BulkConfig.Workers != 0 {
		p = p.Workers(conf.BulkConfig.Workers)
	}

	if conf.BulkConfig.FlushInterval != 0 {
		p = p.FlushInterval(time.Duration(conf.BulkConfig.FlushInterval) * time.Second)
	}

	if conf.BulkConfig.BulkActions != 0 {
		p = p.BulkActions(conf.BulkConfig.BulkActions)
	}

	if conf.BulkConfig.BulkSize != 0 {
		p = p.BulkSize(conf.BulkConfig.BulkSize)
	}

	e.bulkProcessor, err = p.Do(ctx)

	return e, errors.Wrap(err, "failed to start bulk processor")
}

// New create an elasticsearch pump instance.
func (e *ElasticsearchPump) New() Pump {
	newPump := ElasticsearchPump{}

	return &newPump
}

// GetName returns the elasticsearch pump name.
func (e *ElasticsearchPump) GetName() string {
	return "Elasticsearch Pump"
}

// Init initialize the elasticsearch pump instance.
func (e *ElasticsearchPump) Init(config interface{}) error {
	e.esConf = &ElasticsearchConf{}
	loadConfigErr := mapstructure.Decode(config, &e.esConf)

	if loadConfigErr != nil {
		log.Fatalf("Failed to decode configuration: %s", loadConfigErr.Error())
	}

	if e.esConf.IndexName == "" {
		e.esConf.IndexName = "iam_analytics"
	}

	if e.esConf.ElasticsearchURL == "" {
		e.esConf.ElasticsearchURL = "http://localhost:9200"
	}

	if e.esConf.DocumentType == "" {
		e.esConf.DocumentType = "iam_analytics"
	}

	re := regexp.MustCompile(`(.*)\/\/(.*):(.*)\@(.*)`)
	printableURL := re.ReplaceAllString(e.esConf.ElasticsearchURL, `$1//***:***@$4`)

	log.Infof("Elasticsearch URL: %s", printableURL)
	log.Infof("Elasticsearch Index: %s", e.esConf.IndexName)
	if e.esConf.RollingIndex {
		log.Infof("Index will have date appended to it in the format %s -YYYY.MM.DD", e.esConf.IndexName)
	}

	e.connect(context.Background())

	return nil
}

func (e *ElasticsearchPump) connect(ctx context.Context) {
	var err error

	e.operator, err = getOperator(ctx, *e.esConf)
	if err != nil {
		log.Errorf("Elasticsearch connection failed: %s", err.Error())
		time.Sleep(5 * time.Second)
		e.connect(ctx)
	}
}

// WriteData write analyzed data to elasticsearch persistent back-end storage.
func (e *ElasticsearchPump) WriteData(ctx context.Context, data []interface{}) error {
	log.Infof("Writing %d records", len(data))

	if e.operator == nil {
		log.Debug("Connecting to analytics store")
		e.connect(ctx)
		_ = e.WriteData(ctx, data)
	} else if len(data) > 0 {
		_ = e.operator.processData(ctx, data, e.esConf)
	}

	return nil
}

func getIndexName(esConf *ElasticsearchConf) string {
	indexName := esConf.IndexName

	if esConf.RollingIndex {
		currentTime := time.Now()
		// This formats the date to be YYYY.MM.DD but Golang makes you use a specific date for its date formatting
		indexName += "-" + currentTime.Format("2006.01.02")
	}

	return indexName
}

func getMapping(datum analytics.AnalyticsRecord) (map[string]interface{}, string) {
	record := datum
	mapping := map[string]interface{}{
		"@timestamp": record.TimeStamp,
		"username":   record.Username,
		"effect":     record.Effect,
		"conclusion": record.Conclusion,
		"request":    record.Request,
		"policies":   record.Policies,
		"deciders":   record.Deciders,
		"expireAt":   record.ExpireAt,
	}

	return mapping, ""
}

func (e Elasticsearch7Operator) processData(ctx context.Context, data []interface{}, esConf *ElasticsearchConf) error {
	index := e.esClient.Index().Index(getIndexName(esConf))

	for dataIndex := range data {
		if ctxErr := ctx.Err(); ctxErr != nil {
			continue
		}

		d, ok := data[dataIndex].(analytics.AnalyticsRecord)
		if !ok {
			log.Errorf("Error while writing %s: data not of type analytics.AnalyticsRecord", data[dataIndex])

			continue
		}

		mapping, id := getMapping(d)

		if !esConf.DisableBulk {
			r := elastic.NewBulkIndexRequest().Index(getIndexName(esConf)).Type(esConf.DocumentType).Id(id).Doc(mapping)
			e.bulkProcessor.Add(r)
		} else {
			//nolint: staticcheck
			_, err := index.BodyJson(mapping).Type(esConf.DocumentType).Id(id).Do(ctx)
			if err != nil {
				log.Errorf("Error while writing %s %s", data[dataIndex], err.Error())
			}
		}
	}

	return nil
}
