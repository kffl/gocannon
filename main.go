package main

import (
	"fmt"
	"runtime"
)

func main() {
	config, err := parseArgs()
	if err != nil {
		exitWithError(err)
	}

	runtime.GOMAXPROCS(*config.CPUs)

	if *config.Format == "default" {
		printHeader(config)
	}

	g, err := NewGocannon(config)

	if err != nil {
		exitWithError(err)
	}

	if *config.Format == "default" {
		fmt.Printf("gocannon goes brr...\n")
	}

	results, err := g.Run()

	if *config.Format == "default" {
		printSummary(results)
	}
	results.PrintReport(*config.Format)

	if err != nil {
		exitWithError(err)
	}
}
