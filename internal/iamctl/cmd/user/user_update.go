// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

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
	updateUsageStr = "update USERNAME"
)

// UpdateOptions is an options struct to support update subcommands.
type UpdateOptions struct {
	Name     string
	Nickname string
	Email    string
	Phone    string

	iamclient iam.IamInterface
	genericclioptions.IOStreams
}

var (
	updateLong = templates.LongDesc(`Update a user resource. 

Can only update nickname, email and phone.

NOTICE: field will be updated to zero value if not specified.`)

	updateExample = templates.Examples(`
		# Update use foo's information
		iamctl user update foo --nickname=foo2 --email=foo@qq.com --phone=1812883xxxx`)

	updateUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nUSERNAME is required arguments for the update command",
		updateUsageStr,
	)
)

// NewUpdateOptions returns an initialized UpdateOptions instance.
func NewUpdateOptions(ioStreams genericclioptions.IOStreams) *UpdateOptions {
	return &UpdateOptions{
		IOStreams: ioStreams,
	}
}

// NewCmdUpdate returns new initialized instance of update sub command.
func NewCmdUpdate(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewUpdateOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   updateUsageStr,
		DisableFlagsInUseLine: true,
		Aliases:               []string{},
		Short:                 "Update a user resource",
		TraverseChildren:      true,
		Long:                  updateLong,
		Example:               updateExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}

	cmd.Flags().StringVar(&o.Nickname, "nickname", o.Nickname, "The nickname of the user.")
	cmd.Flags().StringVar(&o.Email, "email", o.Email, "The email of the user.")
	cmd.Flags().StringVar(&o.Phone, "phone", o.Phone, "The phone number of the user.")

	return cmd
}

// Complete completes all the required options.
func (o *UpdateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) == 0 {
		return cmdutil.UsageErrorf(cmd, updateUsageErrStr)
	}

	o.Name = args[0]
	o.iamclient, err = f.IAMClient()
	if err != nil {
		return err
	}

	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *UpdateOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes an update subcommand using the specified options.
func (o *UpdateOptions) Run(args []string) error {
	user, err := o.iamclient.APIV1().Users().Get(context.TODO(), o.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if o.Nickname != "" {
		user.Nickname = o.Nickname
	}
	if o.Email != "" {
		user.Email = o.Email
	}
	if o.Phone != "" {
		user.Phone = o.Phone
	}

	ret, err := o.iamclient.APIV1().Users().Update(context.TODO(), user, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, "user/%s updated\n", ret.Name)

	return nil
}
