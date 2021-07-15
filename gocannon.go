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

func runGocannon() error {

	c, err := newHTTPClient(*target, *timeout, *connections)

	if err != nil {
		return err
	}

	n := *connections

	stats, scErr := newStatsCollector(*mode, n, *preallocate, *timeout)

	if scErr != nil {
		return scErr
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

				stats.RecordResponse(cid, code, start, end)
			}
			wg.Done()
		}(c, connectionID)
	}

	wg.Wait()

	err = stats.CalculateStats(start, stop, *interval, *outputFile)

	if err != nil {
		return err
	}

	printSummary(stats)
	stats.PrintReport()

	return nil
}

func main() {
	parseArgs()
	printHeader()

	err := runGocannon()

	if err != nil {
		exitWithError(err)
	}
}
