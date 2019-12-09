package telemetry

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// NewGauge creates a Gauge for telemetry purpose.
func NewGauge(subsystem, name string, tags []string, help string) Gauge {
	g := &promGauge{
		pg: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      fmt.Sprintf("_%s", name),
				Help:      help,
			},
			tags,
		),
	}
	prometheus.MustRegister(g.pg)
	return g
}

// Gauge implementation using Prometheus.
type promGauge struct {
	pg   *prometheus.GaugeVec
	once sync.Once
}

// Set stores the value for the given tags.
func (g *promGauge) Set(value float64, tags ...string) {
	g.pg.WithLabelValues(tags...).Set(value)
}
