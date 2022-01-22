package rollinglog

import (
	"os"

	"github.com/marmotedu/iam/pkg/rollinglog/rolling"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// generate zapcore WriteSyncer, The category must be OutputPaths or ErrorOutputPaths.
func generateWriterSyncer(category string, opts *Options) []zapcore.WriteSyncer {
	var (
		ws    zapcore.WriteSyncer
		err   error
		paths = map[string][]string{
			"OutputPaths":      opts.OutputPaths,
			"ErrorOutputPaths": opts.ErrorOutputPaths,
		}
	)
	wss := make([]zapcore.WriteSyncer, 0, len(paths[category]))
	osStds := map[string]*os.File{"stdout": os.Stdout, "stderr": os.Stderr}

	for _, path := range paths[category] {
		if path == "stdout" || path == "stderr" {
			wss = append(wss, zapcore.AddSync(osStds[path]))
			continue
		}

		if opts.Rolling {
			ws = zapcore.Lock(zapcore.AddSync(buildRollingLogger(path, opts)))
		} else {
			ws, _, err = zap.Open(path)
			if err != nil {
				panic(err.Error())
			}

		}
		wss = append(wss, ws)
	}

	return wss
}

// generate rotate logger
func buildRollingLogger(path string, opts *Options) *rolling.Logger {
	return rolling.NewLogger(
		path,
		rolling.WithMaxAge(opts.RollingMaxAge),
		rolling.WithMaxBackups(opts.RollingMaxBackups),
		rolling.WithMaxSize(opts.RollingMaxSize),
		rolling.WithLocaltime(opts.RollingLocalTime),
		rolling.WithCompress(opts.RollingCompress),
	)
}

func multiWriteSyncer(ws ...zapcore.WriteSyncer) zapcore.WriteSyncer {
	return zapcore.NewMultiWriteSyncer(ws...)
}
