package main

import (
	"fmt"

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
	preallocate = kingpin.Flag("preallocate", "Number of requests in log to preallocate memory for per connection").
			Default("1000").
			Int()
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
