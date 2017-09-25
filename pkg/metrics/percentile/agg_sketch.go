// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

package percentile

import (
	"math"
)

// AggSketch tracks the aggregate values of samples added over one flush period.
type AggSketch struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Count int64   `json:"cnt"`
	Sum   float64 `json:"sum"`
	Avg   float64 `json:"avg"`
}

// NewAggSketch returns a newly initialized AggSketch
func NewAggSketch() *AggSketch {
	return &AggSketch{Min: math.Inf(1), Max: math.Inf(-1)}
}

// Add adds a value to the summary
func (s *AggSketch) Add(v float64) {
	s.Count++
	s.Sum += v
	s.Avg += (v - s.Avg) / float64(s.Count)
	if v < s.Min {
		s.Min = v
	}
	if v > s.Max {
		s.Max = v
	}
}

// Merge combines another summary with this.
func (s *AggSketch) Merge(o *AggSketch) {
	if o.Count == 0 {
		return
	}
	if s.Count == 0 {
		s.Min = o.Min
		s.Max = o.Max
		s.Count = o.Count
		s.Sum = o.Sum
		s.Avg = o.Avg
		return
	}

	s.Count += o.Count
	s.Sum += o.Sum
	s.Avg = s.Avg + (o.Avg-s.Avg)*float64(o.Count)/float64(s.Count)
	if o.Min < s.Min {
		s.Min = o.Min
	}
	if o.Max > s.Max {
		s.Max = o.Max
	}
}
