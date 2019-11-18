package dogstatsd

import (
	"sync"

	"github.com/DataDog/datadog-agent/pkg/metrics"
)

// MetricSample is a metric sample originating from DogStatsD
// Structuraly, this is similar to metrics.MetricSample with []byte slices
// instead of strings. Under the hood those []byte slices are pointing to
// memory allocated in the packet they were received from.
type MetricSample struct {
	Name       []byte
	Value      float64
	SetValue   []byte
	MetricType metrics.MetricType
	Tags       [][]byte
	Hostname   []byte
	SampleRate float64
	Timestamp  float64
	pool       *metricSamplePool
}

func (m *MetricSample) release() {
	m.pool.Put(m)
}

type metricSamplePool struct {
	pool *sync.Pool
}

func newSamplePool() *metricSamplePool {
	pool := &metricSamplePool{
		pool: &sync.Pool{},
	}
	pool.pool.New = func() interface{} {
		return &MetricSample{pool: pool}
	}
	return pool
}

func (m *metricSamplePool) Get() *MetricSample {
	return m.pool.Get().(*MetricSample)
}

func (m *metricSamplePool) Put(sample *MetricSample) {
	sample.Tags = sample.Tags[:]
	m.pool.Put(sample)
}

// Event is an event originating from DogStatsD
// Structuraly, this is similar to metrics.Event with []byte slices
// instead of strings. Under the hood those []byte slices are pointing to
// memory allocated in the packet they were received from.
type Event struct {
	Title          []byte
	Text           []byte
	Timestamp      int64
	Priority       metrics.EventPriority
	Hostname       []byte
	Tags           [][]byte
	ExtraTags      []string
	AlertType      metrics.EventAlertType
	AggregationKey []byte
	SourceTypeName []byte
	pool           *eventPool
}

func (e *Event) release() {
	e.pool.Put(e)
}

func newEventPool() *eventPool {
	pool := &eventPool{
		pool: &sync.Pool{},
	}
	pool.pool.New = func() interface{} {
		return &Event{pool: pool}
	}
	return pool
}

type eventPool struct {
	pool *sync.Pool
}

func (m *eventPool) Get() *Event {
	return m.pool.Get().(*Event)
}

func (m *eventPool) Put(event *Event) {
	event.Tags = event.Tags[:]
	m.pool.Put(event)
}

// ServiceCheck is a service check originating from DogStatsD
// Structuraly, this is similar to metrics.ServiceCheck with []byte slices
// instead of strings. Under the hood those []byte slices are pointing to
// memory allocated in the packet they were received from.
type ServiceCheck struct {
	Name      []byte
	Hostname  []byte
	Timestamp int64
	Status    metrics.ServiceCheckStatus
	Message   []byte
	Tags      [][]byte
	ExtraTags []string
	pool      *serviceCheckPool
}

func (sc *ServiceCheck) release() {
	sc.pool.Put(sc)
}

func newServiceCheckPool() *serviceCheckPool {
	pool := &serviceCheckPool{
		pool: &sync.Pool{},
	}
	pool.pool.New = func() interface{} {
		return &ServiceCheck{pool: pool}
	}
	return pool
}

type serviceCheckPool struct {
	pool *sync.Pool
}

func (m *serviceCheckPool) Get() *ServiceCheck {
	return m.pool.Get().(*ServiceCheck)
}

func (m *serviceCheckPool) Put(serviceCheck *ServiceCheck) {
	serviceCheck.Tags = serviceCheck.Tags[:]
	m.pool.Put(serviceCheck)
}

// ParsedPacket is the parsed content of a packet
type ParsedPacket struct {
	Samples       []*MetricSample
	Events        []Event
	ServiceChecks []ServiceCheck
	packet        *Packet
}

// Release releases the parsed packet memory to it's pool
func (p *ParsedPacket) Release() {
	p.packet.release2()
	for i := 0; i < len(p.Samples); i++ {
		p.Samples[i].release()
	}
	for i := 0; i < len(p.Events); i++ {
		p.Events[i].release()
	}
	for i := 0; i < len(p.ServiceChecks); i++ {
		p.ServiceChecks[i].release()
	}
}
