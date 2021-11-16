// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pumps

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/marmotedu/errors"
	"github.com/mitchellh/mapstructure"

	"github.com/marmotedu/iam/internal/pump/analytics"
	"github.com/marmotedu/iam/pkg/log"
)

// CSVPump defines a csv pump with csv specific options and common options.
type CSVPump struct {
	csvConf *CSVConf
	CommonPumpConfig
}

// CSVConf defines csv specific options.
type CSVConf struct {
	// Specify the directory used to store automatically generated csv file which contains analyzed data.
	CSVDir string `mapstructure:"csv_dir"`
}

// New create a csv pump instance.
func (c *CSVPump) New() Pump {
	newPump := CSVPump{}

	return &newPump
}

// GetName returns the csv pump name.
func (c *CSVPump) GetName() string {
	return "CSV Pump"
}

// Init initialize the csv pump instance.
func (c *CSVPump) Init(conf interface{}) error {
	c.csvConf = &CSVConf{}
	err := mapstructure.Decode(conf, &c.csvConf)
	if err != nil {
		log.Fatalf("Failed to decode configuration: %s", err.Error())
	}

	ferr := os.MkdirAll(c.csvConf.CSVDir, 0o777)
	if ferr != nil {
		log.Error(ferr.Error())
	}

	log.Debug("CSV Initialized")

	return nil
}

// WriteData write analyzed data to csv persistent back-end storage.
func (c *CSVPump) WriteData(ctx context.Context, data []interface{}) error {
	curtime := time.Now()
	fname := fmt.Sprintf("%d-%s-%d-%d.csv", curtime.Year(), curtime.Month().String(), curtime.Day(), curtime.Hour())
	fname = path.Join(c.csvConf.CSVDir, fname)

	var outfile *os.File
	var appendHeader bool

	if _, err := os.Stat(fname); os.IsNotExist(err) {
		var createErr error
		outfile, createErr = os.Create(fname)
		if createErr != nil {
			log.Errorf("Failed to create new CSV file: %s", createErr.Error())
		}
		appendHeader = true
	} else {
		var appendErr error
		outfile, appendErr = os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0o600)
		if appendErr != nil {
			log.Errorf("Failed to open CSV file: %s", appendErr.Error())
		}
	}

	defer outfile.Close()
	writer := csv.NewWriter(outfile)

	if appendHeader {
		startRecord := analytics.AnalyticsRecord{}
		headers := startRecord.GetFieldNames()

		err := writer.Write(headers)
		if err != nil {
			log.Errorf("Failed to write file headers: %s", err.Error())

			return errors.Wrap(err, "failed to write file headers")
		}
	}

	for _, v := range data {
		decoded, _ := v.(analytics.AnalyticsRecord)

		toWrite := decoded.GetLineValues()
		err := writer.Write(toWrite)
		if err != nil {
			log.Error("File write failed!")
			log.Error(err.Error())
		}
	}

	writer.Flush()

	return nil
}
