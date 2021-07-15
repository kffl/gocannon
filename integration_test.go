package main

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestGocannon(t *testing.T) {

	target := "http://localhost:3000/hello"
	timeout := time.Duration(200) * time.Millisecond
	duration := time.Duration(3) * time.Second
	conns := 50

	c, err := newHTTPClient(target, timeout, conns)

	var ops int32
	var wg sync.WaitGroup

	wg.Add(conns)

	reqStats, _ := newStatsCollector("reqlog", conns, 1000, timeout)
	histStats, _ := newStatsCollector("hist", conns, 1000, timeout)

	start := makeTimestamp()
	stop := start + duration.Nanoseconds()

	for connectionID := 0; connectionID < conns; connectionID++ {
		go func(c *fasthttp.HostClient, cid int) {
			for {
				code, start, end := performRequest(c, target)
				if end >= stop {
					break
				}

				if code != 0 {
					atomic.AddInt32(&ops, 1)
					reqStats.RecordResponse(cid, code, start, end)
					histStats.RecordResponse(cid, code, start, end)
				}
			}
			wg.Done()
		}(c, connectionID)
	}

	wg.Wait()

	reqStats.CalculateStats(start, stop, time.Duration(250)*time.Millisecond, "")
	histStats.CalculateStats(start, stop, time.Duration(250)*time.Millisecond, "")

	assert.Nil(t, err, "the http client should be created without an error")
	assert.Equal(
		t,
		int64(ops),
		reqStats.GetReqCount(),
		"request count calculated by reqlog and by atomic counter should be equal",
	)
	assert.Equal(
		t,
		int64(ops),
		histStats.GetReqCount(),
		"request count calculated by hist and by atomic counter should be equal",
	)
	assert.Equal(
		t,
		reqStats.GetReqPerSec(),
		histStats.GetReqPerSec(),
		"requests per second calculated by reqlog and by hist should be equal",
	)
	deltaL := math.Abs(reqStats.GetLatencyAvg() - histStats.GetLatencyAvg())
	assert.LessOrEqual(
		t,
		deltaL,
		float64(2),
		"average latencies calculated by reqlog and hist should be within the error margin",
	)
	assert.Greater(
		t,
		reqStats.GetReqCount(),
		int64(8000),
		"should send over 8k requests in 3 seconds",
	)
}

func TestGocannonDefaultValues(t *testing.T) {
	*duration = time.Second * 1
	*connections = 50
	*timeout = time.Millisecond * 200
	*mode = "reqlog"
	*outputFile = ""
	*interval = time.Millisecond * 250
	*preallocate = 1000
	*target = "http://localhost:3000/hello"

	assert.Nil(t, runGocannon(), "the load test should be completed without errors")
}
