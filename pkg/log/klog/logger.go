/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

// Package klog init klog logger. klog is used by kubernetes, this can compatible the kubernetes packages.
package klog

import (
	"flag"

	"go.uber.org/zap"
	"k8s.io/klog"
)

// InitLogger init klog by zap logger.
func InitLogger(zapLogger *zap.Logger) {
	fs := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(fs)
	defer klog.Flush()
	klog.SetOutputBySeverity("INFO", &infoLogger{logger: zapLogger})
	klog.SetOutputBySeverity("WARNING", &warnLogger{logger: zapLogger})
	klog.SetOutputBySeverity("FATAL", &fatalLogger{logger: zapLogger})
	klog.SetOutputBySeverity("ERROR", &errorLogger{logger: zapLogger})
	_ = fs.Set("skip_headers", "true")
	_ = fs.Set("logtostderr", "false")
}

type infoLogger struct {
	logger *zap.Logger
}

func (l *infoLogger) Write(p []byte) (n int, err error) {
	l.logger.Info(string(p[:len(p)-1]))
	return len(p), nil
}

type warnLogger struct {
	logger *zap.Logger
}

func (l *warnLogger) Write(p []byte) (n int, err error) {
	l.logger.Warn(string(p[:len(p)-1]))
	return len(p), nil
}

type fatalLogger struct {
	logger *zap.Logger
}

func (l *fatalLogger) Write(p []byte) (n int, err error) {
	l.logger.Fatal(string(p[:len(p)-1]))
	return len(p), nil
}

type errorLogger struct {
	logger *zap.Logger
}

func (l *errorLogger) Write(p []byte) (n int, err error) {
	l.logger.Error(string(p[:len(p)-1]))
	return len(p), nil
}
