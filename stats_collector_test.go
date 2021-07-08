package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewStatsCollector(t *testing.T) {
	timeout := time.Millisecond * 200
	_, reqLogErr := newStatsCollector("reqlog", 50, 1000, timeout)
	_, histErr := newStatsCollector("hist", 50, 1000, timeout)
	_, wrongErr := newStatsCollector("sthelse", 50, 1000, timeout)

	assert.Nil(t, reqLogErr, "request log stats collector shall be created w/o error")
	assert.Nil(t, histErr, "histogram stats collector shall be created w/o error")
	assert.NotNil(t, wrongErr, "an error shall be returned when wrong mode is stats collector mode is provided")
}
