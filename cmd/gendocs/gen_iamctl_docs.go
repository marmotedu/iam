// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/marmotedu/iam/internal/iamctl/cmd"
	"github.com/marmotedu/iam/pkg/util/genutil"
)

func main() {
	// use os.Args instead of "flags" because "flags" will mess up the man pages!
	path := "docs/"
	if len(os.Args) == 2 {
		path = os.Args[1]
	} else if len(os.Args) > 2 {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s [output directory]\n", os.Args[0])
		os.Exit(1)
	}

	outDir, err := genutil.OutDir(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}

	// Set environment variables used by iamctl so the output is consistent,
	// regardless of where we run.
	_ = os.Setenv("HOME", "/home/username")
	// TODO os.Stdin should really be something like ioutil.Discard, but a Reader
	iamctl := cmd.NewIAMCtlCommand(os.Stdin, ioutil.Discard, ioutil.Discard)
	_ = doc.GenMarkdownTree(iamctl, outDir)
}
