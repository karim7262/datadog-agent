package telemetry

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// NewGauge creates a gauge telemetry Gauge.
// TODO(remy): doc
func NewGauge(subsystem, name string, tags []string, help string) Gauge {
	g := &promGauge{
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

// promGauge is TODO(remy):
type promGauge struct {
	pg   *prometheus.GaugeVec
	once sync.Once
}

// Set sets the gauge with the given value.
func (g *promGauge) Set(value float64, tags ...string) {
	g.pg.WithLabelValues(tags...).Set(value)
}
