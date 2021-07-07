package main

import "fmt"

func printHeader() {
	fmt.Printf("Attacking %s with %d connections over %s\n", *target, *connections, *duration)
	fmt.Printf("gocannon goes brr...\n")
}

func printSummary(s statsCollector) {
	fmt.Printf("Total Req:  %8d\n", s.GetReqCount())
	fmt.Printf("Req/s:      %11.2f\n", s.GetReqPerSec())
}
