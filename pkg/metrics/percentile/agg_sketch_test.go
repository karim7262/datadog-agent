// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

package percentile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var generators = []Generator{
	NewConstant(42),
	NewUniform(),
	NewNormal(35, 1),
	NewExponential(2),
}

func AssertAggSketchCorrect(t *testing.T,
	dataset *Dataset, aggSketch *AggSketch) {

	eps := float64(1e-6)
	assert.Equal(t, dataset.Min(), aggSketch.Min)
	assert.Equal(t, dataset.Max(), aggSketch.Max)
	assert.Equal(t, dataset.Count, aggSketch.Count)
	assert.InEpsilon(t, dataset.Sum(), aggSketch.Sum, eps)
	assert.InEpsilon(t, dataset.Avg(), aggSketch.Avg, eps)

}

func TestAggSketchAdd(t *testing.T) {
	n := 1000
	for _, gen := range generators {
		a := NewAggSketch()
		d := NewDataset()

		for i := 0; i < n; i++ {
			value := gen.Generate()
			a.Add(value)
			d.Add(value)
		}
		AssertAggSketchCorrect(t, d, a)
	}
}

func TestAggSketchMergeEmpty(t *testing.T) {
	n := 1000
	d := NewDataset()
	a1 := NewAggSketch()
	a2 := NewAggSketch()
	a3 := NewAggSketch()

	generator := NewExponential(2)
	for i := 0; i < n; i++ {
		value := generator.Generate()
		a2.Add(value)
		d.Add(value)
	}

	a1.Merge(a2)
	AssertAggSketchCorrect(t, d, a1)

	a2.Merge(a3)
	AssertAggSketchCorrect(t, d, a2)
}

func TestAggSketchMergeMixed(t *testing.T) {
	n := 1000
	d := NewDataset()
	a1 := NewAggSketch()
	generator1 := NewNormal(100, 1)
	for i := 0; i < n; i++ {
		value := generator1.Generate()
		a1.Add(value)
		d.Add(value)
	}

	a2 := NewAggSketch()
	generator2 := NewExponential(5)
	for i := 0; i < n; i += 10 {
		value := generator2.Generate()
		a2.Add(value)
		d.Add(value)
	}
	a1.Merge(a2)

	a3 := NewAggSketch()
	generator3 := NewExponential(0.1)
	for i := 0; i < n; i += 2 {
		value := generator3.Generate()
		a3.Add(value)
		d.Add(value)
	}
	a1.Merge(a3)

	AssertAggSketchCorrect(t, d, a1)
}
