package telemetry

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// NewCounter creates a Counter for telemetry purpose.
func NewCounter(subsystem, name string, tags []string, help string) Counter {
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

// Counter implementation using Prometheus.
type promCounter struct {
	pc   *prometheus.CounterVec
	once sync.Once
}

// Add adds the given value to the counter for the given tags.
func (c *promCounter) Add(value float64, tags ...string) {
	c.pc.WithLabelValues(tags...).Add(value)
}

// Inc increments the counter for the given tags.
func (c *promCounter) Inc(tags ...string) {
	c.pc.WithLabelValues(tags...).Inc()
}
