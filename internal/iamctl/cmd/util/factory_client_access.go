// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// this file contains factories with no other dependencies

package util

import (
	"github.com/marmotedu/iam/pkg/cli/genericclioptions"
	"github.com/marmotedu/marmotedu-sdk-go/marmotedu"
	restclient "github.com/marmotedu/marmotedu-sdk-go/rest"
	"github.com/marmotedu/marmotedu-sdk-go/tools/clientcmd"
)

type factoryImpl struct {
	clientGetter genericclioptions.RESTClientGetter
}

func NewFactory(clientGetter genericclioptions.RESTClientGetter) Factory {
	if clientGetter == nil {
		panic("attempt to instantiate client_access_factory with nil clientGetter")
	}

	f := &factoryImpl{
		clientGetter: clientGetter,
	}

	return f
}

func (f *factoryImpl) ToRESTConfig() (*restclient.Config, error) {
	return f.clientGetter.ToRESTConfig()
}

func (f *factoryImpl) ToRawIAMConfigLoader() clientcmd.ClientConfig {
	return f.clientGetter.ToRawIAMConfigLoader()
}

func (f *factoryImpl) IAMClientSet() (*marmotedu.Clientset, error) {
	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return marmotedu.NewForConfig(clientConfig)
}

func (f *factoryImpl) RESTClient() (*restclient.RESTClient, error) {
	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	setIAMDefaults(clientConfig)
	return restclient.RESTClientFor(clientConfig)
}
