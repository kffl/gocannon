package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/kffl/gocannon/stats"
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

	reqs := stats.NewRequests(n, *preallocate)

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
					reqs.RecordResponse(cid, code, start, end)
				}
			}
			wg.Done()
		}(c, connectionID)
	}

	wg.Wait()

	stats, err := reqs.CalculateStats(start, stop, *interval, *outputFile)
	stats.Print()

	if err != nil {
		exitWithError(err)
	}
}
