// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jwt

import (
	"fmt"
	"regexp"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cobra"

	cmdutil "github.com/marmotedu/iam/internal/iamctl/cmd/util"
	"github.com/marmotedu/iam/internal/iamctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

const (
	veirfyUsageStr = "veirfy SECRETKEY TOKEN"
)

// VerifyOptions is an options struct to support verify subcommands.
type VerifyOptions struct {
	Compact bool
	Debug   bool

	genericclioptions.IOStreams
}

var (
	verifyExample = templates.Examples(`
		# Verify a JWT token
		iamctl jwt verify XXX xxxxx.yyyyy.zzzzz`)

	verifyUsageErrStr = fmt.Sprintf(
		"expected '%s'.\nSECRETKEY and TOKEN are required arguments for the subcmd1 command",
		veirfyUsageStr,
	)
)

// NewVerifyOptions returns an initialized VerifyOptions instance.
func NewVerifyOptions(ioStreams genericclioptions.IOStreams) *VerifyOptions {
	return &VerifyOptions{
		Compact: false,
		Debug:   false,

		IOStreams: ioStreams,
	}
}

// NewCmdVerify returns new initialized instance of verify sub command.
func NewCmdVerify(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewVerifyOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "verify",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"sub2"},
		Short:                 "Verify a JWT token",
		Long:                  "Verify a JWT token",
		TraverseChildren:      true,
		Example:               verifyExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return cmdutil.UsageErrorf(cmd, verifyUsageErrStr)
			}

			return nil
		},
	}

	// mark flag as deprecated
	cmd.Flags().BoolVar(&o.Compact, "compact", o.Compact, "Output compact JSON.")
	cmd.Flags().BoolVar(&o.Debug, "debug", o.Debug, "Print out all kinds of debug data.")

	return cmd
}

// Complete completes all the required options.
func (o *VerifyOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *VerifyOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Run executes a verify subcommand using the specified options.
func (o *VerifyOptions) Run(args []string) error {
	// get the token
	tokenData := []byte(args[1])

	// trim possible whitespace from token
	tokenData = regexp.MustCompile(`\s*$`).ReplaceAll(tokenData, []byte{})

	// Parse the token.  Load the key from command line option
	token, err := jwt.Parse(string(tokenData), func(t *jwt.Token) (interface{}, error) {
		return []byte(args[0]), nil
	})

	// Print some debug data
	if o.Debug && token != nil {
		fmt.Println("Header:")
		if pErr := printJSON(o.Compact, token.Header); pErr != nil {
			return fmt.Errorf("failed to output header: %w", pErr)
		}

		fmt.Println("Claims:")
		if pErr := printJSON(o.Compact, token.Claims); pErr != nil {
			return fmt.Errorf("failed to output claims: %w", pErr)
		}
	}

	// Print an error if we can't parse for some reason
	if err != nil {
		return fmt.Errorf("couldn't parse token: %w", err)
	}

	// Is token invalid?
	if !token.Valid {
		return fmt.Errorf("token is invalid")
	}

	if !o.Debug {
		// Print the token details
		if err := printJSON(o.Compact, token.Claims); err != nil {
			return fmt.Errorf("failed to output claims: %w", err)
		}
	}

	return nil
}
