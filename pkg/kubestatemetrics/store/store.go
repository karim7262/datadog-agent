package store

import (
	"k8s.io/kube-state-metrics/pkg/metric"
	"sync"
	"k8s.io/apimachinery/pkg/types"

	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)


// MetricsStore implements the k8s.io/client-go/tools/cache.Store
// interface. Instead of storing entire Kubernetes objects, it stores metrics
// generated based on those objects.
type MetricsStore struct {
	// Protects metrics
	mutex sync.RWMutex
	// metrics is a map indexed by Kubernetes object id, containing a slice of
	// metric families, containing a slice of metrics. We need to keep metrics
	// grouped by metric families in order to zip families with their help text in
	// MetricsStore.WriteAll().
	metrics map[types.UID][]DDMetricsFam
	// generateMetricsFunc generates metrics based on a given Kubernetes object
	// and returns them grouped by metric family.
	generateMetricsFunc func(interface{}) []metric.FamilyInterface
}

// NewMetricsStore returns a new MetricsStore
func NewMetricsStore(generateFunc func(interface{}) []metric.FamilyInterface) *MetricsStore {
	return &MetricsStore{
		generateMetricsFunc: generateFunc,
		metrics:             map[types.UID][]DDMetricsFam{},
	}
}

type DDMetric struct {
	Labels []string
	Val float64
}

type DDMetricsFam struct {
	Name string
	listMetrics []DDMetric
}

func (d *DDMetricsFam) extract(f metric.Family) {
	d.Name = f.Name
	for _, m := range f.Metrics {
		var err error
		s := DDMetric{}
		s.Val = m.Value
		s.Labels, err = buildTags(m)
		if err != nil {
			// TODO test how verbose that could be.
			log.Errorf("Could not retrieve the labels for %s: %v", f.Name, err)
			continue
		}
		d.listMetrics = append(d.listMetrics, s)
	}
}

// Implementing k8s.io/client-go/tools/cache.Store interface

// Add inserts adds to the MetricsStore by calling the metrics generator functions and
// adding the generated metrics to the metrics map that underlies the MetricStore.
func (s *MetricsStore) Add(obj interface{}) error {
	o, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	metricsForUID := s.generateMetricsFunc(obj)
	convertedMetricsForUID := make([]DDMetricsFam, len(metricsForUID))
	for i, f := range metricsForUID {
		metricConvertedList := DDMetricsFam{}
		f.Inspect(metricConvertedList.extract)
		convertedMetricsForUID[i] = metricConvertedList
	}

	s.metrics[o.GetUID()] = convertedMetricsForUID

	return nil
}

func buildTags(metrics *metric.Metric) ([]string, error) {
	var tags []string
	if len(metrics.LabelKeys) != len(metrics.LabelValues) {
		return nil, fmt.Errorf("LabelKeys and LabelValues not same size")
	}
	for i := range metrics.LabelKeys {
		tags = append(tags, fmt.Sprintf("%s:%s", metrics.LabelKeys[i], metrics.LabelValues[i]))
	}
	return tags, nil
}

// Update updates the existing entry in the MetricsStore.
func (s *MetricsStore) Update(obj interface{}) error {
	// TODO: For now, just call Add, in the future one could check if the resource version changed?
	return s.Add(obj)
}

// Delete deletes an existing entry in the MetricsStore.
func (s *MetricsStore) Delete(obj interface{}) error {
	o, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.metrics, o.GetUID())

	return nil
}

// List implements the List method of the store interface.
func (s *MetricsStore) List() []interface{} {
	return nil
}

// ListKeys implements the ListKeys method of the store interface.
func (s *MetricsStore) ListKeys() []string {
	return nil
}

// Get implements the Get method of the store interface.
func (s *MetricsStore) Get(obj interface{}) (item interface{}, exists bool, err error) {
	return nil, false, nil
}

// GetByKey implements the GetByKey method of the store interface.
func (s *MetricsStore) GetByKey(key string) (item interface{}, exists bool, err error) {
	return nil, false, nil
}

// Replace will delete the contents of the store, using instead the
// given list.
func (s *MetricsStore) Replace(list []interface{}, _ string) error {
	for _, o := range list {
		err := s.Add(o)
		if err != nil {
			return err
		}
	}

	return nil
}

// Resync implements the Resync method of the store interface.
func (s *MetricsStore) Resync() error {
	return nil
}

func (s *MetricsStore) Push() map[string][]DDMetric {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	res := make(map[string][]DDMetric)
	// to return f

	for u := range s.metrics {
		for _, m := range s.metrics[u] {
			// Metrics of a specific object UID
			for _, me := range m.listMetrics {
				uidAdd := append(me.Labels, fmt.Sprintf("uid:%s", u))
				res[m.Name] = append(res[m.Name], DDMetric{
					Val: me.Val,
					Labels: uidAdd,
				})
			}
		}
	}
	return res

}
