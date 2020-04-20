package cluster

import (
	core "github.com/DataDog/datadog-agent/pkg/collector/corechecks"
	"context"
	"gopkg.in/yaml.v2"
	"github.com/DataDog/datadog-agent/pkg/collector/check"
	"github.com/DataDog/datadog-agent/pkg/autodiscovery/integration"
	"github.com/DataDog/datadog-agent/pkg/aggregator"
	kubestatemetrics "github.com/DataDog/datadog-agent/pkg/kubestatemetrics/builder"
	"k8s.io/client-go/tools/cache"
	"github.com/DataDog/datadog-agent/pkg/util/kubernetes/apiserver"
	"k8s.io/kube-state-metrics/pkg/options"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	ksmstore "github.com/DataDog/datadog-agent/pkg/kubestatemetrics/store"
	"k8s.io/kube-state-metrics/pkg/allowdenylist"
	"time"
	"k8s.io/apimachinery/pkg/types"
	"github.com/clamoriniere/ddksm/pkg/store"
	"strings"
)

const (
	kubeStateMetricsCheckName = "kube-state-metrics"
)

type KSMConfig struct {
	// TODO fill in all the configurations.
	Collectors                             []string  `yaml:"collectors"`
	//Namespaces                           kubestatemetrics.NamespaceList `yaml:"collectors"`
	//Shard                                int32
	//TotalShards                          int
	//Pod                                  string
	//Namespace                            string
	//MetricBlacklist                      kubestatemetrics.MetricSet
	MetricWhitelist                        []string  `yaml:"metrics"`
	//Version                              bool
	//DisablePodNonGenericResourceMetrics  bool
	//DisableNodeNonGenericResourceMetrics bool
}

type KSMCheck struct {
	ac       *apiserver.APIClient
	core.CheckBase
	instance *KSMConfig
	store []cache.Store
	cancelF context.CancelFunc
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

	builder := kubestatemetrics.New()

	if err := builder.WithEnabledResources(k.instance.Collectors); err != nil {
		log.Errorf("Failed to set up collectors: %v", err)
		return nil
	}
	// All namespaces
	builder.WithNamespaces(options.DefaultNamespaces)

	// Metrics exclusion/inclusion
	allowDenyList, err := allowdenylist.New(options.MetricSet{}, options.MetricSet{})
	if err != nil {
		log.Errorf("Error %v", err)
		return nil
	}

	err = allowDenyList.Parse()
	if err != nil {
		log.Errorf("error initializing the allowDenyList list : %v", err)
		return nil
	}
	builder.WithAllowDenyList(allowDenyList)

	c, err := apiserver.GetAPIClient()
	if err != nil {
		return err
	}

	builder.WithKubeClient(c.Cl)
	builder.WithContext(context.Background())

	builder.WithResync(30 * time.Second)  //k.instance.ResyncPeriod ? Even necessary ?

	builder.WithGenerateStoreFunc(builder.GenerateStore)

	k.store = builder.Build()

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

		metrics =  store.(*ksmstore.MetricsStore).Push()

		processMetrics(sender, metrics)
		// TODO identify how I can extrac tthe store name to convert later on.
	}
	return nil
}

func processMetrics(sender aggregator.Sender, metrics map[types.UID][]store.DDMetricsFam) {
	for u := range metrics {
		for _, mfam := range metrics[u] {
			for _, m := range mfam. {
			}
		}
	}
	for name, metric := range metrics {
		for _, m := range metric {
			sender.Gauge(strings.Replace(name, "_", ".", -1), m.Val, "", m.Labels)
		}
	}
}

func KubeStateMEtricsFactory() check.Check {
	return newKSMCheck(core.NewCheckBase(kubeStateMetricsCheckName), &KSMConfig{})
}

func newKSMCheck(base core.CheckBase, instance *KSMConfig) *KSMCheck {
	return &KSMCheck{
		CheckBase: base,
		instance: instance,
	}
}

func init() {
	// create the KSM builder
	core.RegisterCheck(kubeStateMetricsCheckName, KubeStateMEtricsFactory)
}
