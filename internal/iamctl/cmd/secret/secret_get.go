// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package secret

import (
	"context"
	"fmt"
	"time"

	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/marmotedu-sdk-go/marmotedu/service/iam"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	cmdutil "github.com/marmotedu/iam/internal/iamctl/cmd/util"
	"github.com/marmotedu/iam/internal/iamctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

const (
	getUsageStr = "get SECRET_NAME"
)

// GetOptions is an options struct to support get subcommands.
type GetOptions struct {
	Name string

	iamclient iam.IamInterface

	genericclioptions.IOStreams
}

var (
	getExample = templates.Examples(`
		# Get a specified secret information
		iamctl secret get foo`)

	getUsageErrStr = fmt.Sprintf("expected '%s'.\nSECRET_NAME is required arguments for the get command", getUsageStr)
)

// NewGetOptions returns an initialized GetOptions instance.
func NewGetOptions(ioStreams genericclioptions.IOStreams) *GetOptions {
	return &GetOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdGet returns new initialized instance of get sub command.
func NewCmdGet(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewGetOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "get SECRET_NAME",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Display a secret resource",
		TraverseChildren:      true,
		Long:                  "Display a secret resource.",
		Example:               getExample,
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
func (o *GetOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, getUsageErrStr)
	}

	o.Name = args[0]

	o.iamclient, err = f.IAMClient()
	if err != nil {
		return err
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *GetOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a get subcommand using the specified options.
func (o *GetOptions) Run(args []string) error {
	secret, err := o.iamclient.APIV1().Secrets().Get(context.TODO(), o.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(o.Out)

	data := [][]string{
		{
			secret.Name,
			secret.SecretID,
			secret.SecretKey,
			time.Unix(secret.Expires, 0).Format("2006-01-02 15:04:05"),
			secret.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	table = setHeader(table)
	table = cmdutil.TableWriterDefaultConfig(table)
	table.AppendBulk(data)
	table.Render()

	return nil
}
