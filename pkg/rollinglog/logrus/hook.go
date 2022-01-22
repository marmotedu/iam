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

func (h hook) Fire(entry *logrus.Entry) error {
	fields := make([]zap.Field, 0, 15)

	for key, val := range entry.Data {
		if key == logrus.ErrorKey {
			fields = append(fields, zap.Error(val.(error)))
		} else {
			fields = append(fields, zap.Any(key, val))
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
	case logrus.DebugLevel:
		h.Write(zapcore.DebugLevel, entry.Message, fields, entry.Caller)
	}

	return nil
}

func (h *hook) Write(lv zapcore.Level, msg string, fields []zap.Field, caller *runtime.Frame) {
	if ce := h.logger.Check(lv, msg); ce != nil {
		if caller != nil {
			ce.Caller = zapcore.NewEntryCaller(caller.PC, caller.File, caller.Line, caller.PC != 0)
		}
		ce.Write(fields...)
	}
}

func newHook(logger *zap.Logger) *hook {
	return &hook{logger: logger}
}
