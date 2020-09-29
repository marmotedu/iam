// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package authzserver does all of the work necessary to create a authzserver
package authzserver

import (
	"fmt"

	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
	"github.com/marmotedu/component-base/pkg/cli/globalflag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marmotedu/component-base/pkg/term"
	"github.com/marmotedu/component-base/pkg/version"
	"github.com/marmotedu/component-base/pkg/version/verflag"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/internal/authzserver/analytics"
	"github.com/marmotedu/iam/internal/authzserver/options"
	"github.com/marmotedu/iam/internal/authzserver/store"
	genericapiserver "github.com/marmotedu/iam/internal/pkg/server"
	"github.com/marmotedu/iam/pkg/storage"
	"github.com/marmotedu/log"
)

const (
	// recommendedFileName defines the configuration used by iam-authz-server.
	// the configuration file is different from other iam service.
	recommendedFileName = "iam-authz-server.yaml"

	// appName defines the executable binary filename for iam-authz-server component.
	appName = "iam-authz-server"

	// RedisKeyPrefix defines the prefix key in redis for analytics data.
	RedisKeyPrefix = "analytics-"
)

// NewAuthzServerCommand creates a *cobra.Command object with default parameters.
func NewAuthzServerCommand() *cobra.Command {
	cliflag.InitFlags()

	s := options.NewServerRunOptions()

	cmd := &cobra.Command{
		Use:   appName,
		Short: "Authorization server to decide who is able to do what on something given some context",
		Long: `Authorization server to run ladon policies which can protecting your resources.
It is written inspired by AWS IAM policiis.

Find more iam-authz-server information at:
    https://github.com/marmotedu/iam/blob/master/docs/admin/iam-authz-server.md,

Find more ladon information at:
    https://github.com/ory/ladon`,

		// stop printing usage when the command errors
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			verflag.PrintAndExitIfRequested(appName)
			cliflag.PrintFlags(cmd.Flags())

			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			// set default options
			completedOptions, err := Complete(s)
			if err != nil {
				return err
			}

			// validate options
			if errs := completedOptions.Validate(); len(errs) != 0 {
				return errors.NewAggregate(errs)
			}

			// setup logger
			log.Init(completedOptions.Log)
			defer log.Flush()

			return Run(completedOptions, genericapiserver.SetupSignalHandler())
		},
		PostRun: func(cmd *cobra.Command, args []string) {
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	fs := cmd.Flags()
	namedFlagSets := s.Flags()
	verflag.AddFlags(namedFlagSets.FlagSet("global"))
	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})

	return cmd
}

// Run runs the specified AuthzServer. This should never exit.
func Run(completedOptions completedServerRunOptions, stopCh <-chan struct{}) error {
	// To help debugging, immediately log config and version
	log.Infof("config: `%s`", completedOptions.String())
	log.Infof("version: %+v", version.Get().ToJSON())

	// start cacheService
	cacheService := store.New(buildStorageConfig(completedOptions), completedOptions.RPCServer, completedOptions.ClientCA)
	cacheService.Start()

	// start analytics service
	// ey := fmt.Sprintf("%s:%d", audit.Request.Context["username"].(string), time.Now().Unix())]"
	if completedOptions.AnalyticsOptions.Enable {
		analyticsStore := storage.RedisCluster{KeyPrefix: RedisKeyPrefix}
		analytics.NewAnalytics(completedOptions.AnalyticsOptions, &analyticsStore).Start(stopCh)
	}

	// create apiserver config from all options
	apiServerConfig, err := CreateAuthzServerConfig(completedOptions.ServerRunOptions)
	if err != nil {
		return err
	}

	// create apiserver according to apiserver config
	server, err := CreateAuthzServer(apiServerConfig)
	if err != nil {
		return err
	}

	return server.Run(stopCh)
}

// ExtraConfig defines extra configuration for the master.
type ExtraConfig struct {
	SecretKey string
}

type completedExtraConfig struct {
	*ExtraConfig
}

