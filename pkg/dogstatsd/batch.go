package dogstatsd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/metrics"
)

// batcher batches multiple metrics before submission
// this struct is not safe for concurrent use
type batcher struct {
	samples      []metrics.MetricSample
	samplesCount int

	events        []*metrics.Event
	serviceChecks []*metrics.ServiceCheck

	// output channels
	choutSamples       chan<- []metrics.MetricSample
	choutEvents        chan<- []*metrics.Event
	choutServiceChecks chan<- []*metrics.ServiceCheck

	chFile chan []metrics.MetricSample

	metricSamplePool *metrics.MetricSamplePool
}

func newBatcher(agg *aggregator.BufferedAggregator) *batcher {
	s, e, sc := agg.GetBufferedChannels()

	chFile := make(chan []metrics.MetricSample)

	go writeToFile(chFile)

	return &batcher{
		samples:            agg.MetricSamplePool.GetBatch(),
		metricSamplePool:   agg.MetricSamplePool,
		choutSamples:       s,
		choutEvents:        e,
		choutServiceChecks: sc,
		// XXX(remy):
		chFile: chFile,
	}
}

func writeToFile(chFile chan []metrics.MetricSample) {
	f, err := ioutil.TempFile("", "dsd")
	if err != nil {
		log.Fatalf("writeToFile: can't open a tmp file: %v", err)
	}

	select {
	case samples := <-chFile:
		for _, sample := range samples {
			fmt.Fprintf(f, fmt.Sprintf("%+v", sample))
		}
	}
}

func (b *batcher) appendSample(sample metrics.MetricSample) {
	if b.samplesCount == len(b.samples) {
		b.flushSamples()
	}
	b.samples[b.samplesCount] = sample
	b.samplesCount++
}

func (b *batcher) appendEvent(event *metrics.Event) {
	b.events = append(b.events, event)
}

func (b *batcher) appendServiceCheck(serviceCheck *metrics.ServiceCheck) {
	b.serviceChecks = append(b.serviceChecks, serviceCheck)
}

func (b *batcher) flushSamples() {
	if b.samplesCount > 0 {
		// b.choutSamples <- b.samples[:b.samplesCount]
		b.chFile <- b.samples[:b.samplesCount]
		b.samplesCount = 0
		b.samples = b.metricSamplePool.GetBatch()
	}
}

func (b *batcher) flush() {
	b.flushSamples()
	if len(b.events) > 0 {
		b.choutEvents <- b.events
		b.events = []*metrics.Event{}
	}
	if len(b.serviceChecks) > 0 {
		b.choutServiceChecks <- b.serviceChecks
		b.serviceChecks = []*metrics.ServiceCheck{}
	}
}
