package telemetry

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// NewCounter creates a new telemetry Counter.
// TODO(remy): documentation
func NewCounter(namespace, subsystem, name string, tags []string, help string) Counter {
	c := &promCounter{
		pc: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
			},
			tags,
		),
	}
	prometheus.MustRegister(c.pc)
	return c
}

// promCounter is
// TODO(remy): doc
type promCounter struct {
	pc   *prometheus.CounterVec
	once sync.Once
}

// Add adds the given value to the counter with the given tags.
func (c *promCounter) Add(value float64, tags ...string) {
	c.pc.WithLabelValues(tags...).Add(value)
}

// Inc increments the counter with the given tags.
func (c *promCounter) Inc(tags ...string) {
	c.pc.WithLabelValues(tags...).Inc()
}
