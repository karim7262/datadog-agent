package ebpf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnWithHigherStats(t *testing.T) {
	for i, tc := range []struct {
		c1, c2     ConnectionStats
		c1Expected bool
	}{
		{
			c1:         ConnectionStats{MonotonicSentBytes: 2},
			c2:         ConnectionStats{MonotonicSentBytes: 1},
			c1Expected: true,
		},
		{
			c1:         ConnectionStats{MonotonicSentBytes: 2, MonotonicRecvBytes: 50},
			c2:         ConnectionStats{MonotonicSentBytes: 1, MonotonicRecvBytes: 50},
			c1Expected: true,
		},
		{
			c1:         ConnectionStats{MonotonicSentBytes: 2, MonotonicRecvBytes: 50},
			c2:         ConnectionStats{MonotonicSentBytes: 2, MonotonicRecvBytes: 50, MonotonicRetransmits: 1},
			c1Expected: false,
		},
	} {
		res := connWithHigherStats(tc.c1, tc.c2)
		if tc.c1Expected {
			assert.Equalf(t, tc.c1, res, "test %d", i)
		} else {
			assert.Equalf(t, tc.c2, res, "test %d", i)
		}
	}
}

func TestCompareConnStats(t *testing.T) {
	for i, tc := range []struct {
		c1, c2   ConnectionStats
		expected bool
	}{
		{
			c1:       ConnectionStats{MonotonicSentBytes: 2},
			c2:       ConnectionStats{MonotonicSentBytes: 1},
			expected: false,
		},
		{
			c1:       ConnectionStats{MonotonicSentBytes: 2, LastUpdateEpoch: 1},
			c2:       ConnectionStats{MonotonicSentBytes: 2, LastUpdateEpoch: 2},
			expected: false,
		},
		{
			c1:       ConnectionStats{MonotonicSentBytes: 2, LastUpdateEpoch: 2},
			c2:       ConnectionStats{MonotonicSentBytes: 2, LastUpdateEpoch: 2},
			expected: true,
		},
		{
			c1:       ConnectionStats{MonotonicSentBytes: 2, MonotonicRecvBytes: 50},
			c2:       ConnectionStats{MonotonicSentBytes: 2, MonotonicRecvBytes: 50},
			expected: true,
		},
		{
			c1:       ConnectionStats{MonotonicSentBytes: 2, MonotonicRecvBytes: 49},
			c2:       ConnectionStats{MonotonicSentBytes: 2, MonotonicRecvBytes: 50},
			expected: false,
		},
		{
			c1:       ConnectionStats{MonotonicSentBytes: 2, MonotonicRecvBytes: 50},
			c2:       ConnectionStats{MonotonicSentBytes: 2, MonotonicRecvBytes: 50, MonotonicRetransmits: 1},
			expected: false,
		},
	} {
		assert.Equalf(t, tc.expected, compareConnsStats(tc.c1, tc.c2), "test %d", i)
	}
}
