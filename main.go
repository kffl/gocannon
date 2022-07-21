package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/kffl/gocannon/lib"
)

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func main() {
	config, err := parseArgs()
	if err != nil {
		exitWithError(err)
	}

	runtime.GOMAXPROCS(*config.CPUs)

	if *config.Format == "default" {
		printHeader(config)
	}

	g, err := lib.NewGocannon(config)

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
