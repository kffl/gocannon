package main

import (
	"fmt"

	"github.com/kffl/gocannon/common"
)

func printHeader(cfg common.Config) {
	fmt.Printf("Attacking %s with %d connections over %s\n", *cfg.Target, *cfg.Connections, *cfg.Duration)
}

func printSummary(s TestResults) {
	fmt.Printf("Total Req:  %8d\n", s.GetReqCount())
	fmt.Printf("Req/s:      %11.2f\n", s.GetReqPerSec())
}
