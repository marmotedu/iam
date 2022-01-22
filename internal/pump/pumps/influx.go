// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pumps

import (
	"context"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/marmotedu/component-base/pkg/json"
	"github.com/mitchellh/mapstructure"

	"github.com/marmotedu/iam/internal/pump/analytics"
	"github.com/marmotedu/iam/pkg/log"
)

// InfluxPump defines an influx pump with influx specific options and common options.
type InfluxPump struct {
	dbConf *InfluxConf
	CommonPumpConfig
}

var table = "analytics"

// InfluxConf defines influx specific options.
type InfluxConf struct {
	DatabaseName string   `mapstructure:"database_name"`
	Addr         string   `mapstructure:"address"`
	Username     string   `mapstructure:"username"`
	Password     string   `mapstructure:"password"`
	Fields       []string `mapstructure:"fields"`
	Tags         []string `mapstructure:"tags"`
}

// New create an influx pump instance.
func (i *InfluxPump) New() Pump {
	newPump := InfluxPump{}

	return &newPump
}

// GetName returns the influx pump name.
func (i *InfluxPump) GetName() string {
	return "InfluxDB Pump"
}

// Init initialize the influx pump instance.
func (i *InfluxPump) Init(config interface{}) error {
	i.dbConf = &InfluxConf{}
	err := mapstructure.Decode(config, &i.dbConf)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	i.connect()

	log.Debugf("Influx DB CS: %s", i.dbConf.Addr)

	return nil
}

func (i *InfluxPump) connect() client.Client {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     i.dbConf.Addr,
		Username: i.dbConf.Username,
		Password: i.dbConf.Password,
	})
	if err != nil {
		log.Errorf("Influx connection failed: %s", err.Error())
		time.Sleep(5 * time.Second)
		i.connect()
	}

	return c
}

// WriteData write analyzed data to influx persistent back-end storage.
func (i *InfluxPump) WriteData(ctx context.Context, data []interface{}) error {
	c := i.connect()
	defer c.Close()

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  i.dbConf.DatabaseName,
		Precision: "us",
	})

	var pt *client.Point
	var err error

	//	 Create a point and add to batch
	for _, v := range data {
		// Convert to AnalyticsRecord
		decoded, _ := v.(analytics.AnalyticsRecord)
		mapping := map[string]interface{}{
			"timestamp":  decoded.TimeStamp,
			"username":   decoded.Username,
			"effect":     decoded.Effect,
			"conclusion": decoded.Conclusion,
			"request":    decoded.Request,
			"policies":   decoded.Policies,
			"deciders":   decoded.Deciders,
			"expireAt":   decoded.ExpireAt,
		}

		tags := make(map[string]string)
		fields := make(map[string]interface{})

		// Select tags from config
		for _, t := range i.dbConf.Tags {
			var tag string

			if b, e := json.Marshal(mapping[t]); e != nil {
				tag = ""
			} else {
				// convert and remove surrounding quotes from tag value
				tag = strings.Trim(string(b), "\"")
			}

			tags[t] = tag
		}

		// Select field from config
		for _, f := range i.dbConf.Fields {
			fields[f] = mapping[f]
		}

		// New record
		if pt, err = client.NewPoint(table, tags, fields, time.Now()); err != nil {
			log.Error(err.Error())

			continue
		}

		// Add point to batch point
		bp.AddPoint(pt)
	}

	// Now that all points are added, write the batch
	_ = c.Write(bp)

	return nil
}
