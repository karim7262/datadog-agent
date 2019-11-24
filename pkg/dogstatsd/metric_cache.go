// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package dogstatsd

import (
	"bytes"
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/metrics"
	"github.com/hashicorp/golang-lru"
)

type metricCache struct {
	cache *lru.Cache
}

func newMetricCache(size int) (*metricCache, error) {
	cache, err := lru.New(size)
	if err != nil {
		return &metricCache{}, err
	}
	return &metricCache{cache: cache}, nil
}

func (m *metricCache) get(metricName string) (*metrics.MetricSample) {
	if result, ok := m.cache.Get(metricName); ok {
		return result.(*metrics.MetricSample)
	}
	return nil
}

func (m *metricCache) add(metricName string, sample *metrics.MetricSample) {
	m.cache.Add(metricName, sample)
}

func buildMetricCacheKey(message []byte) (string, error) {
	valStartIndex := bytes.Index(message, valueSeparator)
	if valStartIndex == -1 {
		return "", fmt.Errorf("separator `%q` not found in %q", valueSeparator, message)
	}
	valStartIndex++  // skip separator

	valEndIndex := bytes.Index(message[valStartIndex:], fieldSeparator)
	if valEndIndex == -1 {
		return "", fmt.Errorf("separator `%q` not found in %q", fieldSeparator, message)
	}
	valEndIndex = valEndIndex + valStartIndex
	return string(message[:valStartIndex]) + string(message[valEndIndex:]), nil
}
