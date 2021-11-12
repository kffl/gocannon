package main

import (
	"fmt"
	"sync/atomic"

	"github.com/kffl/gocannon/common"
)

type plugin string

var config common.Config

// you can use global variables within the plugin to persist its state between invocations of BeforeRequest
var reqCounter int64 = 0

func (p plugin) Startup(cfg common.Config) {
	// saving the config for later
	// make sure not to mutate the config contents
	config = cfg
}

func (p plugin) BeforeRequest(cid int) (target string, method string, body common.RawRequestBody, headers common.RequestHeaders) {
	// there can be multiple invocations of BeforeRequest (with different connections ids) happening in parallel
	// therefore it's necessary to ensure thread-safe usage of plugin global variables

	// as an example, we are going to use atomic add operation
	// in order to track how many requests were sent so far
	reqNum := atomic.AddInt64(&reqCounter, 1)

	headers = *config.Headers

	// every 10th request, we want to add a special header
	if reqNum%10 == 0 {
		headers = append(*config.Headers, common.RequestHeader{"X-Special-Header", "gocannon"})
	}

	// appending connectionID to the target (i.e. http://target:123/?connection=5)
	target = fmt.Sprintf("%s?connection=%d", *config.Target, cid)

	// we leave the HTTP method unchanged (from config)
	method = *config.Method

	// and the same for body
	// the body shall not be mutated after being passed as a return value, as gocannon uses fasthttp's SetBodyRaw
	body = *config.Body

	return
}

func (p plugin) GetName() string {
	return string(p)
}

// GocannonPlugin is an exported instance of the plugin (the name "GocannonPlugin" is mandatory)
var GocannonPlugin plugin = "Sample Plugin"
