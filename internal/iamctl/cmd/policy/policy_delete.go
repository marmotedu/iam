// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package policy

import (
	"context"
	"fmt"

	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/marmotedu-sdk-go/marmotedu/service/iam"
	"github.com/spf13/cobra"

	cmdutil "github.com/marmotedu/iam/internal/iamctl/cmd/util"
	"github.com/marmotedu/iam/internal/iamctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

const (
	deleteUsageStr = "delete POLICY_NAME"
)

// DeleteOptions is an options struct to support delete subcommands.
type DeleteOptions struct {
	Name string

	iamclient iam.IamInterface
	genericclioptions.IOStreams
}

var (
	deleteExample = templates.Examples(`
		# Delete a policy resource
		iamctl policy delete foo`)

	deleteUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nPOLICY_NAME is required arguments for the delete command",
		deleteUsageStr,
	)
)

// NewDeleteOptions returns an initialized DeleteOptions instance.
func NewDeleteOptions(ioStreams genericclioptions.IOStreams) *DeleteOptions {
	return &DeleteOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdDelete returns new initialized instance of delete sub command.
func NewCmdDelete(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewDeleteOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "delete",
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Delete a authorization policy resource",
		TraverseChildren:      true,
		Long:                  "Delete a authorization policy resource.",
		Example:               deleteExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run())
		},
		SuggestFor: []string{},
	}

	return cmd
}

// Complete completes all the required options.
func (o *DeleteOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, deleteUsageErrStr)
	}

	o.Name = args[0]

	o.iamclient, err = f.IAMClient()
	if err != nil {
		return err
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *DeleteOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a delete subcommand using the specified options.
func (o *DeleteOptions) Run() error {
	if err := o.iamclient.APIV1().Policies().Delete(context.TODO(), o.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "policy/%s deleted\n", o.Name)

	return nil
}
