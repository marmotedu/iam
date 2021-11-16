// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cronlog

import (
	"fmt"

	"go.uber.org/zap"
)

type logger struct {
	zapLogger *zap.SugaredLogger
}

// NewLogger create a logger which implement `github.com/robfig/cron.Logger`.
func NewLogger(zapLogger *zap.SugaredLogger) logger {
	return logger{zapLogger: zapLogger}
}

func (l logger) Info(msg string, args ...interface{}) {
	l.zapLogger.Infow(msg, args...)
}

func (l logger) Error(err error, msg string, args ...interface{}) {
	l.zapLogger.Errorw(fmt.Sprintf(msg, args...), "error", err.Error())
}

func (l logger) Flush() {
	_ = l.zapLogger.Sync()
}
