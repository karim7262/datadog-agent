package metrics

import (
	"sync"
)

// MetricSamplePool is a pool of metrics sample
type MetricSamplePool struct {
	pool      *sync.Pool
	batchSize int
}

// NewMetricSamplePool creates a new MetricSamplePool
func NewMetricSamplePool(batchSize int) *MetricSamplePool {
	return &MetricSamplePool{
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]MetricSample, batchSize)
			},
		},
		batchSize: batchSize,
	}
}

// GetBatch gets a batch of metric samples from the pool
func (m *MetricSamplePool) GetBatch() []MetricSample {
	if m == nil {
		return nil
	}
	return make([]MetricSample, m.batchSize)
}

// PutBatch puts a batch back into the pool
func (m *MetricSamplePool) PutBatch(batch []MetricSample) {
	if m == nil {
		return
	}
}
