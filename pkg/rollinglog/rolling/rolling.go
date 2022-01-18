package rolling

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*lumberjack.Logger
}

type Opt func(*Logger)

// WithMaxSize set log maximum size in MB, defaults to 100 MB.
func WithMaxSize(s int) Opt {
	return func(logger *Logger) {
		logger.MaxSize = s
	}
}

// WithMaxAge set maximum Reserve of days, the default is not to remove old logs.
func WithMaxAge(a int) Opt {
	return func(logger *Logger) {
		logger.MaxAge = a
	}
}

// WithMaxBackups set maximum amount old log files to retain,
// the default is to retain all old logs.
func WithMaxBackups(bu int) Opt {
	return func(logger *Logger) {
		logger.MaxBackups = bu
	}
}

// WithLocaltime use local time.
func WithLocaltime(lt bool) Opt {
	return func(logger *Logger) {
		logger.LocalTime = lt
	}
}

func WithCompress(c bool) Opt {
	return func(logger *Logger) {
		logger.Compress = c
	}
}

// NewLogger  return rotate Logger.
func NewLogger(filename string, opts ...Opt) *Logger {
	l := &Logger{
		Logger: &lumberjack.Logger{
			Filename: filename,
		},
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}
