package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequests(t *testing.T) {
	reqs := newRequests(5)

	assert.Len(t, reqs, 5, "should have the specified length")
	assert.Len(t, reqs[0], 0, "should have empty slices as elements")
}

func TestRequestSort(t *testing.T) {
	reqs := flattenedRequests{
		{200, 123, 223},
		{200, 123, 223},
		{200, 234, 235},
		{200, 234, 235},
		{200, 534, 535},
	}

	reqs.sort()

	var endTimes []int64

	for _, r := range reqs {
		endTimes = append(endTimes, r.end)
	}

	assert.IsNonDecreasing(t, endTimes, "requests should be sorted by end times in ascending order")
}
