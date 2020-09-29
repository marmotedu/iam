// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package options is the public flags and options used by a generic api
// server. It takes a minimal set of dependencies and does not reference
// implementations, in order to ensure it may be reused by multiple components
// (such as CLI commands that wish to generate or validate config).
package options // import "github.com/marmotedu/iam/internal/pkg/options"
