// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package genericclioptions

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/marmotedu/marmotedu-sdk-go/rest"
	"github.com/marmotedu/marmotedu-sdk-go/tools/clientcmd"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Defines flag for iamctl.
const (
	FlagIAMConfig     = "iamconfig"
	FlagBearerToken   = "user.token"
	FlagUsername      = "user.username"
	FlagPassword      = "user.password"
	FlagSecretID      = "user.secret-id"
	FlagSecretKey     = "user.secret-key"
	FlagCertFile      = "user.client-certificate"
	FlagKeyFile       = "user.client-key"
	FlagTLSServerName = "server.tls-server-name"
	FlagInsecure      = "server.insecure-skip-tls-verify"
	FlagCAFile        = "server.certificate-authority"
	FlagAPIServer     = "server.address"
	FlagTimeout       = "server.timeout"
	FlagMaxRetries    = "server.max-retries"
	FlagRetryInterval = "server.retry-interval"
)

// RESTClientGetter is an interface that the ConfigFlags describe to provide an easier way to mock for commands
// and eliminate the direct coupling to a struct type.  Users may wish to duplicate this type in their own packages
// as per the golang type overlapping.
type RESTClientGetter interface {
	// ToRESTConfig returns restconfig
	ToRESTConfig() (*rest.Config, error)
	// ToRawIAMConfigLoader return iamconfig loader as-is
	ToRawIAMConfigLoader() clientcmd.ClientConfig
}

var _ RESTClientGetter = &ConfigFlags{}

// ConfigFlags composes the set of values necessary
// for obtaining a REST client config.
type ConfigFlags struct {
	IAMConfig *string

	BearerToken *string
	Username    *string
	Password    *string
	SecretID    *string
	SecretKey   *string

	Insecure      *bool
	TLSServerName *string
	CertFile      *string
	KeyFile       *string
	CAFile        *string

	APIServer     *string
	Timeout       *time.Duration
	MaxRetries    *int
	RetryInterval *time.Duration

	clientConfig clientcmd.ClientConfig
	lock         sync.Mutex
	// If set to true, will use persistent client config and
	// propagate the config to the places that need it, rather than
	// loading the config multiple times
	usePersistentConfig bool
}

// ToRESTConfig implements RESTClientGetter.
// Returns a REST client configuration based on a provided path
// to a .iamconfig file, loading rules, and config flag overrides.
// Expects the AddFlags method to have been called.
func (f *ConfigFlags) ToRESTConfig() (*rest.Config, error) {
	return f.ToRawIAMConfigLoader().ClientConfig()
}

// ToRawIAMConfigLoader binds config flag values to config overrides
// Returns an interactive clientConfig if the password flag is enabled,
// or a non-interactive clientConfig otherwise.
func (f *ConfigFlags) ToRawIAMConfigLoader() clientcmd.ClientConfig {
	if f.usePersistentConfig {
		return f.toRawIAMPersistentConfigLoader()
	}

	return f.toRawIAMConfigLoader()
}

func (f *ConfigFlags) toRawIAMConfigLoader() clientcmd.ClientConfig {
	config := clientcmd.NewConfig()
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return clientcmd.NewClientConfigFromConfig(config)
}

// toRawIAMPersistentConfigLoader binds config flag values to config overrides
// Returns a persistent clientConfig for propagation.
func (f *ConfigFlags) toRawIAMPersistentConfigLoader() clientcmd.ClientConfig {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.clientConfig == nil {
		f.clientConfig = f.toRawIAMConfigLoader()
	}

	return f.clientConfig
}

// AddFlags binds client configuration flags to a given flagset.
func (f *ConfigFlags) AddFlags(flags *pflag.FlagSet) {
	if f.IAMConfig != nil {
		flags.StringVar(f.IAMConfig, FlagIAMConfig, *f.IAMConfig,
			fmt.Sprintf("Path to the %s file to use for CLI requests", FlagIAMConfig))
	}

	if f.BearerToken != nil {
		flags.StringVar(
			f.BearerToken,
			FlagBearerToken,
			*f.BearerToken,
			"Bearer token for authentication to the API server",
		)
	}

	if f.Username != nil {
		flags.StringVar(f.Username, FlagUsername, *f.Username, "Username for basic authentication to the API server")
	}

	if f.Password != nil {
		flags.StringVar(f.Password, FlagPassword, *f.Password, "Password for basic authentication to the API server")
	}

	if f.SecretID != nil {
		flags.StringVar(f.SecretID, FlagSecretID, *f.SecretID, "SecretID for JWT authentication to the API server")
	}

	if f.SecretKey != nil {
		flags.StringVar(f.SecretKey, FlagSecretKey, *f.SecretKey, "SecretKey for jwt authentication to the API server")
	}

	if f.CertFile != nil {
		flags.StringVar(f.CertFile, FlagCertFile, *f.CertFile, "Path to a client certificate file for TLS")
	}
	if f.KeyFile != nil {
		flags.StringVar(f.KeyFile, FlagKeyFile, *f.KeyFile, "Path to a client key file for TLS")
	}
	if f.TLSServerName != nil {
		flags.StringVar(f.TLSServerName, FlagTLSServerName, *f.TLSServerName, ""+
			"Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used")
	}
	if f.Insecure != nil {
		flags.BoolVar(f.Insecure, FlagInsecure, *f.Insecure, ""+
			"If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure")
	}
	if f.CAFile != nil {
		flags.StringVar(f.CAFile, FlagCAFile, *f.CAFile, "Path to a cert file for the certificate authority")
	}

	if f.APIServer != nil {
		flags.StringVarP(f.APIServer, FlagAPIServer, "s", *f.APIServer, "The address and port of the IAM API server")
	}

	if f.Timeout != nil {
		flags.DurationVar(
			f.Timeout,
			FlagTimeout,
			*f.Timeout,
			"The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests.",
		)
	}

	if f.MaxRetries != nil {
		flag.IntVar(f.MaxRetries, FlagMaxRetries, *f.MaxRetries, "Maximum number of retries.")
	}

	if f.RetryInterval != nil {
		flags.DurationVar(
			f.RetryInterval,
			FlagRetryInterval,
			*f.RetryInterval,
			"The interval time between each attempt.",
		)
	}
}

// WithDeprecatedPasswordFlag enables the username and password config flags.
func (f *ConfigFlags) WithDeprecatedPasswordFlag() *ConfigFlags {
	f.Username = pointer.ToString("")
	f.Password = pointer.ToString("")

	return f
}

// WithDeprecatedSecretFlag enables the secretID and secretKey config flags.
func (f *ConfigFlags) WithDeprecatedSecretFlag() *ConfigFlags {
	f.SecretID = pointer.ToString("")
	f.SecretKey = pointer.ToString("")

	return f
}

// NewConfigFlags returns ConfigFlags with default values set.
func NewConfigFlags(usePersistentConfig bool) *ConfigFlags {
	return &ConfigFlags{
		IAMConfig: pointer.ToString(""),

		BearerToken:   pointer.ToString(""),
		Insecure:      pointer.ToBool(false),
		TLSServerName: pointer.ToString(""),
		CertFile:      pointer.ToString(""),
		KeyFile:       pointer.ToString(""),
		CAFile:        pointer.ToString(""),

		APIServer:           pointer.ToString(""),
		Timeout:             pointer.ToDuration(30 * time.Second),
		MaxRetries:          pointer.ToInt(0),
		RetryInterval:       pointer.ToDuration(1 * time.Second),
		usePersistentConfig: usePersistentConfig,
	}
}
