package dogstatsd

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/DataDog/datadog-agent/pkg/metrics"
)

var (
	gaugeSymbol        = []byte("g")
	countSymbol        = []byte("c")
	histogramSymbol    = []byte("h")
	distributionSymbol = []byte("d")
	setSymbol          = []byte("s")
	timingSymbol       = []byte("ms")

	tagsFieldPrefix       = []byte("#")
	sampleRateFieldPrefix = []byte("@")
)

// sanity checks a given message against the metric sample format
func hasMetricSampleFormat(message []byte) bool {
	if message == nil {
		return false
	}
	separatorCount := bytes.Count(message, fieldSeparator)
	if separatorCount < 1 || separatorCount > 3 {
		return false
	}
	return true
}

func parseMetricSampleNameAndRawValue(rawNameAndValue []byte) ([]byte, []byte, error) {
	sepIndex := bytes.Index(rawNameAndValue, colonSeparator)
	if sepIndex == -1 {
		return nil, nil, fmt.Errorf("invalid name and value: %q", rawNameAndValue)
	}
	rawName := rawNameAndValue[:sepIndex]
	rawValue := rawNameAndValue[sepIndex+1:]
	if len(rawName) == 0 || len(rawValue) == 0 {
		return nil, nil, fmt.Errorf("invalid name and value: %q", rawNameAndValue)
	}
	return rawName, rawValue, nil
}

func parseMetricSampleMetricType(rawMetricType []byte) (metrics.MetricType, error) {
	switch {
	case bytes.Equal(rawMetricType, gaugeSymbol):
		return metrics.GaugeType, nil
	case bytes.Equal(rawMetricType, countSymbol):
		return metrics.CounterType, nil
	case bytes.Equal(rawMetricType, histogramSymbol):
		return metrics.HistogramType, nil
	case bytes.Equal(rawMetricType, distributionSymbol):
		return metrics.DistributionType, nil
	case bytes.Equal(rawMetricType, setSymbol):
		return metrics.SetType, nil
	case bytes.Equal(rawMetricType, timingSymbol):
		return metrics.HistogramType, nil
	}
	return 0, fmt.Errorf("invalid metric type: %q", rawMetricType)
}

func parseMetricSampleSampleRate(rawSampleRate []byte) (float64, error) {
	return strconv.ParseFloat(string(rawSampleRate), 64)
}

func parseMetricSample(pool *metricSamplePool, message []byte) (*MetricSample, error) {
	// fast path to eliminate most of the gibberish
	// especially important here since all the unidentified garbage gets
	// identified as metrics
	if !hasMetricSampleFormat(message) {
		return nil, fmt.Errorf("invalid dogstatsd message format: %q", message)
	}

	rawNameAndValue, message := nextField(message)
	name, rawValue, err := parseMetricSampleNameAndRawValue(rawNameAndValue)
	if err != nil {
		return nil, err
	}

	rawMetricType, message := nextField(message)
	metricType, err := parseMetricSampleMetricType(rawMetricType)
	if err != nil {
		return nil, err
	}

	var setValue []byte
	var value float64
	if metricType == metrics.SetType {
		setValue = rawValue
	} else {
		value, err = strconv.ParseFloat(string(rawValue), 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse dogstatsd metric value: %v", err)
		}
	}

	sample := pool.Get()
	sample.Name = name
	sample.Value = value
	sample.SetValue = setValue
	sample.MetricType = metricType
	sample.SampleRate = 1.0

	var optionalField []byte
	for message != nil {
		optionalField, message = nextField(message)
		if bytes.HasPrefix(optionalField, tagsFieldPrefix) {
			sample.Tags = appendTags(sample.Tags, optionalField[1:])
		} else if bytes.HasPrefix(optionalField, sampleRateFieldPrefix) {
			sample.SampleRate, err = parseMetricSampleSampleRate(optionalField[1:])
			if err != nil {
				return nil, fmt.Errorf("could not parse dogstatsd sample rate %q", optionalField)
			}
		}
	}

	return sample, nil
}
