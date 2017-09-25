// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

package metrics

import (
	"github.com/DataDog/datadog-agent/pkg/metrics/percentile"
)

// Summary keeps track of the aggregate values of samples added in one
// flush period.
type Summary struct {
	aggSketch *percentile.AggSketch
}

// NewSummary returns a newly initialized summary
func NewSummary() *Summary {
	return &Summary{aggSketch: percentile.NewAggSketch()}
}

func (s *Summary) addSample(sample *MetricSample, timestamp float64) {
	s.aggSketch.Add(sample.Value)
}

func (s *Summary) flush(timestamp float64) ([]*Serie, error) {
	if s.aggSketch.Count == 0 {
		return []*Serie{}, NoSerieError{}
	}

	// Aggregations provided by Summary
	aggregations := []string{maxAgg, minAgg, avgAgg, sumAgg, countAgg}

	series := make([]*Serie, 0, len(aggregations))

	for _, aggregate := range aggregations {
		var value float64
		var mType APIMetricType
		switch aggregate {
		case maxAgg:
			value = s.aggSketch.Max
			mType = APIGaugeType
		case minAgg:
			value = s.aggSketch.Min
			mType = APIGaugeType
		case countAgg:
			value = float64(s.aggSketch.Count)
			mType = APICountType
		case sumAgg:
			value = s.aggSketch.Sum
			mType = APIGaugeType
		case avgAgg:
			value = s.aggSketch.Avg
			mType = APIGaugeType
		default:
			continue
		}

		series = append(series, &Serie{
			Points:     []Point{{Ts: timestamp, Value: value}},
			MType:      mType,
			NameSuffix: "." + aggregate,
		})
	}

	// reset the summary
	s.aggSketch = percentile.NewAggSketch()

	return series, nil

}
