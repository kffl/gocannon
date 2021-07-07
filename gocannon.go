package main

import (
	"fmt"
	"os"
	"sync"

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

	n := *connections

	stats, scErr := newStatsCollector(*mode, n, *preallocate, *timeout)

	if scErr != nil {
		exitWithError(scErr)
	}

	var wg sync.WaitGroup

	wg.Add(n)

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
					stats.RecordResponse(cid, code, start, end)
				}
			}
			wg.Done()
		}(c, connectionID)
	}

	wg.Wait()

	stats.CalculateStats(start, stop, *interval)

	printSummary(stats)
	stats.PrintReport()

	if *outputFile != "" {
		err = stats.SaveRawData(*outputFile)
		if err != nil {
			exitWithError(err)
		}
	}
}
