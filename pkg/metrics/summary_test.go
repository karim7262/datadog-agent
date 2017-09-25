// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

package metrics

import (
	//	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSummary(t *testing.T) {
	summary := NewSummary()

	// Empty flush
	_, err := summary.flush(10)
	assert.NotNil(t, err)

	// Add samples
	summary.addSample(&MetricSample{Value: 1}, 20)
	summary.addSample(&MetricSample{Value: 1.5}, 21)
	summary.addSample(&MetricSample{Value: 0.01}, 22)
	summary.addSample(&MetricSample{Value: 5.3}, 22)
	summary.addSample(&MetricSample{Value: 2.1}, 25)
	summary.addSample(&MetricSample{Value: 7.0}, 29)

	series, err := summary.flush(30)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(series))

	for _, serie := range series {
		assert.Equal(t, 1, len(serie.Points))
		assert.Equal(t, float64(30), serie.Points[0].Ts)
	}

	assert.Equal(t, float64(7.0), series[0].Points[0].Value)
	assert.Equal(t, ".max", series[0].NameSuffix)

	assert.Equal(t, float64(0.01), series[1].Points[0].Value)
	assert.Equal(t, ".min", series[1].NameSuffix)

	assert.Equal(t, float64(2.8183333333333334), series[2].Points[0].Value)
	assert.Equal(t, ".avg", series[2].NameSuffix)

	assert.Equal(t, float64(16.91), series[3].Points[0].Value)
	assert.Equal(t, ".sum", series[3].NameSuffix)

	assert.Equal(t, float64(6), series[4].Points[0].Value)
	assert.Equal(t, ".count", series[4].NameSuffix)

}
