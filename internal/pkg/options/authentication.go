/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"github.com/spf13/pflag"
)

// ClientCertAuthenticationOptions provides different options for client cert auth.
type ClientCertAuthenticationOptions struct {
	// ClientCA is the certificate bundle for all the signers that you'll recognize for incoming client certificates
	ClientCA string `json:"client-ca-file" mapstructure:"client-ca-file"`
}

// NewClientCertAuthenticationOptions creates a ClientCertAuthenticationOptions object with default parameters.
func NewClientCertAuthenticationOptions() *ClientCertAuthenticationOptions {
	return &ClientCertAuthenticationOptions{
		ClientCA: "",
	}
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (o *ClientCertAuthenticationOptions) Validate() []error {
	return []error{}
}

// AddFlags adds flags related to ClientCertAuthenticationOptions for a specific server to the
// specified FlagSet.
func (o *ClientCertAuthenticationOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ClientCA, "client-ca-file", o.ClientCA, ""+
		"If set, any request presenting a client certificate signed by one of "+
		"the authorities in the client-ca-file is authenticated with an identity "+
		"corresponding to the CommonName of the client certificate.")
}
