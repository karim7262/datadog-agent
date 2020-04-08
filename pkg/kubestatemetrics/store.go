package kubestatemetrics

import (
	"k8s.io/kube-state-metrics/pkg/metric"
	"sync"
	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
	"github.com/DataDog/datadog-agent/pkg/aggregator"
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
	metrics map[string]map[string]*metric.Family
	// generateMetricsFunc generates metrics based on a given Kubernetes object
	// and returns them grouped by metric family.
	generateMetricsFunc func(interface{}) []*metric.Family
}

// NewMetricsStore returns a new MetricsStore
func NewMetricsStore(generateFunc func(interface{}) []*metric.Family) *MetricsStore {
	return &MetricsStore{
		generateMetricsFunc: generateFunc,
		metrics:             make(map[string]map[string]*metric.Family),
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

	families := s.generateMetricsFunc(obj)

	uid := string(o.GetUID())
	if _, ok := s.metrics[uid]; !ok {
		s.metrics[uid] = make(map[string]*metric.Family)
	}

	var errs []error
	for _, f := range families {
		s.metrics[uid][f.Name] = f
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func buildTags(uid string, metrics *metric.Metric) ([]string, error) {
	tags := []string{fmt.Sprintf("uid:%s", uid)}
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

	delete(s.metrics, string(o.GetUID()))

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

func (s *MetricsStore) Push(sender aggregator.Sender) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var errs []error
	for uid := range s.metrics {
		for _, f := range s.metrics[uid] {
			s.metrics[uid][f.Name] = f

			switch f.Type {
			case metric.Gauge:
				for _, m := range f.Metrics {
					tags, err := buildTags(uid, m)
					if err != nil {
						errs = append(errs, err)
						continue
					}
					sender.Gauge(f.Name, m.Value, "", tags)
				}
			case metric.Counter:
				for _, m := range f.Metrics {
					tags, err := buildTags(uid, m)
					if err != nil {
						errs = append(errs, err)
						continue
					}
					sender.Gauge(f.Name, m.Value, "", tags)
				}
			default:
				errs = append(errs, fmt.Errorf("metric type: %s not supported", f.Type))
			}

		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
