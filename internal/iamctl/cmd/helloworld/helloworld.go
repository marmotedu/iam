// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package helloworld

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdutil "github.com/marmotedu/iam/internal/iamctl/cmd/util"
	"github.com/marmotedu/iam/internal/iamctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

const (
	helloworldUsageStr = "helloworld USERNAME PASSWORD"
	maxStringLength    = 17
)

// HelloworldOptions is an options struct to support 'helloworld' sub command.
type HelloworldOptions struct {
	// options
	StringOption      string
	StringSliceOption []string
	IntOption         int
	BoolOption        bool

	// args
	Username string
	Password string

	genericclioptions.IOStreams
}

var (
	helloworldLong = templates.LongDesc(`A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`)

	helloworldExample = templates.Examples(`
		# Print all option values for helloworld 
		iamctl helloworld marmotedu marmotedupass`)

	helloworldUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nUSERNAME and PASSWORD are required arguments for the helloworld command",
		helloworldUsageStr,
	)
)

// NewHelloworldOptions returns an initialized HelloworldOptions instance.
func NewHelloworldOptions(ioStreams genericclioptions.IOStreams) *HelloworldOptions {
	return &HelloworldOptions{
		StringOption: "default",
		IOStreams:    ioStreams,
	}
}

// NewCmdHelloworld returns new initialized instance of 'helloworld' sub command.
func NewCmdHelloworld(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewHelloworldOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   helloworldUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "A brief description of your command",
		TraverseChildren:      true,
		Long:                  helloworldLong,
		Example:               helloworldExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
		Args: func(cmd *cobra.Command, args []string) error {
			// nolint: gomnd // no need
			if len(args) < 2 {
				return cmdutil.UsageErrorf(cmd, helloworldUsageErrStr)
			}

			// if need args equal to zero, uncomment the following code
			/*
				if len(args) != 0 {
					return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
				}
			*/

			return nil
		},
	}

	// mark flag as deprecated
	_ = cmd.Flags().MarkDeprecated("deprecated-opt", "This flag is deprecated and will be removed in future.")
	cmd.Flags().StringVarP(&o.StringOption, "string", "", o.StringOption, "String option.")
	cmd.Flags().StringSliceVar(&o.StringSliceOption, "slice", o.StringSliceOption, "String slice option.")
	cmd.Flags().IntVarP(&o.IntOption, "int", "i", o.IntOption, "Int option.")
	cmd.Flags().BoolVarP(&o.BoolOption, "bool", "b", o.BoolOption, "Bool option.")

	return cmd
}

// Complete completes all the required options.
func (o *HelloworldOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	if o.StringOption != "" {
		o.StringOption += "(complete)"
	}

	o.Username = args[0]
	o.Password = args[1]

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *HelloworldOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(o.StringOption) > maxStringLength {
		return cmdutil.UsageErrorf(cmd, "--string length must less than 18")
	}

	if o.IntOption < 0 {
		return cmdutil.UsageErrorf(cmd, "--int must be a positive integer: %v", o.IntOption)
	}

	return nil
}

// Run executes a helloworld sub command using the specified options.
func (o *HelloworldOptions) Run(args []string) error {
	fmt.Fprintf(o.Out, "The following is option values:\n")
	fmt.Fprintf(o.Out, "==> --string: %v\n==> --slice: %v\n==> --int: %v\n==> --bool: %v\n",
		o.StringOption, o.StringSliceOption, o.IntOption, o.BoolOption)

	fmt.Fprintf(o.Out, "\nThe following is args values:\n")
	fmt.Fprintf(o.Out, "==> username: %v\n==> password: %v\n", o.Username, o.Password)

	return nil
}
