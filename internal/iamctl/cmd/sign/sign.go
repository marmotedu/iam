// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package sign used to sign a jwt token with given secretId and secretKey.
package sign

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	jwt "github.com/dgrijalva/jwt-go/v4"
	cmdutil "github.com/marmotedu/iam/internal/iamctl/cmd/util"
	"github.com/marmotedu/iam/internal/iamctl/util/templates"
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

// ErrSigningMethod defines invalid signing method error.
var ErrSigningMethod = errors.New("invalid signing method")

// SignOptions is an options struct to support 'sign' sub command.
type SignOptions struct {
	Timeout   time.Duration
	Algorithm string

	genericclioptions.IOStreams
}

var (
	signExample = templates.Examples(`
		# Sign a token with secretID and secretKey
		iamctl sign tgydj8d9EQSnFqKf iBdEdFNBLN1nR3fV

		# Sign a token with expires and sign method
		iamctl sign tgydj8d9EQSnFqKf iBdEdFNBLN1nR3fV --timeout=2h --algorithm=HS256`)
)

// NewSignOptions returns an initialized SignOptions instance.
func NewSignOptions(ioStreams genericclioptions.IOStreams) *SignOptions {
	return &SignOptions{
		Timeout:   2 * time.Hour,
		Algorithm: "HS256",
		IOStreams: ioStreams,
	}
}

// NewCmdSign returns new initialized instance of 'sign' sub command.
func NewCmdSign(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewSignOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "sign SECRETID SECRETKEY",
		DisableFlagsInUseLine: true,
		Short:                 "Sign a jwt token with given secretId and secretKey",
		Long:                  "Sign a jwt token with given secretId and secretKey",
		Example:               signExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		Aliases:    []string{},
		SuggestFor: []string{},
	}

	cmd.Flags().DurationVar(&o.Timeout, "timeout", o.Timeout, "JWT token expires time")
	cmd.Flags().StringVar(&o.Algorithm, "algorithm", o.Algorithm, "Signing algorithm - possible values are HS256, HS384, HS512")

	return cmd
}

// Complete completes all the required options.
func (o *SignOptions) Complete(f cmdutil.Factory, cmd *cobra.Command) error {
	return nil
}

// Validate makes sure there is no discrepency in command options.
func (o *SignOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return cmdutil.UsageErrorf(cmd, "Unexpected args: %v", args)
	}

	switch o.Algorithm {
	case "HS256", "HS384", "HS512":
	default:
		return ErrSigningMethod
	}

	return nil
}

// Run executes a sign sub command using the specified options.
func (o *SignOptions) Run(args []string) error {
	tokenString, err := createJWTToken(o.Algorithm, o.Timeout, args[0], args[1])
	if err != nil {
		return err
	}

	fmt.Fprintf(o.Out, tokenString+"\n")

	return nil
}

// createJWTToken create a jwt token by the given parameters.
func createJWTToken(algorithm string, timeout time.Duration, secretID, secretKey string) (string, error) {
	expire := time.Now().Add(timeout)

	// The token content.
	token := jwt.NewWithClaims(jwt.GetSigningMethod(algorithm), jwt.MapClaims{
		"jti": secretID,
		"exp": expire.Unix(),
		"iat": time.Now().Unix(),
	})

	// Sign the token with the specified secret.
	return token.SignedString([]byte(secretKey))
}
