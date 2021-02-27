// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"github.com/marmotedu/marmotedu-sdk-go/marmotedu/service/iam"
	restclient "github.com/marmotedu/marmotedu-sdk-go/rest"

	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
)

// Factory provides abstractions that allow the IAM command to be extended across multiple types
// of resources and different API sets.
// The rings are here for a reason. In order for composers to be able to provide alternative factory implementations
// they need to provide low level pieces of *certain* functions so that when the factory calls back into itself
// it uses the custom version of the function. Rather than try to enumerate everything that someone would want to
// override
// we split the factory into rings, where each ring can depend on methods in an earlier ring, but cannot depend
// upon peer methods in its own ring.
// TODO: make the functions interfaces
// TODO: pass the various interfaces on the factory directly into the command constructors (so the
// commands are decoupled from the factory).
type Factory interface {
	genericclioptions.RESTClientGetter

	// IAMClient gives you back an external iamclient
	IAMClient() (*iam.IamClient, error)

	// Returns a RESTClient for accessing IAM resources or an error.
	RESTClient() (*restclient.RESTClient, error)
}
