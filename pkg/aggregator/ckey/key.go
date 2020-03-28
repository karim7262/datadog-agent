// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

package ckey

import (
	"sort"

	"github.com/cespare/xxhash"
)

// ContextKey is a non-cryptographic hash that allows to
// aggregate metrics from a same context together.
//
// This implementation has been designed to remove all heap
// allocations from the intake to reduce GC pressure on high
// volumes.
//
// It uses the 128bit murmur3 hash, that is already successfully
// used on other products. 128bit is probably overkill for avoiding
// collisions, but it's better to err on the safe side, as we do not
// have a collision mitigation mechanism.
// TODO(remy): comment me
type ContextKey uint64

// KeyGenerator generates key
// Not safe for concurrent usage
type KeyGenerator struct {
	buf []byte
}

// NewKeyGenerator creates a new key generator
func NewKeyGenerator() *KeyGenerator {
	return &KeyGenerator{
		buf: make([]byte, 0, 1024),
	}
}

// Generate returns the ContextKey hash for the given parameters.
// The tags array is sorted in place to avoid heap allocations.
// func (g *KeyGenerator) Generate(name, hostname string, tags []string) ContextKey {
// 	g.buf = g.buf[:0]
//
// 	// Sort the tags in place. For typical tag slices, we use
// 	// the in-place section sort to avoid heap allocations.
// 	// We default to stdlib's sort package for longer slices.
// 	if len(tags) < 20 {
// 		selectionSort(tags)
// 	} else {
// 		sort.Strings(tags)
// 	}
//
// 	g.buf = append(g.buf, name...)
// 	g.buf = append(g.buf, ',')
// 	for i := 0; i < len(tags); i++ {
// 		g.buf = append(g.buf, tags[i]...)
// 		g.buf = append(g.buf, ',')
// 	}
// 	g.buf = append(g.buf, hostname...)
//
// 	return ContextKey(murmur3.Sum64(g.buf))
// }

// Generate returns the ContextKey hash for the given parameters.
// The tags array is sorted in place to avoid heap allocations.
func (g *KeyGenerator) Generate(name, hostname string, tags []string) ContextKey {
	g.buf = g.buf[:0]

	// Sort the tags in place. For typical tag slices, we use
	// the in-place section sort to avoid heap allocations.
	// We default to stdlib's sort package for longer slices.
	if len(tags) < 20 {
		selectionSort(tags)
	} else {
		sort.Strings(tags)
	}

	g.buf = append(g.buf, name...)
	g.buf = append(g.buf, ',')
	for i := 0; i < len(tags); i++ {
		g.buf = append(g.buf, tags[i]...)
		g.buf = append(g.buf, ',')
	}
	g.buf = append(g.buf, hostname...)

	return ContextKey(xxhash.Sum64(g.buf))
}

// func (g *KeyGenerator) Generate(name, hostname string, tags []string) ContextKey {
// 	// Sort the tags in place. For typical tag slices, we use
// 	// the in-place section sort to avoid heap allocations.
// 	// We default to stdlib's sort package for longer slices.
// 	if len(tags) < 20 {
// 		selectionSort(tags)
// 	} else {
// 		sort.Strings(tags)
// 	}
//
// 	hash := fnv1a.Init64
//
// 	hash = fnv1a.AddString64(hash, name)
// 	for i := 0; i < len(tags); i++ {
// 		hash = fnv1a.AddString64(hash, tags[i])
// 	}
// 	hash = fnv1a.AddString64(hash, hostname)
//
// 	return ContextKey(hash)
// }

// Equals returns whether the two context keys are equal or not.
func Equals(a, b ContextKey) bool {
	return a == b
}

// IsZero returns true if the key is at zero value
func (k ContextKey) IsZero() bool {
	return k == 0
}
