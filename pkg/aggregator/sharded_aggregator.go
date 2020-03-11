package aggregator

import (
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/metrics"
	"github.com/DataDog/datadog-agent/pkg/serializer"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"github.com/segmentio/fasthash/fnv1a"
)

type ShardedAggregator struct {
	aggs []*BufferedAggregator
}

func InitShardedAggregator(s serializer.MetricSerializer, hostname, agentName string) *ShardedAggregator {
	sa := &ShardedAggregator{}
	n := config.Datadog.GetInt("dogstatsd_sharded_aggregator") // TODO(remy): do that properly
	if n == 0 {
		n = 4
	}
	for i := 0; i < n; i++ {
		sa.aggs = append(sa.aggs, InitAggregatorWithFlushInterval(s, hostname, agentName, DefaultFlushInterval))
	}
	log.Infof("Initialized %d aggregators\n", len(sa.aggs))
	return sa
}

func (a *ShardedAggregator) First() *BufferedAggregator {
	return a.aggs[0]
}

func (a *ShardedAggregator) PushSamples(samples []metrics.MetricSample) {
	// XXX(remy): erroneous because we're only using the first metric name, we could send
	// the same metric name to different aggregators.
	a.aggs[fnv1a.HashString32(samples[0].Name)%uint32(len(a.aggs))].bufferedMetricIn <- samples
}
