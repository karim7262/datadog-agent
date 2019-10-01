package telemetry

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// NewGauge creates a gauge telemetry Gauge.
// TODO(remy): doc
func NewGauge(namespace, subsystem, name string, tags []string, help string) Gauge {
	g := &PromGauge{
		pg: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      name,
				Help:      help,
			},
			tags,
		),
	}
	prometheus.MustRegister(g.pg)
	return g
}

// PromGauge is TODO(remy):
type PromGauge struct {
	pg   *prometheus.GaugeVec
	once sync.Once
}

// Reset resets the counter to 0.
func (g *PromGauge) Reset() {
	g.pg.Reset()
}

// Set sets the gauge with the given value.
func (g *PromGauge) Set(value float64, tags ...string) {
	g.pg.WithLabelValues(tags...).Add(value)
}
