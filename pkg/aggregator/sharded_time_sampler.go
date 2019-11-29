package aggregator

import (
	"github.com/DataDog/datadog-agent/pkg/metrics"
	"github.com/segmentio/fasthash/fnv1a"
)

type shardedTimeSampler struct {
	timeSamplers []*TimeSampler
}

func newShardedTimeSampler(shardCount int, interval int64) *shardedTimeSampler {
	var samplers []*TimeSampler
	for i := 0; i < shardCount; i++ {
		samplers = append(samplers, NewTimeSampler(interval))
	}
	return &shardedTimeSampler{timeSamplers: samplers}
}

func (s *shardedTimeSampler) addSample(sample *metrics.MetricSample) {
	samplerIndex := fnv1a.HashString32(sample.Name) % uint32(len(s.timeSamplers))
	s.timeSamplers[samplerIndex].sampleChann <- sample
}