// authzServerConfig defines configuration for the iam-authz-server.
type authzServerConfig struct {
	GenericConfig *genericapiserver.Config
	ExtraConfig   ExtraConfig
}

// AuthzServer is only responsible for serving the APIs for iam-authz-server.
type AuthzServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   completedExtraConfig
}

// Complete fills in any fields not set that are required to have valid data and can be derived from other fields.
func (c *ExtraConfig) Complete() completedExtraConfig {
	return completedExtraConfig{c}
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *authzServerConfig) Complete() completedConfig {
	return completedConfig{
		c.GenericConfig.Complete(),
		c.ExtraConfig.Complete(),
	}
}

// New returns a new instance of AuthzServer from the given config.
// Certain config fields will be set to a default value if unset.
func (c completedConfig) New() (*AuthzServer, error) {
	genericServer, err := c.GenericConfig.New()
	if err != nil {
		return nil, err
	}

	s := &AuthzServer{
		GenericAPIServer: genericServer,
	}

	installHandler(s.GenericAPIServer.Engine)

	return s, nil
}

// Run start to run AuthzServer.
func (s *AuthzServer) Run(stopCh <-chan struct{}) error {
	// run generic server
	return s.GenericAPIServer.Run(stopCh)
}

// CreateAuthzServer create AuthzServer with authzServerConfig.
func CreateAuthzServer(apiServerConfig *authzServerConfig) (*AuthzServer, error) {
	apiServer, err := apiServerConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	return apiServer, nil
}

func buildGenericConfig(s *options.ServerRunOptions) (genericConfig *genericapiserver.Config, lastErr error) {
	genericConfig = genericapiserver.NewConfig()
	if lastErr = s.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = s.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = s.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = s.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}

func buildStorageConfig(completedOptions completedServerRunOptions) *storage.Config {
	return &storage.Config{
		Host:                  completedOptions.RedisOptions.Host,
		Port:                  completedOptions.RedisOptions.Port,
		Addrs:                 completedOptions.RedisOptions.Addrs,
		MasterName:            completedOptions.RedisOptions.MasterName,
		Username:              completedOptions.RedisOptions.Username,
		Password:              completedOptions.RedisOptions.Password,
		Database:              completedOptions.RedisOptions.Database,
		MaxIdle:               completedOptions.RedisOptions.MaxIdle,
		MaxActive:             completedOptions.RedisOptions.MaxActive,
		Timeout:               completedOptions.RedisOptions.Timeout,
		EnableCluster:         completedOptions.RedisOptions.EnableCluster,
		UseSSL:                completedOptions.RedisOptions.UseSSL,
		SSLInsecureSkipVerify: completedOptions.RedisOptions.SSLInsecureSkipVerify,
	}
}

// CreateAuthzServerConfig create authzServerConfig based on options.ServerRunOptions.
func CreateAuthzServerConfig(s *options.ServerRunOptions) (*authzServerConfig, error) {
	genericConfig, err := buildGenericConfig(s)
	if err != nil {
		return nil, err
	}
	// If specified, all requests except those which match the LongRunningFunc predicate will timeout
	// after this duration.
	config := &authzServerConfig{
		GenericConfig: genericConfig,
		ExtraConfig:   ExtraConfig{
			//SecretKey: s.SecretKey,
		},
	}

	return config, nil
}

// completedServerRunOptions is a private wrapper that enforces a call of complete() before Run can be invoked.
type completedServerRunOptions struct {
	*options.ServerRunOptions
}

// Complete set default ServerRunOptions.
// Should be called after authzserver flags parsed.
func Complete(s *options.ServerRunOptions) (completedServerRunOptions, error) {
	var options completedServerRunOptions

	genericapiserver.LoadConfig(s.AuthzConfig, recommendedFileName)

	if err := viper.Unmarshal(s); err != nil {
		return options, err
	}

	if err := s.SecureServing.Complete(); err != nil {
		return options, err
	}

	options.ServerRunOptions = s

	return options, nil
}
