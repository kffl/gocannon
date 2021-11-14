package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/kffl/gocannon/common"
	"github.com/valyala/fasthttp"
)

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func runGocannon(cfg common.Config) error {
	var gocannonPlugin common.GocannonPlugin
	var err error

	if *cfg.Plugin != "" {
		gocannonPlugin, err = loadPlugin(*cfg.Plugin)
		if err != nil {
			return err
		}
		gocannonPlugin.Startup(cfg)
	}

	c, err := newHTTPClient(*cfg.Target, *cfg.Timeout, *cfg.Connections, *cfg.TrustAll, true)

	if err != nil {
		return err
	}

	n := *cfg.Connections

	stats, scErr := newStatsCollector(*cfg.Mode, n, *cfg.Preallocate, *cfg.Timeout)

	if scErr != nil {
		return scErr
	}

	var wg sync.WaitGroup

	wg.Add(n)

	start := makeTimestamp()
	stop := start + cfg.Duration.Nanoseconds()

	fmt.Printf("gocannon goes brr...\n")

	for connectionID := 0; connectionID < n; connectionID++ {
		go func(c *fasthttp.HostClient, cid int, p common.GocannonPlugin) {
			for {
				var code int
				var start int64
				var end int64
				if p != nil {
					plugTarget, plugMethod, plugBody, plugHeaders := p.BeforeRequest(cid)
					code, start, end = performRequest(c, plugTarget, plugMethod, plugBody, plugHeaders)
				} else {
					code, start, end = performRequest(c, *cfg.Target, *cfg.Method, *cfg.Body, *cfg.Headers)
				}
				if end >= stop {
					break
				}

				stats.RecordResponse(cid, code, start, end)
			}
			wg.Done()
		}(c, connectionID, gocannonPlugin)
	}

	wg.Wait()

	err = stats.CalculateStats(start, stop, *cfg.Interval, *cfg.OutputFile)

	if err != nil {
		return err
	}

	printSummary(stats)
	stats.PrintReport()

	return nil
}

func main() {
	err := parseArgs()
	if err != nil {
		exitWithError(err)
	}

	printHeader(config)

	err = runGocannon(config)

	if err != nil {
		exitWithError(err)
	}
}
