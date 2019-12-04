// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package util

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortUniqInPlace(t *testing.T) {
	elements := []string{"tag2:tagval", "tag1:tagval", "tag2:tagval"}
	elements = SortUniqInPlace(elements)

	assert.ElementsMatch(t, elements, []string{"tag1:tagval", "tag2:tagval"})
}

func TestSortUniq(t *testing.T) {
	elements := []string{"tag2:tagval", "tag1:tagval", "tag2:tagval"}
	sortUniqElements := make([]string, len(elements))
	sortUniqElements = SortUniq(sortUniqElements, elements)

	assert.ElementsMatch(t, sortUniqElements, []string{"tag1:tagval", "tag2:tagval"})
}

func TestSortUniqOriginalSliceUntouched(t *testing.T) {
	elements := []string{"tag2:tagval", "tag1:tagval", "tag2:tagval"}
	originalElements := make([]string, len(elements))
	copy(originalElements, elements)
	sortUniqElements := make([]string, len(elements))
	sortUniqElements = SortUniq(sortUniqElements, elements)

	assert.ElementsMatch(t, sortUniqElements, []string{"tag1:tagval", "tag2:tagval"})
	assert.Equal(t, originalElements, elements)
}

func benchmarkDeduplicateTags(b *testing.B, numberOfTags int, inPlace bool) {
	tags := make([]string, 0, numberOfTags+1)
	for i := 0; i < numberOfTags; i++ {
		tags = append(tags, fmt.Sprintf("aveeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeerylong:tag%d", i))
	}
	// this is the worst case for the insertion sort we are using
	sort.Sort(sort.Reverse(sort.StringSlice(tags)))

	tempTags := make([]string, len(tags))
	copy(tempTags, tags)
	sortedTags := make([]string, len(tempTags))
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if inPlace {
			SortUniqInPlace(tempTags)
		} else {
			SortUniq(sortedTags, tempTags)
		}
	}
}
func BenchmarkDeduplicateTags(b *testing.B) {
	for i := 1; i <= 128; i *= 2 {
		b.Run(fmt.Sprintf("deduplicate-%d-tags-clone", i), func(b *testing.B) {
			benchmarkDeduplicateTags(b, i, false)
		})
	}
	for i := 1; i <= 128; i *= 2 {
		b.Run(fmt.Sprintf("deduplicate-%d-tags-in-place", i), func(b *testing.B) {
			benchmarkDeduplicateTags(b, i, true)
		})
	}
}
