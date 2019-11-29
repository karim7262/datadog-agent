package aggregator

import (
	"github.com/DataDog/datadog-agent/pkg/metrics"
	"github.com/segmentio/fasthash/fnv1a"
)

type shardedTimeSampler struct {
	dispatcher   chan []*metrics.MetricSample
	timeSamplers []*TimeSampler
}

func newShardedTimeSampler(shardCount int, interval int64) *shardedTimeSampler {
	var samplers []*TimeSampler
	for i := 0; i < shardCount; i++ {
		samplers = append(samplers, NewTimeSampler(interval))
	}
	s := &shardedTimeSampler{timeSamplers: samplers}
	go s.dispatchLoop()
	return s
}

func (s *shardedTimeSampler) addSample(sample *metrics.MetricSample) {
	s.dispatcher <- []*metrics.MetricSample{sample}
}

func (s *shardedTimeSampler) addSamples(samples []*metrics.MetricSample) {
	s.dispatcher <- samples
}

func (s *shardedTimeSampler) dispatchLoop() {
	for samples := range s.dispatcher {
		dispatchedSamples := make([][]*metrics.MetricSample, len(s.timeSamplers))
		for _, sample := range samples {
			samplerIndex := fnv1a.HashString32(sample.Name) % uint32(len(s.timeSamplers))
			dispatchedSamples[samplerIndex] = append(dispatchedSamples[samplerIndex], sample)
		}
		for i := 0; i < len(dispatchedSamples); i++ {
			s.timeSamplers[i].sampleChann <- dispatchedSamples[i]
		}
	}
}
