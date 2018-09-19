// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

// +build kubeapiserver

package custommetrics

import (
	//"fmt"
	//"os"
	//"time"

	//"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/cmd/server"

	//apimeta "k8s.io/apimachinery/pkg/api/meta"
	//"k8s.io/client-go/discovery"
	//"k8s.io/client-go/dynamic"
	//"k8s.io/client-go/rest"

	"github.com/DataDog/datadog-agent/pkg/clusteragent/custommetrics"
	as "github.com/DataDog/datadog-agent/pkg/util/kubernetes/apiserver"
	"github.com/DataDog/datadog-agent/pkg/util/kubernetes/apiserver/common"
	basecmd "github.com/CharlyF/custom-metrics-apiserver/pkg/cmd"
	"github.com/CharlyF/custom-metrics-apiserver/pkg/provider"
	"github.com/golang/glog"
	"github.com/prometheus/common/log"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"github.com/CharlyF/custom-metrics-apiserver/pkg/apiserver"
	"k8s.io/apimachinery/pkg/runtime"
	"net"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	//"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/wait"
	"flag"
	"os"
)

//var opts *server.CustomMetricsAdapterServerOptions
var cmd *DatadogMetricsAdapter

var stopCh chan struct{}

func init() {
//	// FIXME: log to seelog
	//opts = server.NewCustomMetricsAdapterServerOptions()
}

type DatadogMetricsAdapter struct {
	basecmd.AdapterBase

	// the message printed on startup
	Message string
}

// AddFlags ensures the required flags exist
//func AddFlags(fs *pflag.FlagSet) {
//
//	options.SecureServing.AddFlags(fs)
//	options.Authentication.AddFlags(fs)
//	options.Authorization.AddFlags(fs)
//	options.Features.AddFlags(fs)
//}

//// ValidateArgs validates the custom metrics arguments passed
//func ValidateArgs(args []string) error {
//	return options.Validate(args)
//}

// StartServer creates and start a k8s custom metrics API server
func StartServer() error {

	cmd = &DatadogMetricsAdapter{}
	cmd.Flags().StringVar(&cmd.Message, "msg", "starting adapter...", "startup message")
	cmd.Flags().AddGoFlagSet(flag.CommandLine) // make sure you get the glog flags
	cmd.Flags().Parse(os.Args)

	provider := cmd.makeProviderOrDie()
	cmd.WithExternalMetrics(provider)
	log.Infof("secure serving is %#v", cmd.SecureServing)
	glog.Infof(cmd.Message)

	cmd.Name = "datadog-custom-metrics-adapter"

	conf, err := cmd.Config()
	if err != nil 																																																	{
		log.Infof("err %#v", err)
	}
	server, err := conf.Complete(nil).New("datadog-custom-metrics-adapter", nil, provider)
	if err != nil {
		log.Infof("err server %#v", err)
	}
	server.GenericAPIServer.PrepareRun().Run(wait.NeverStop)

	return nil
	////
	//test := cmd.CustomMetricsAdapterServerOptions.SecureServing
	//log.Infof("adapter init %#v", test)
	//fs := cmd.Flags()
	//args := os.Args
	//log.Debugf("os args is %#v", args)
	//fs.Parse(args) // What does this do ?
	//
	//log.Infof("adapter installed flags %#v", cmd)
	////NewDelegatingAuthenticationOptions()
	////auth := options.NewDelegatingAuthenticationOptions()
	//
	////auth.RequestHeader =
	//_, e := cmd.Authentication.ToAuthenticationConfig()
	//cmd.SecureServing.Validate()
	//
	//if e != nil {
	//	log.Infof("err while authenticating %#v", e)
	//}
	////log.Infof("Authenticated with %#v", a.RequestHeaderConfig)
	//errL := cmd.Authorization.Validate()
	//for _, e := range errL {
	//	log.Infof("err: %v", e)
	//}
	//
	//provider := cmd.makeProviderOrDie()
	//cmd.WithExternalMetrics(provider)
	//
	//config, _ := cmd.Config()
//	s := options.NewSecureServingOptions()
	//
	////config.GenericConfig.SecureServing
	////options.NewRecommendedOptions("foo", runtime.Codec)
	////rec := genericapiserver.NewRecommendedConfig(apiserver.Codecs)
	//
	//server, _ := config.Complete(nil).New("datadog-custom-metrics-adapter", nil, provider)
	//stopCh = make(chan struct{})
	//return server.GenericAPIServer.PrepareRun().Run(stopCh)

	///
	//if err := cmd.Run(wait.NeverStop); err != nil {
		//log.Errorf("unable to run custom metrics adapter: %v", err)
		//return err
	//}
    //return nil


	//config, err := options.Config()
	//if err != nil {
	//	return err
	//}
	//var clientConfig *rest.Config
	//clientConfig, err = rest.InClusterConfig()
	//if err != nil {
	//	return err
	//}
	//
	//discoveryClient, err := discovery.NewDiscoveryClientForConfig(clientConfig)
	//if err != nil {
	//	return fmt.Errorf("unable to construct discovery client for dynamic client: %v", err)
	//}

	//dynamicMapper, err := dynamicmapper.NewRESTMapper(discoveryClient, apimeta.InterfacesForUnstructured, time.Second*5)
	//if err != nil {
	//	return fmt.Errorf("unable to construct dynamic discovery mapper: %v", err)
	//}
	//
	//clientPool := dynamic.NewClientPool(clientConfig, dynamicMapper, dynamic.LegacyAPIPathResolverFunc)
	//if err != nil {
	//	return fmt.Errorf("unable to construct lister client to initialize provider: %v", err)
	//}
	//
	//client, err := as.GetAPIClient()
	//if err != nil {
	//	return err
	//}
	//datadogHPAConfigMap := custommetrics.GetConfigmapName()
	//store, err := custommetrics.NewConfigMapStore(client.Cl, common.GetResourcesNamespace(), datadogHPAConfigMap)
	//if err != nil {
	//	return err
	//}
	//emProvider := custommetrics.NewDatadogProvider(clientPool, dynamicMapper, store)
	//// As the Custom Metrics Provider is introduced, change the first emProvider to a cmProvider.
	//server, err := config.Complete().New("datadog-custom-metrics-adapter", emProvider, emProvider)
	//if err != nil {
	//	return err
	//}
	//stopCh = make(chan struct{})
	//return server.GenericAPIServer.PrepareRun().Run(stopCh)
}

