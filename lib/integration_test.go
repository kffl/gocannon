package lib

import (
	"math"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kffl/gocannon/common"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestGocannon(t *testing.T) {

	target := "http://localhost:3000/hello"
	timeout := time.Duration(200) * time.Millisecond
	duration := time.Duration(3) * time.Second
	conns := 50
	body := []byte("")
	headers := common.RequestHeaders{}

	c, err := newHTTPClient(target, timeout, conns, true, true)

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
				code, start, end := performRequest(c, target, "GET", body, headers)
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
	maxDelta := math.Min(reqStats.GetLatencyAvg(), histStats.GetLatencyAvg()) * 0.0001
	assert.LessOrEqual(
		t,
		deltaL,
		maxDelta,
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
	duration := time.Second * 1
	connections := 50
	cpus := runtime.NumCPU()
	timeout := time.Millisecond * 200
	mode := "reqlog"
	outputFile := ""
	interval := time.Millisecond * 250
	preallocate := 1000
	method := "GET"
	body := common.RawRequestBody{}
	header := common.RequestHeaders{}
	trustAll := false
	format := "default"
	plugin := ""
	target := "http://localhost:3000/hello"

	cfg := common.Config{
		Duration:    &duration,
		Connections: &connections,
		CPUs:        &cpus,
		Timeout:     &timeout,
		Mode:        &mode,
		OutputFile:  &outputFile,
		Interval:    &interval,
		Preallocate: &preallocate,
		Method:      &method,
		Body:        &body,
		Headers:     &header,
		TrustAll:    &trustAll,
		Format:      &format,
		Plugin:      &plugin,
		Target:      &target,
	}

	g, creationErr := NewGocannon(cfg)

	assert.Nil(t, creationErr, "gocannon instance should be created without errors")

	if creationErr == nil {
		results, execErr := g.Run()
		assert.Nil(t, execErr, "the load test should be completed without errors")
		assert.Greater(
			t,
			results.GetReqPerSec(),
			100.0,
			"a throughput of at least 100 req/s should be achieved",
		)
	}
}

func TestGocanonWithPlugin(t *testing.T) {

	err := exec.Command("go", "build", "-race", "-buildmode=plugin", "-o", "../_example_plugin/plugin.so", "../_example_plugin/plugin.go").
		Run()

	assert.Nil(t, err, "the plugin should compile without an error")

	duration := time.Second * 1
	connections := 50
	cpus := runtime.NumCPU()
	timeout := time.Millisecond * 200
	mode := "hist"
	outputFile := ""
	interval := time.Millisecond * 250
	preallocate := 1000
	method := "GET"
	body := common.RawRequestBody{}
	header := common.RequestHeaders{}
	trustAll := true
	format := "json"
	plugin := "../_example_plugin/plugin.so"
	target := "http://localhost:3000/hello"

	cfg := common.Config{
		Duration:    &duration,
		Connections: &connections,
		CPUs:        &cpus,
		Timeout:     &timeout,
		Mode:        &mode,
		OutputFile:  &outputFile,
		Interval:    &interval,
		Preallocate: &preallocate,
		Method:      &method,
		Body:        &body,
		Headers:     &header,
		TrustAll:    &trustAll,
		Format:      &format,
		Plugin:      &plugin,
		Target:      &target,
	}

	g, creationErr := NewGocannon(cfg)

	assert.Nil(t, creationErr, "gocannon instance with a plugin should be created without errors")

	if creationErr == nil {
		results, execErr := g.Run()

		assert.Nil(t, execErr, "the load test should be completed without errors")

		assert.Greater(
			t,
			results.GetReqPerSec(),
			100.0,
			"a throughput of at least 100 req/s should be achieved",
		)
	}

}
