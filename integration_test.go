package main

import (
	"github.com/kffl/gocannon/stats"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGocannon(t *testing.T) {

	target := "http://localhost:3000/hello"
	timeout := time.Duration(200) * time.Millisecond
	duration := time.Duration(1) * time.Second
	conns := 50

	c, err := newHTTPClient(target, timeout, conns)

	var ops int32
	var wg sync.WaitGroup

	wg.Add(conns)

	reqs := stats.NewRequests(conns)

	start := makeTimestamp()
	stop := start + duration.Nanoseconds()

	for connectionID := 0; connectionID < conns; connectionID++ {
		go func(c *fasthttp.HostClient, cid int) {
			for {
				code, start, end := performRequest(c, target)
				if end >= stop {
					break
				}
				atomic.AddInt32(&ops, 1)

				if code != -1 {
					reqs.RecordResponse(cid, code, start, end)
				}
			}
			wg.Done()
		}(c, connectionID)
	}

	wg.Wait()

	stats, saveErr := reqs.CalculateStats(start, stop, time.Duration(250)*time.Millisecond, "")

	assert.Nil(t, err, "the http client should be created without an error")
	assert.Nil(t, saveErr, "the save error should be nil")
	assert.Equal(
		t,
		int(ops),
		stats.GetReqCount(),
		"request count calculated by stats and by atomic counter should be equal",
	)
	assert.Greater(t, stats.GetReqCount(), 1000, "should send over 1k requests in one second")
}
