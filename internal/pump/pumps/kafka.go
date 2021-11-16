// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pumps

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/marmotedu/component-base/pkg/json"
	"github.com/marmotedu/errors"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/segmentio/kafka-go/snappy"

	"github.com/marmotedu/iam/internal/pump/analytics"
	"github.com/marmotedu/iam/pkg/log"
)

// KafkaPump defines a kafka pump with kafka specific options and common options.
type KafkaPump struct {
	kafkaConf    *KafkaConf
	writerConfig kafka.WriterConfig
	CommonPumpConfig
}

// Message contains the messages need to push to pump.
type Message map[string]interface{}

// KafkaConf defines kafka specific options.
type KafkaConf struct {
	Broker                []string          `mapstructure:"broker"`
	ClientID              string            `mapstructure:"client_id"`
	Topic                 string            `mapstructure:"topic"`
	SSLCertFile           string            `mapstructure:"ssl_cert_file"`
	SSLKeyFile            string            `mapstructure:"ssl_key_file"`
	SASLMechanism         string            `mapstructure:"sasl_mechanism"`
	Username              string            `mapstructure:"sasl_username"`
	Password              string            `mapstructure:"sasl_password"`
	Algorithm             string            `mapstructure:"sasl_algorithm"`
	Timeout               time.Duration     `mapstructure:"timeout"`
	MetaData              map[string]string `mapstructure:"meta_data"`
	Compressed            bool              `mapstructure:"compressed"`
	UseSSL                bool              `mapstructure:"use_ssl"`
	SSLInsecureSkipVerify bool              `mapstructure:"ssl_insecure_skip_verify"`
}

// New create a kafka pump instance.
func (k *KafkaPump) New() Pump {
	newPump := KafkaPump{}

	return &newPump
}

// GetName returns the kafka pump name.
func (k *KafkaPump) GetName() string {
	return "Kafka Pump"
}

// Init initialize the kafka pump instance.
func (k *KafkaPump) Init(config interface{}) error {
	// Read configuration file
	k.kafkaConf = &KafkaConf{}
	err := mapstructure.Decode(config, &k.kafkaConf)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	var tlsConfig *tls.Config
	// nolint: nestif
	if k.kafkaConf.UseSSL {
		if k.kafkaConf.SSLCertFile != "" && k.kafkaConf.SSLKeyFile != "" {
			var cert tls.Certificate
			log.Debug("Loading certificates for mTLS.")
			cert, err = tls.LoadX509KeyPair(k.kafkaConf.SSLCertFile, k.kafkaConf.SSLKeyFile)
			if err != nil {
				log.Debugf("Error loading mTLS certificates: %s", err.Error())

				return errors.Wrap(err, "failed loading mTLS certificates")
			}
			tlsConfig = &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: k.kafkaConf.SSLInsecureSkipVerify,
			}
		} else if k.kafkaConf.SSLCertFile != "" || k.kafkaConf.SSLKeyFile != "" {
			log.Error("Only one of ssl_cert_file and ssl_cert_key configuration option is setted, you should set both to enable mTLS.")
		} else {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: k.kafkaConf.SSLInsecureSkipVerify,
			}
		}
	} else if k.kafkaConf.SASLMechanism != "" {
		log.Warn("SASL-Mechanism is setted but use_ssl is false.", log.String("SASL-Mechanism", k.kafkaConf.SASLMechanism))
	}

	var mechanism sasl.Mechanism

	switch k.kafkaConf.SASLMechanism {
	case "":
		break
	case "PLAIN", "plain":
		mechanism = plain.Mechanism{Username: k.kafkaConf.Username, Password: k.kafkaConf.Password}
	case "SCRAM", "scram":
		algorithm := scram.SHA256
		if k.kafkaConf.Algorithm == "sha-512" || k.kafkaConf.Algorithm == "SHA-512" {
			algorithm = scram.SHA512
		}
		var mechErr error
		mechanism, mechErr = scram.Mechanism(algorithm, k.kafkaConf.Username, k.kafkaConf.Password)
		if mechErr != nil {
			log.Fatalf("Failed initialize kafka mechanism: %s", mechErr.Error())
		}
	default:
		log.Warn(
			"IAM pump doesn't support this SASL mechanism.",
			log.String("SASL-Mechanism", k.kafkaConf.SASLMechanism),
		)
	}

	// Kafka writer connection config
	dialer := &kafka.Dialer{
		Timeout:       k.kafkaConf.Timeout,
		ClientID:      k.kafkaConf.ClientID,
		TLS:           tlsConfig,
		SASLMechanism: mechanism,
	}

	k.writerConfig.Brokers = k.kafkaConf.Broker
	k.writerConfig.Topic = k.kafkaConf.Topic
	k.writerConfig.Balancer = &kafka.LeastBytes{}
	k.writerConfig.Dialer = dialer
	k.writerConfig.WriteTimeout = k.kafkaConf.Timeout
	k.writerConfig.ReadTimeout = k.kafkaConf.Timeout
	if k.kafkaConf.Compressed {
		k.writerConfig.CompressionCodec = snappy.NewCompressionCodec()
	}

	log.Infof("Kafka config: %s", k.writerConfig)

	return nil
}

// WriteData write analyzed data to kafka persistent back-end storage.
func (k *KafkaPump) WriteData(ctx context.Context, data []interface{}) error {
	startTime := time.Now()
	log.Infof("Writing %d records ...", len(data))
	kafkaMessages := make([]kafka.Message, len(data))
	for i, v := range data {
		// Build message format
		decoded, _ := v.(analytics.AnalyticsRecord)
		message := Message{
			"timestamp":  decoded.TimeStamp,
			"username":   decoded.Username,
			"effect":     decoded.Effect,
			"conclusion": decoded.Conclusion,
			"request":    decoded.Request,
			"policies":   decoded.Policies,
			"deciders":   decoded.Deciders,
			"expireAt":   decoded.ExpireAt,
		}
		// Add static metadata to json
		for key, value := range k.kafkaConf.MetaData {
			message[key] = value
		}

		// Transform object to json string
		json, jsonError := json.Marshal(message)
		if jsonError != nil {
			log.Error("unable to marshal message", log.String("error", jsonError.Error()))
		}

		// Kafka message structure
		kafkaMessages[i] = kafka.Message{
			Time:  time.Now(),
			Value: json,
		}
	}
	// Send kafka message
	kafkaError := k.write(ctx, kafkaMessages)
	if kafkaError != nil {
		log.Error("unable to write message", log.String("error", kafkaError.Error()))
	}
	log.Debugf("ElapsedTime in seconds for %d records %v", len(data), time.Since(startTime))

	return nil
}

func (k *KafkaPump) write(ctx context.Context, messages []kafka.Message) error {
	kafkaWriter := kafka.NewWriter(k.writerConfig)
	defer kafkaWriter.Close()

	return kafkaWriter.WriteMessages(ctx, messages...)
}