func (a *DatadogMetricsAdapter) makeProviderOrDie() provider.ExternalMetricsProvider {
	client, err := a.DynamicClient()


	if err != nil {
		glog.Fatalf("unable to construct dynamic client: %v", err)
	}
	apiCl, err := as.GetAPIClient()
	if err != nil {
		return nil
	}
	datadogHPAConfigMap := custommetrics.GetConfigmapName()
	store, err := custommetrics.NewConfigMapStore(apiCl.Cl, common.GetResourcesNamespace(), datadogHPAConfigMap)
	if err != nil {
		return nil
	}

	mapper, err := a.RESTMapper()
	if err != nil {
		glog.Fatalf("unable to construct discovery REST mapper: %v", err)
	}

	return custommetrics.NewDatadogProvider(client, mapper, store)
}

func (o DatadogMetricsAdapter) Config() (*apiserver.Config, error) {
	if err := o.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		log.Infof("fail maybedefault %#v", err)
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}
	//nso := options.NewSecureServingOptions()

	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	serverConfig := genericapiserver.NewConfig(codecs)
	//serverConfig.SecureServing
	//o.CustomMetricsAdapterServerOptions.SecureServing
	log.Infof("serverConfig %#v", serverConfig)
	err := o.SecureServing.ApplyTo(&serverConfig.SecureServing, &serverConfig.LoopbackClientConfig)
	if err != nil {
		log.Infof("conversion %#v", err)
	}
	if err := o.Authentication.ApplyTo(&serverConfig.Authentication, serverConfig.SecureServing, nil); err != nil {
		log.Infof("fail AUTHN %#v", err)
		return nil, err
	}
	if err := o.Authorization.ApplyTo(&serverConfig.Authorization); err != nil {
		log.Infof("fail AUTHZ %#v", err)
		return nil, err
	}

	return &apiserver.Config{
		GenericConfig:  serverConfig,
	}, nil
}

// StopServer closes the connection and the server
// stops listening to new commands.
func StopServer() {
	if stopCh != nil {
		close(stopCh)
	}
}
