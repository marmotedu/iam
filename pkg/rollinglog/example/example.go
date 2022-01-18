package main

import "github.com/marmotedu/iam/pkg/rollinglog"

func main() {
	opts := &rollinglog.Options{
		Level:            "debug",
		Format:           "json",
		EnableColor:      false,
		EnableCaller:     false,
		OutputPaths:      []string{"test.log"},
		ErrorOutputPaths: []string{"error.log"},
		Rolling:          true,
		RollingMaxSize:   1,
	}
	// 初始化全局logger
	rollinglog.Init(opts)
	defer rollinglog.Flush()

	for i := 0; i < 30000; i++ {
		// rollinglog.Debug("This is a debug message")
		rollinglog.Warnf("This is a formatted %s message", "hello")
	}
}
