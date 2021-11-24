package main

import (
	"fmt"
	"runtime"

	"github.com/kffl/gocannon/common"
)

func printHeader(cfg common.Config) {
	fmt.Printf("Attacking %s with %d connections over %s using %d CPUs\n", *cfg.Target, *cfg.Connections, *cfg.Duration, runtime.GOMAXPROCS(0))
}

func printSummary(s TestResults) {
	fmt.Printf("Total Req:  %8d\n", s.GetReqCount())
	fmt.Printf("Req/s:      %11.2f\n", s.GetReqPerSec())
}
