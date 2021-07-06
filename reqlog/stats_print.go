package reqlog

import (
	"fmt"
	"time"
)

func printStatsHeader() {
	fmt.Println("|--REQS--|    |------------------------LATENCY-------------------------|")
	fmt.Println("     Count           AVG         P50         P75         P90         P99")
}

func formatLatency(latency float64) string {
	d := time.Duration(latency) * time.Nanosecond
	return d.String()
}

func formatLatencyI64(latency int64) string {
	d := time.Duration(latency) * time.Nanosecond
	return d.String()
}

func (s *statistics) print() {
	fmt.Printf("%10d %13v", s.count, formatLatency(s.latencyAVG))
	for _, v := range s.latencyPercentiles {
		fmt.Printf(" %11v", formatLatencyI64(v))
	}
	fmt.Printf("\n")
}

func (s *fullStatistics) Print() {
	fmt.Printf("Total Req: %8d\n", s.reqCount)
	fmt.Printf("Req/s:     %11.2f\n", s.reqPerSec)

	fmt.Printf("Interval stats: (interval = %v) \n", s.interval)
	printStatsHeader()

	for _, stats := range s.detailed {
		stats.print()
	}

	fmt.Println("----------")

	s.summary.print()
}

func (s *fullStatistics) GetReqCount() int {
	return s.reqCount
}
