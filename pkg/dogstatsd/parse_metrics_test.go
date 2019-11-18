package dogstatsd

import (
	"testing"

	"github.com/DataDog/datadog-agent/pkg/metrics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const epsilon = 0.00001

func TestParseGauge(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("daemon:666|g"))

	assert.NoError(t, err)

	assert.Equal(t, []byte("daemon"), sample.Name)
	assert.InEpsilon(t, 666.0, sample.Value, epsilon)
	assert.Equal(t, metrics.GaugeType, sample.MetricType)
	assert.Equal(t, 0, len(sample.Tags))
	assert.InEpsilon(t, 1.0, sample.SampleRate, epsilon)
}

func TestParseCounter(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("daemon:21|c"))

	assert.NoError(t, err)

	assert.Equal(t, []byte("daemon"), sample.Name)
	assert.Equal(t, 21.0, sample.Value)
	assert.Equal(t, metrics.CounterType, sample.MetricType)
	assert.Equal(t, 0, len(sample.Tags))
	assert.InEpsilon(t, 1.0, sample.SampleRate, epsilon)
}

func TestParseCounterWithTags(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("custom_counter:1|c|#protocol:http,bench"))

	assert.NoError(t, err)

	assert.Equal(t, []byte("custom_counter"), sample.Name)
	assert.Equal(t, 1.0, sample.Value)
	assert.Equal(t, metrics.CounterType, sample.MetricType)
	assert.Equal(t, 2, len(sample.Tags))
	assert.Equal(t, []byte("protocol:http"), sample.Tags[0])
	assert.Equal(t, []byte("bench"), sample.Tags[1])
	assert.InEpsilon(t, 1.0, sample.SampleRate, epsilon)
}

func TestParseHistogram(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("daemon:21|h"))

	assert.NoError(t, err)

	assert.Equal(t, []byte("daemon"), sample.Name)
	assert.Equal(t, 21.0, sample.Value)
	assert.Equal(t, metrics.HistogramType, sample.MetricType)
	assert.Equal(t, 0, len(sample.Tags))
	assert.InEpsilon(t, 1.0, sample.SampleRate, epsilon)
}

func TestParseTimer(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("daemon:21|ms"))

	assert.NoError(t, err)

	assert.Equal(t, []byte("daemon"), sample.Name)
	assert.Equal(t, 21.0, sample.Value)
	assert.Equal(t, metrics.HistogramType, sample.MetricType)
	assert.Equal(t, 0, len(sample.Tags))
	assert.InEpsilon(t, 1.0, sample.SampleRate, epsilon)
}

func TestParseSet(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("daemon:abc|s"))

	assert.NoError(t, err)

	assert.Equal(t, []byte("daemon"), sample.Name)
	assert.Equal(t, []byte("abc"), sample.SetValue)
	assert.Equal(t, metrics.SetType, sample.MetricType)
	assert.Equal(t, 0, len(sample.Tags))
	assert.InEpsilon(t, 1.0, sample.SampleRate, epsilon)
}

func TestSampleistribution(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("daemon:3.5|d"))

	assert.NoError(t, err)

	assert.Equal(t, []byte("daemon"), sample.Name)
	assert.Equal(t, 3.5, sample.Value)
	assert.Equal(t, metrics.DistributionType, sample.MetricType)
	assert.Equal(t, 0, len(sample.Tags))
}

func TestParseSetUnicode(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("daemon:♬†øU†øU¥ºuT0♪|s"))

	assert.NoError(t, err)

	assert.Equal(t, []byte("daemon"), sample.Name)
	assert.Equal(t, []byte("♬†øU†øU¥ºuT0♪"), sample.SetValue)
	assert.Equal(t, metrics.SetType, sample.MetricType)
	assert.Equal(t, 0, len(sample.Tags))
	assert.InEpsilon(t, 1.0, sample.SampleRate, epsilon)
}

