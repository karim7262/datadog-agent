package telemetry

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// NewCounter creates a new telemetry Counter.
// TODO(remy): documentation
func NewCounter(namespace, subsystem, name string, tags []string, help string) Counter {
	c := &PromCounter{
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

// PromCounter is
// TODO(remy): doc
type PromCounter struct {
	pc   *prometheus.CounterVec
	once sync.Once
}

// Reset resets the counter to 0.
func (c *PromCounter) Reset() {
	c.pc.Reset()
}

// Add adds the given value to the counter with the given tags.
func (c *PromCounter) Add(value float64, tags ...string) {
	c.pc.WithLabelValues(tags...).Add(value)
}

// Inc increments the counter with the given tags.
func (c *PromCounter) Inc(tags ...string) {
	c.pc.WithLabelValues(tags...).Inc()
}
