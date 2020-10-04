// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package apiserver does all of the work necessary to create a iam APIServer.
package apiserver

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	pb "github.com/marmotedu/api/proto/apiserver/v1"
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
	"github.com/marmotedu/component-base/pkg/cli/globalflag"
	"github.com/marmotedu/component-base/pkg/term"
	"github.com/marmotedu/component-base/pkg/util/idutil"
	"github.com/marmotedu/component-base/pkg/version"
	"github.com/marmotedu/component-base/pkg/version/verflag"
	"github.com/marmotedu/errors"
	cachev1 "github.com/marmotedu/iam/internal/apiserver/api/v1/cache"
	"github.com/marmotedu/iam/internal/apiserver/options"
	"github.com/marmotedu/iam/internal/apiserver/store"
	"github.com/marmotedu/iam/internal/apiserver/store/datastore"
	genericoptions "github.com/marmotedu/iam/internal/pkg/options"
	genericapiserver "github.com/marmotedu/iam/internal/pkg/server"
	"github.com/marmotedu/iam/pkg/storage"
	"github.com/marmotedu/log"
)

const (
	// recommendedFileName defines the configuration used by iam-apiserver.
	// the configuration file is different from other iam service.
	recommendedFileName = "iam-apiserver.yaml"

	// appName defines the executable binary filename for iam-apiserver component.
	appName = "iam-apiserver"
)

// NewAPIServerCommand creates a *cobra.Command object with default parameters.
func NewAPIServerCommand() *cobra.Command {
	cliflag.InitFlags()

	s := options.NewServerRunOptions()

	cmd := &cobra.Command{
		Use:   appName,
		Short: "The IAM API server to validates and configures data for the api objects",
		Long: `The IAM API server validates and configures data
for the api objects which include users, policies, secrets, and
others. The API Server services REST operations to do the api objects management.

Find more iam-apiserver information at:
    https://github.com/marmotedu/iam/blob/master/docs/admin/iam-apiserver.md`,

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

			log.Init(completedOptions.Log)
			defer log.Flush()

			return Run(completedOptions, genericapiserver.SetupSignalHandler())
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return nil
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

	namedFlagSets := s.Flags()
	verflag.AddFlags(namedFlagSets.FlagSet("global"))
	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())
	fs := cmd.Flags()
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

// Run runs the specified APIServer. This should never exit.
func Run(completedOptions completedServerRunOptions, stopCh <-chan struct{}) error {
	// To help debugging, immediately log config and version
	log.Infof("config: `%s`", completedOptions.String())
	log.Infof("version: %+v", version.Get().ToJSON())

	if err := completedOptions.InitDataStore(); err != nil {
		return err
	}

	serverConfig, err := createAPIServerConfig(completedOptions.ServerRunOptions)
	if err != nil {
		return err
	}

	server, err := createAPIServer(serverConfig)
	if err != nil {
		return err
	}

	return server.Run(stopCh)
}

// ExtraConfig defines extra configuration for the iam-apiserver.
type ExtraConfig struct {
	Addr       string
	ServerCert genericoptions.GeneratableKeyCert
}

type completedExtraConfig struct {
	*ExtraConfig
}

// apiServerConfig defines configuration for the iam-apiserver.
type apiServerConfig struct {
	GenericConfig *genericapiserver.Config
	ExtraConfig   ExtraConfig
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   completedExtraConfig
}

// APIServer is only responsible for serving the APIs for iam-apiserver.
type APIServer struct {
	GrpcAPIServer    *grpcAPIServer
	GenericAPIServer *genericapiserver.GenericAPIServer
}

// Complete fills in any fields not set that are required to have valid data and can be derived from other fields.
func (c *ExtraConfig) Complete() completedExtraConfig {
	if c.Addr == "" {
		c.Addr = "127.0.0.1:8081"
	}

	return completedExtraConfig{c}
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *apiServerConfig) Complete() completedConfig {
	return completedConfig{
		c.GenericConfig.Complete(),
		c.ExtraConfig.Complete(),
	}
}

// New returns a new instance of APIServer from the given config.
// Certain config fields will be set to a default value if unset.
func (c completedConfig) New() (*APIServer, error) {
	genericServer, err := c.GenericConfig.New()
	if err != nil {
		return nil, err
	}

	initRouter(genericServer.Engine)

	grpcServer := c.ExtraConfig.New()

	s := &APIServer{
		GenericAPIServer: genericServer,
		GrpcAPIServer:    grpcServer,
	}

	return s, nil
}

// New create a grpcAPIServer instance.
func (c *ExtraConfig) New() *grpcAPIServer {
	creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err.Error())
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterCacheServer(grpcServer, &cachev1.Cache{})

	reflection.Register(grpcServer)

	return &grpcAPIServer{grpcServer, c.Addr}
}

// Run start the APIServer.
func (s *APIServer) Run(stopCh <-chan struct{}) error {
	// run grpc server
	go s.GrpcAPIServer.Run(stopCh)

	// run generic server
	return s.GenericAPIServer.Run(stopCh)
}

// createAPIServer create apiserver according to apiserver config.
func createAPIServer(apiServerConfig *apiServerConfig) (*APIServer, error) {
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

// createAPIServerConfig create apiserver config from all options.
func createAPIServerConfig(s *options.ServerRunOptions) (*apiServerConfig, error) {
	genericConfig, err := buildGenericConfig(s)
	if err != nil {
		return nil, err
	}

	config := &apiServerConfig{
		GenericConfig: genericConfig,
		ExtraConfig: ExtraConfig{
			Addr:       fmt.Sprintf("%s:%d", s.GrpcOptions.BindAddress, s.GrpcOptions.BindPort),
			ServerCert: s.SecureServing.ServerCert,
		},
	}

	return config, nil
}

// completedServerRunOptions is a private wrapper that enforces a call of Complete() before Run can be invoked.
type completedServerRunOptions struct {
	*options.ServerRunOptions
}

// Complete set default ServerRunOptions.
// Should be called after iam-apiserver flags parsed.
func Complete(s *options.ServerRunOptions) (completedServerRunOptions, error) {
	var options completedServerRunOptions

	genericapiserver.LoadConfig(s.APIConfig, recommendedFileName)

	if err := viper.Unmarshal(s); err != nil {
		return options, err
	}

	if s.JwtOptions.Key == "" {
		s.JwtOptions.Key = idutil.NewSecretKey()
	}

	if err := s.SecureServing.Complete(); err != nil {
		return options, err
	}

	options.ServerRunOptions = s

	return options, nil
}

func (completedOptions completedServerRunOptions) InitDataStore() error {
	completedOptions.InitRedisStore()

	return completedOptions.InitMySQLStore()
}

func (completedOptions completedServerRunOptions) InitMySQLStore() error {
	mysqlStore, err := datastore.NewMySQLStore(completedOptions.MySQLOptions)
	if err != nil {
		return err
	}

	store.SetClient(mysqlStore)

	return nil
}

func (completedOptions completedServerRunOptions) InitRedisStore() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := &storage.Config{
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

	// try to connect to redis
	go storage.ConnectToRedis(ctx, config)
}