func TestParseGaugeWithTags(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("daemon:666|g|#sometag1:somevalue1,sometag2:somevalue2"))

	assert.NoError(t, err)

	assert.Equal(t, []byte("daemon"), sample.Name)
	assert.InEpsilon(t, 666.0, sample.Value, epsilon)
	assert.Equal(t, metrics.GaugeType, sample.MetricType)
	require.Equal(t, 2, len(sample.Tags))
	assert.Equal(t, []byte("sometag1:somevalue1"), sample.Tags[0])
	assert.Equal(t, []byte("sometag2:somevalue2"), sample.Tags[1])
	assert.InEpsilon(t, 1.0, sample.SampleRate, epsilon)
}

func TestParseGaugeWithNoTags(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("daemon:666|g"))
	assert.NoError(t, err)

	assert.Equal(t, []byte("daemon"), sample.Name)
	assert.InEpsilon(t, 666.0, sample.Value, epsilon)
	assert.Equal(t, metrics.GaugeType, sample.MetricType)
	assert.Empty(t, sample.Tags)
	assert.InEpsilon(t, 1.0, sample.SampleRate, epsilon)
}

func TestParseGaugeWithSampleRate(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("daemon:666|g|@0.21"))

	assert.NoError(t, err)

	assert.Equal(t, []byte("daemon"), sample.Name)
	assert.InEpsilon(t, 666.0, sample.Value, epsilon)
	assert.Equal(t, metrics.GaugeType, sample.MetricType)
	assert.Equal(t, 0, len(sample.Tags))
	assert.InEpsilon(t, 0.21, sample.SampleRate, epsilon)
}

func TestParseGaugeWithPoundOnly(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("daemon:666|g|#"))

	assert.NoError(t, err)

	assert.Equal(t, []byte("daemon"), sample.Name)
	assert.InEpsilon(t, 666.0, sample.Value, epsilon)
	assert.Equal(t, metrics.GaugeType, sample.MetricType)
	assert.Equal(t, 0, len(sample.Tags))
	assert.InEpsilon(t, 1.0, sample.SampleRate, epsilon)
}

func TestParseGaugeWithUnicode(t *testing.T) {
	sample, err := parseMetricSample(newSamplePool(), []byte("♬†øU†øU¥ºuT0♪:666|g|#intitulé:T0µ"))

	assert.NoError(t, err)

	assert.Equal(t, []byte("♬†øU†øU¥ºuT0♪"), sample.Name)
	assert.InEpsilon(t, 666.0, sample.Value, epsilon)
	assert.Equal(t, metrics.GaugeType, sample.MetricType)
	require.Equal(t, 1, len(sample.Tags))
	assert.Equal(t, []byte("intitulé:T0µ"), sample.Tags[0])
	assert.InEpsilon(t, 1.0, sample.SampleRate, epsilon)
}

func TestParseMetricError(t *testing.T) {
	// not enough information
	_, err := parseMetricSample(newSamplePool(), []byte("daemon:666"))
	assert.Error(t, err)

	_, err = parseMetricSample(newSamplePool(), []byte("daemon:666|"))
	assert.Error(t, err)

	_, err = parseMetricSample(newSamplePool(), []byte("daemon:|g"))
	assert.Error(t, err)

	_, err = parseMetricSample(newSamplePool(), []byte(":666|g"))
	assert.Error(t, err)

	_, err = parseMetricSample(newSamplePool(), []byte("abc666|g"))
	assert.Error(t, err)

	// too many value
	_, err = parseMetricSample(newSamplePool(), []byte("daemon:666:777|g"))
	assert.Error(t, err)

	// unknown metadata prefix
	_, err = parseMetricSample(newSamplePool(), []byte("daemon:666|g|m:test"))
	assert.NoError(t, err)

	// invalid value
	_, err = parseMetricSample(newSamplePool(), []byte("daemon:abc|g"))
	assert.Error(t, err)

	// invalid metric type
	_, err = parseMetricSample(newSamplePool(), []byte("daemon:666|unknown"))
	assert.Error(t, err)

	// invalid sample rate
	_, err = parseMetricSample(newSamplePool(), []byte("daemon:666|g|@abc"))
	assert.Error(t, err)
}
