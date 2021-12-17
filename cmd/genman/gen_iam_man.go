// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	mangen "github.com/cpuguy83/go-md2man/v2/md2man"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/marmotedu/iam/internal/apiserver"
	"github.com/marmotedu/iam/internal/authzserver"
	"github.com/marmotedu/iam/internal/iamctl/cmd"
	"github.com/marmotedu/iam/internal/pump"
	"github.com/marmotedu/iam/internal/watcher"
	"github.com/marmotedu/iam/pkg/util/genutil"
)

func main() {
	// use os.Args instead of "flags" because "flags" will mess up the man pages!
	path := "docs/man/man1"
	module := ""
	if len(os.Args) == 3 {
		path = os.Args[1]
		module = os.Args[2]
	} else {
		fmt.Fprintf(os.Stderr, "usage: %s [output directory] [module] \n", os.Args[0])
		os.Exit(1)
	}

	outDir, err := genutil.OutDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}

	// Set environment variables used by command so the output is consistent,
	// regardless of where we run.
	_ = os.Setenv("HOME", "/home/username")

	switch module {
	case "iam-apiserver":
		// generate manpage for iam-apiserver
		apiServer := apiserver.NewApp("iam-apiserver").Command()
		genMarkdown(apiServer, "", outDir)
		for _, c := range apiServer.Commands() {
			genMarkdown(c, "iam-apiserver", outDir)
		}
	case "iam-authz-server":
		// generate manpage for iam-authz-server
		authzServer := authzserver.NewApp("iam-authz-server").Command()
		genMarkdown(authzServer, "", outDir)
		for _, c := range authzServer.Commands() {
			genMarkdown(c, "iam-authz-server", outDir)
		}
	case "iam-pump":
		// generate manpage for iam-pump
		pump := pump.NewApp("iam-pump").Command()
		genMarkdown(pump, "", outDir)
		for _, c := range pump.Commands() {
			genMarkdown(c, "iam-pump", outDir)
		}
	case "iam-watcher":
		// generate manpage for iam-watcher
		watcher := watcher.NewApp("iam-watcher").Command()
		genMarkdown(watcher, "", outDir)
		for _, c := range watcher.Commands() {
			genMarkdown(c, "iam-watcher", outDir)
		}
	case "iamctl":
		// generate manpage for iamctl
		// TODO os.Stdin should really be something like ioutil.Discard, but a Reader
		iamctl := cmd.NewDefaultIAMCtlCommand()
		genMarkdown(iamctl, "", outDir)
		for _, c := range iamctl.Commands() {
			genMarkdown(c, "iamctl", outDir)
		}
	default:
		fmt.Fprintf(os.Stderr, "Module %s is not supported", module)
		os.Exit(1)
	}
}

func preamble(out *bytes.Buffer, name, short, long string) {
	out.WriteString(`% IAM(1) iam User Manuals
% Eric Paris
% Jan 2015
# NAME
`)
	fmt.Fprintf(out, "%s \\- %s\n\n", name, short)
	fmt.Fprintf(out, "# SYNOPSIS\n")
	fmt.Fprintf(out, "**%s** [OPTIONS]\n\n", name)
	fmt.Fprintf(out, "# DESCRIPTION\n")
	fmt.Fprintf(out, "%s\n\n", long)
}

func printFlags(out io.Writer, flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		format := "**--%s**=%s\n\t%s\n\n"
		if flag.Value.Type() == "string" {
			// put quotes on the value
			format = "**--%s**=%q\n\t%s\n\n"
		}

		// Todo, when we mark a shorthand is deprecated, but specify an empty message.
		// The flag.ShorthandDeprecated is empty as the shorthand is deprecated.
		// Using len(flag.ShorthandDeprecated) > 0 can't handle this, others are ok.
		if !(len(flag.ShorthandDeprecated) > 0) && len(flag.Shorthand) > 0 {
			format = "**-%s**, " + format
			fmt.Fprintf(out, format, flag.Shorthand, flag.Name, flag.DefValue, flag.Usage)
		} else {
			fmt.Fprintf(out, format, flag.Name, flag.DefValue, flag.Usage)
		}
	})
}

func printOptions(out io.Writer, command *cobra.Command) {
	flags := command.NonInheritedFlags()
	if flags.HasFlags() {
		fmt.Fprintf(out, "# OPTIONS\n")
		printFlags(out, flags)
		fmt.Fprintf(out, "\n")
	}
	flags = command.InheritedFlags()
	if flags.HasFlags() {
		fmt.Fprintf(out, "# OPTIONS INHERITED FROM PARENT COMMANDS\n")
		printFlags(out, flags)
		fmt.Fprintf(out, "\n")
	}
}

func genMarkdown(command *cobra.Command, parent, docsDir string) {
	dparent := strings.ReplaceAll(parent, " ", "-")
	name := command.Name()

	dname := name
	if len(parent) > 0 {
		dname = dparent + "-" + name
		name = parent + " " + name
	}

	out := new(bytes.Buffer)

	short, long := command.Short, command.Long
	if len(long) == 0 {
		long = short
	}

	preamble(out, name, short, long)
	printOptions(out, command)

	if len(command.Example) > 0 {
		fmt.Fprintf(out, "# EXAMPLE\n")
		fmt.Fprintf(out, "```\n%s\n```\n", command.Example)
	}

	if len(command.Commands()) > 0 || len(parent) > 0 {
		fmt.Fprintf(out, "# SEE ALSO\n")

		if len(parent) > 0 {
			fmt.Fprintf(out, "**%s(1)**, ", dparent)
		}

		for _, c := range command.Commands() {
			fmt.Fprintf(out, "**%s-%s(1)**, ", dname, c.Name())
			genMarkdown(c, name, docsDir)
		}

		fmt.Fprintf(out, "\n")
	}

	out.WriteString(`
# HISTORY
January 2015, Originally compiled by Eric Paris (eparis at redhat dot com) based on the marmotedu source material, but hopefully they have been automatically generated since!
`)

	final := mangen.Render(out.Bytes())

	filename := docsDir + dname + ".1"

	outFile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()

	_, err = outFile.Write(final)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
