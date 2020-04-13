package cluster

import (
	core "github.com/DataDog/datadog-agent/pkg/collector/corechecks"
	"context"
	"gopkg.in/yaml.v2"
	"github.com/DataDog/datadog-agent/pkg/collector/check"
	"github.com/DataDog/datadog-agent/pkg/autodiscovery/integration"
	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/kubestatemetrics"
	"k8s.io/client-go/tools/cache"
	"github.com/DataDog/datadog-agent/pkg/util/kubernetes/apiserver"
	"k8s.io/kube-state-metrics/pkg/allowdenylist"
	"k8s.io/kube-state-metrics/pkg/options"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	ksmstore "github.com/DataDog/datadog-agent/pkg/store"
)

const (
	kubeStateMetricsCheckName = "kube-state-metrics"
)

type KSMConfig struct {
	// TODO fill in all the configurations.
	Collectors                             []string  `yaml:"collectors"` //options.CollectorSet
	//Namespaces                           kubestatemetrics.NamespaceList `yaml:"collectors"`
	//Shard                                int32
	//TotalShards                          int
	//Pod                                  string
	//Namespace                            string
	//MetricBlacklist                      kubestatemetrics.MetricSet
	MetricWhitelist                        []string  `yaml:"metrics"` //options.MetricSet
	//Version                              bool
	//DisablePodNonGenericResourceMetrics  bool
	//DisableNodeNonGenericResourceMetrics bool
}

type KSMCheck struct {
	ac       *apiserver.APIClient
	core.CheckBase
	instance *KSMConfig
	builder *kubestatemetrics.Builder
	store []cache.Store
}

func (k *KSMCheck) Configure(config, initConfig integration.Data, source string) error {
	return nil
	//
	//err := k.CommonConfigure(config, source)
	//if err != nil {
	//	return err
	//}
	//err = k.instance.parse(config)
	//if err != nil {
	//	log.Error("could not parse the config for the API server")
	//	return err
	//}
	//
	//whiteBlackList, err := whiteblacklist.New(options.MetricSet{}, nil) // k.instance.MetricBlacklist)
	//if err != nil {
	//	klog.Fatal(err)
	//}
	//whiteBlackList.Include([]string{"kube_node_*"})
	//err = whiteBlackList.Parse()
	//if err != nil {
	//	log.Errorf("error initializing the whiteblack list : %v", err)
	//	return err
	//}
	//log.Infof("All metrics are %v ", whiteBlackList.Status())
	//k.builder.WithWhiteBlackList(whiteBlackList)
	//
	//var collectors []string
	//if len(k.instance.Collectors) == 0 {
	//collectors = options.DefaultCollectors.AsSlice()
	//} else {
	//	collectors = k.instance.Collectors
	//}
	//if err := k.builder.WithEnabledResources(collectors); err != nil {
	//	log.Errorf("Failed to set up collectors: %v", err)
	//	return err
	//}
	//log.Infof("KSM configured with %s", collectors)
	//
	//k.builder.WithNamespaces(options.DefaultNamespaces)
	//
	//k.store = k.builder.Build()
	//
	//for _, store := range k.store {
	//	store.(*kubestatemetrics.MetricsStore).List()
	//}
	//
	//log.Infof("k store is %#V", k.store)
	//
	//return nil
}


func (c *KSMConfig) parse(data []byte) error {
	// default values

	return yaml.Unmarshal(data, c)
}

func (k *KSMCheck) Run() error {
	sender, err := aggregator.GetSender(k.ID())
	if err != nil {
		return err
	}
	defer sender.Commit()

	for _, store := range k.store {
		store.(*ksmstore.MetricsStore).Push(sender)
	}

	return nil
}

func KubeStateMEtricsFactory() check.Check {
	return newKSMCheck(core.NewCheckBase(kubeStateMetricsCheckName), &KSMConfig{})
}

func newKSMCheck(base core.CheckBase, instance *KSMConfig) *KSMCheck {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	builder := kubestatemetrics.New()

	//enabling collectors first ?
	if err := builder.WithEnabledResources([]string{"nodes"}); err != nil {
		log.Errorf("Failed to set up collectors: %v", err)
		return nil
	}
	// All namespaces
	builder.WithNamespaces(options.DefaultNamespaces)

	// Metrics exclusion/inclusion
	allowDenyList, err := allowdenylist.New(options.MetricSet{}, nil)
	if err != nil {
		log.Errorf("Error %v", err)
		return nil
	}
	builder.WithAllowDenyList(allowDenyList)

	allowDenyList.Exclude([]string{"kube_node_status_allocatable_cpu_cores"})

	err = allowDenyList.Parse()
	if err != nil {
		log.Errorf("error initializing the allowDenyList list : %v", err)
		return nil
	}
	builder.WithAllowDenyList(allowDenyList)


	// We start the API Server Client.
	ac, err := apiserver.GetAPIClient()
	if err != nil {
		return nil
	}
	v, err := ac.Cl.Discovery().ServerVersion()
	if err != nil {
		log.Errorf("Could not get server version %v",err)
		return nil
	}
	log.Infof("Connected to Server res %s", v.String())
//	vpaClient, err := vpaclientset.NewForConfig(config)

	builder.WithKubeClient(ac.Cl)
	builder.WithContext(ctx)
	//builder.WithVPAClient(nil)
	// Key part ?
	builder.WithGenerateStoreFunc(builder.GenerateStore)

	return &KSMCheck{
		CheckBase: base,
		ac: ac,
		instance: instance,
		//store: store,
		builder: builder,
	}
}

func init() {
	// create the KSM builder
	core.RegisterCheck(kubeStateMetricsCheckName, KubeStateMEtricsFactory)
}
