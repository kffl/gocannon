package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/valyala/fasthttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	duration = kingpin.Flag("duration", "Load test duration").
			Short('d').
			Default("10s").
			Duration()
	connections = kingpin.Flag("connections", "Maximum number of concurrent connections").
			Short('c').
			Default("50").
			Int()
	timeout = kingpin.Flag("timeout", "HTTP client timeout").
		Short('t').
		Default("200ms").
		Duration()
	outputFile = kingpin.Flag("output", "File to save the CSV output (raw req data)").
			PlaceHolder("file.csv").
			Short('o').
			String()
	interval = kingpin.Flag("interval", "Interval for statistics calculation").
			Default("250ms").
			Short('i').
			Duration()
	target = kingpin.Arg("target", "HTTP target URL").Required().String()
)

func parseArgs() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
}

func printHeader() {
	fmt.Printf("Attacking %s with %d connections over %s\n", *target, *connections, *duration)
	fmt.Printf("gocannon goes brr...\n")
}

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
	var ops uint64
	var wg sync.WaitGroup

	n := *connections

	wg.Add(n)

	reqs := newRequests(n)

	start := makeTimestamp()
	stop := start + duration.Nanoseconds()

	for connectionID := 0; connectionID < n; connectionID++ {
		go func(c *fasthttp.HostClient, cid int) {
			for {
				code, start, end := performRequest(c, *target)
				if end >= stop {
					break
				}
				atomic.AddUint64(&ops, 1)

				if code != -1 {
					reqs[cid] = append(reqs[cid], request{code, start, end})
				}
			}
			wg.Done()
		}(c, connectionID)
	}

	wg.Wait()

	fmt.Printf("Total Req: %8d\n", ops)
	fmt.Printf("Req/s:     %11.2f\n", float64(ops)/(duration.Seconds()))

	reqsFlattened := reqs.flatten()
	reqsFlattened.sort()

	stats, detailedStats := calculateStats(reqsFlattened, start, stop, int64(*interval))

	fmt.Printf("Interval stats: (interval = %v) \n", *interval)
	stats.printHeader()

	for _, s := range detailedStats {
		s.print()
	}

	fmt.Println("----------")

	stats.print()

	if *outputFile != "" {
		err := reqsFlattened.saveCSV(start)
		if err != nil {
			exitWithError(err)
		}
		fmt.Printf("Raw request data saved to CSV file: %s\n", *outputFile)
	}
}
