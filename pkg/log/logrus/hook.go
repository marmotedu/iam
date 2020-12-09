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

package logrus

import (
	"runtime"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type hook struct {
	logger *zap.Logger
}

func (h *hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *hook) Fire(entry *logrus.Entry) error {
	fields := make([]zap.Field, 0, 10)

	for key, value := range entry.Data {
		if key == logrus.ErrorKey {
			fields = append(fields, zap.Error(value.(error)))
		} else {
			fields = append(fields, zap.Any(key, value))
		}
	}

	switch entry.Level {
	case logrus.PanicLevel:
		h.Write(zapcore.PanicLevel, entry.Message, fields, entry.Caller)
	case logrus.FatalLevel:
		h.Write(zapcore.FatalLevel, entry.Message, fields, entry.Caller)
	case logrus.ErrorLevel:
		h.Write(zapcore.ErrorLevel, entry.Message, fields, entry.Caller)
	case logrus.WarnLevel:
		h.Write(zapcore.WarnLevel, entry.Message, fields, entry.Caller)
	case logrus.InfoLevel:
		h.Write(zapcore.InfoLevel, entry.Message, fields, entry.Caller)
	case logrus.DebugLevel, logrus.TraceLevel:
		h.Write(zapcore.DebugLevel, entry.Message, fields, entry.Caller)
	}

	return nil
}

func (h *hook) Write(lvl zapcore.Level, msg string, fields []zap.Field, caller *runtime.Frame) {
	if ce := h.logger.Check(lvl, msg); ce != nil {
		if caller != nil {
			ce.Caller = zapcore.NewEntryCaller(caller.PC, caller.File, caller.Line, caller.PC != 0)
		}
		ce.Write(fields...)
	}
}

func newHook(logger *zap.Logger) logrus.Hook {
	return &hook{logger: logger}
}
