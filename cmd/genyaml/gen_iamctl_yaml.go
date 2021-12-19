// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/marmotedu/iam/internal/iamctl/cmd"
	"github.com/marmotedu/iam/pkg/util/genutil"
)

type cmdOption struct {
	Name         string
	Shorthand    string `yaml:",omitempty"`
	DefaultValue string `yaml:"defaultValue,omitempty"`
	Usage        string `yaml:",omitempty"`
}

type cmdDoc struct {
	Name             string
	Synopsis         string      `yaml:",omitempty"`
	Description      string      `yaml:",omitempty"`
	Options          []cmdOption `yaml:",omitempty"`
	InheritedOptions []cmdOption `yaml:"inheritedOptions,omitempty"`
	Example          string      `yaml:",omitempty"`
	SeeAlso          []string    `yaml:"seeAlso,omitempty"`
}

func main() {
	path := "docs/yaml/iamctl"
	if len(os.Args) == 2 {
		path = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [output directory]\n", os.Args[0])
		os.Exit(1)
	}

	outDir, err := genutil.OutDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}

	// Set environment variables used by iamctl so the output is consistent,
	// regardless of where we run.
	os.Setenv("HOME", "/home/username")
	iamctl := cmd.NewDefaultIAMCtlCommand()
	genYaml(iamctl, "", outDir)
	for _, c := range iamctl.Commands() {
		genYaml(c, "iamctl", outDir)
	}
}

// Temporary workaround for yaml lib generating incorrect yaml with long strings
// that do not contain \n.
func forceMultiLine(s string) string {
	if len(s) > 60 && !strings.Contains(s, "\n") {
		s += "\n"
	}

	return s
}

func genFlagResult(flags *pflag.FlagSet) []cmdOption {
	result := []cmdOption{}

	flags.VisitAll(func(flag *pflag.Flag) {
		// Todo, when we mark a shorthand is deprecated, but specify an empty message.
		// The flag.ShorthandDeprecated is empty as the shorthand is deprecated.
		// Using len(flag.ShorthandDeprecated) > 0 can't handle this, others are ok.
		if !(len(flag.ShorthandDeprecated) > 0) && len(flag.Shorthand) > 0 {
			opt := cmdOption{
				flag.Name,
				flag.Shorthand,
				flag.DefValue,
				forceMultiLine(flag.Usage),
			}
			result = append(result, opt)
		} else {
			opt := cmdOption{
				Name:         flag.Name,
				DefaultValue: forceMultiLine(flag.DefValue),
				Usage:        forceMultiLine(flag.Usage),
			}
			result = append(result, opt)
		}
	})

	return result
}

func genYaml(command *cobra.Command, parent, docsDir string) {
	doc := cmdDoc{}

	doc.Name = command.Name()
	doc.Synopsis = forceMultiLine(command.Short)
	doc.Description = forceMultiLine(command.Long)

	flags := command.NonInheritedFlags()
	if flags.HasFlags() {
		doc.Options = genFlagResult(flags)
	}
	flags = command.InheritedFlags()
	if flags.HasFlags() {
		doc.InheritedOptions = genFlagResult(flags)
	}

	if len(command.Example) > 0 {
		doc.Example = command.Example
	}

	if len(command.Commands()) > 0 || len(parent) > 0 {
		result := []string{}
		if len(parent) > 0 {
			result = append(result, parent)
		}
		for _, c := range command.Commands() {
			result = append(result, c.Name())
		}
		doc.SeeAlso = result
	}

	final, err := yaml.Marshal(&doc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var filename string

	if parent == "" {
		filename = docsDir + doc.Name + ".yaml"
	} else {
		filename = docsDir + parent + "_" + doc.Name + ".yaml"
	}

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
