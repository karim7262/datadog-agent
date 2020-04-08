package cluster

import (
	core "github.com/DataDog/datadog-agent/pkg/collector/corechecks"
	"context"
	"gopkg.in/yaml.v2"
	"github.com/DataDog/datadog-agent/pkg/collector/check"
	"github.com/DataDog/datadog-agent/pkg/autodiscovery/integration"
	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/kubestatemetrics"
	"k8s.io/kube-state-metrics/pkg/whiteblacklist"
	"k8s.io/client-go/tools/cache"
	"github.com/DataDog/datadog-agent/pkg/util/kubernetes/apiserver"
	"github.com/DataDog/datadog-agent/pkg/util/log"

	"k8s.io/klog"
	"k8s.io/kube-state-metrics/pkg/options"
)

const (
	kubeStateMetricsCheckName = "kube-state-metrics"
)

type KSMConfig struct {
	// TODO fill in all the configurations.
	Collectors                           options.CollectorSet  `yaml:"collectors"`
	//Namespaces                           kubestatemetrics.NamespaceList `yaml:"collectors"`
	//Shard                                int32
	//TotalShards                          int
	//Pod                                  string
	//Namespace                            string
	//MetricBlacklist                      kubestatemetrics.MetricSet
	MetricWhitelist                      options.MetricSet
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
	err := k.CommonConfigure(config, source)
	if err != nil {
		return err
	}

	err = k.instance.parse(config)
	if err != nil {
		log.Error("could not parse the config for the API server")
		return err
	}

	whiteBlackList, err := whiteblacklist.New(k.instance.MetricWhitelist, nil) // k.instance.MetricBlacklist)
	if err != nil {
		klog.Fatal(err)
	}
	err = whiteBlackList.Parse()
	if err != nil {
		log.Errorf("error initializing the whiteblack list : %v", err)
		return err
	}
	k.builder.WithWhiteBlackList(whiteBlackList)

	var collectors []string
	if len(k.instance.Collectors) == 0 {
	collectors = options.DefaultCollectors.AsSlice()
	} else {
		collectors = k.instance.Collectors.AsSlice()
	}
	if err := k.builder.WithEnabledResources(collectors); err != nil {
		log.Errorf("Failed to set up collectors: %v", err)
		return err
	}
	log.Infof("KSM configured with %s", collectors)

	k.builder.WithNamespaces(options.DefaultNamespaces)

	k.store = k.builder.Build()

	return nil
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
		store.(*kubestatemetrics.MetricsStore).Push(sender)
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

	//whiteBlackList, err := whiteblacklist.New(instance.MetricWhitelist, instance.MetricBlacklist)

	//builder.WithWhiteBlackList(instance.MetricWhitelist)
	// We start the API Server Client.
	ac, err := apiserver.GetAPIClient()
	if err != nil {
		return nil
	}

	builder.WithKubeClient(ac.Cl)
	builder.WithVPAClient(nil)

	builder.WithContext(ctx)
	builder.WithCustomGenerateStoreFunc(builder.GenerateStore)

	ksmCheck := &KSMCheck{
		CheckBase: base,
		ac: ac,
		builder: builder,
	}

	return ksmCheck
}

func init() {
	// create the KSM builder
	core.RegisterCheck(kubeStateMetricsCheckName, KubeStateMEtricsFactory)
}
