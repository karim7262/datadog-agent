package dogstatsd

import "github.com/DataDog/datadog-agent/pkg/metrics"

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
}

// ParsedPacket is the parsed content of a packet
type ParsedPacket struct {
	Samples       []MetricSample
	Events        []Event
	ServiceChecks []ServiceCheck
	packet        *Packet
}

// Release releases the parsed packet memory to it's pool
func (p *ParsedPacket) Release() {
	p.packet.release2()
}
