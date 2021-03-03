// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package validate validate the basic environment for iamctl to run.
package validate

import (
	"fmt"
	"net"
	"net/url"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marmotedu/iam/internal/iamctl"
	cmdutil "github.com/marmotedu/iam/internal/iamctl/cmd/util"
	"github.com/marmotedu/iam/internal/iamctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

// ValidateOptions is an options struct to support 'validate' sub command.
type ValidateOptions struct {
	genericclioptions.IOStreams
}

// ValidateInfo defines the validate information.
type ValidateInfo struct {
	ItemName string
	Status   string
	Message  string
}

var validateExample = templates.Examples(`
		# Validate the basic environment for iamctl to run
		iamctl validate`)

// NewValidateOptions returns an initialized ValidateOptions instance.
func NewValidateOptions(ioStreams genericclioptions.IOStreams) *ValidateOptions {
	return &ValidateOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdValidate returns new initialized instance of 'validate' sub command.
func NewCmdValidate(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewValidateOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "validate",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Validate the basic environment for iamctl to run",
		TraverseChildren:      true,
		Long:                  "Validate the basic environment for iamctl to run.",
		Example:               validateExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	return cmd
}

// Complete completes all the required options.
func (o *ValidateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *ValidateOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a validate sub command using the specified options.
func (o *ValidateOptions) Run(args []string) error {
	data := [][]string{}
	FAIL := color.RedString("FAIL")
	PASS := color.GreenString("PASS")
	validateInfo := ValidateInfo{}

	// check if can access db
	validateInfo.ItemName = "iam-apiserver"
	target, err := url.Parse(viper.GetString("server.address"))
	if err != nil {
		return err
	}
	_, err = net.Dial("tcp", target.Host)
	// defer client.Close()
	if err != nil {
		validateInfo.Status = FAIL
		validateInfo.Message = fmt.Sprintf("%v", err)
	} else {
		validateInfo.Status = PASS
		validateInfo.Message = ""
	}

	data = append(data, []string{validateInfo.ItemName, validateInfo.Status, validateInfo.Message})

	table := tablewriter.NewWriter(o.Out)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(iamctl.TableWidth)
	table.SetHeader([]string{"ValidateItem", "Result", "Message"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render()

	return nil
}
