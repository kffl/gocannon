package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMakeTimestamp(t *testing.T) {
	timestampA := makeTimestamp()
	time.Sleep(time.Millisecond)
	timestampB := makeTimestamp()

	assert.Greater(t, timestampB, timestampA, "timestamp should be increasing")
}
