package main

import (
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
	mode = kingpin.Flag("mode", "Statistics collection mode: reqlog (logs each request) or hist (stores histogram of completed requests latencies)").
		Default("reqlog").
		Short('m').
		String()
	outputFile = kingpin.Flag("output", "File to save the request log in CSV format (reqlog mode) or a text file with raw histogram data (hist mode)").
			PlaceHolder("file.csv").
			Short('o').
			String()
	interval = kingpin.Flag("interval", "Interval for statistics calculation (reqlog mode)").
			Default("250ms").
			Short('i').
			Duration()
	preallocate = kingpin.Flag("preallocate", "Number of requests in req log to preallocate memory for per connection (reqlog mode)").
			Default("1000").
			Int()
	target = kingpin.Arg("target", "HTTP target URL").Required().String()
)

func parseArgs() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
}
