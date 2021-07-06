package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/kffl/gocannon/reqlog"
	"github.com/valyala/fasthttp"
)

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func main() {
	parseArgs()
	printHeader()

	c, err := newHTTPClient(*target, *timeout, *connections)
	if err != nil {
		exitWithError(err)
	}
	var wg sync.WaitGroup

	n := *connections

	wg.Add(n)

	reqLog := reqlog.NewRequests(n, *preallocate)

	start := makeTimestamp()
	stop := start + duration.Nanoseconds()

	for connectionID := 0; connectionID < n; connectionID++ {
		go func(c *fasthttp.HostClient, cid int) {
			for {
				code, start, end := performRequest(c, *target)
				if end >= stop {
					break
				}

				if code != -1 {
					reqLog.RecordResponse(cid, code, start, end)
				}
			}
			wg.Done()
		}(c, connectionID)
	}

	wg.Wait()

	stats, err := reqLog.CalculateStats(start, stop, *interval, *outputFile)
	stats.Print()

	if err != nil {
		exitWithError(err)
	}
}
