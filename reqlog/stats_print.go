package reqlog

import (
	"fmt"
)

func printStatsHeader() {
	fmt.Println("|--REQS--|    |------------------------LATENCY-------------------------|")
	fmt.Println("     Count           AVG         P50         P75         P90         P99")
}

var durationUnits = []string{"s ", "ms", "μs", "ns"}
var durationValues = []int64{1000000000, 1000000, 1000, 1}

func formatDuration(d float64) string {
	if d == -1 {
		return "-"
	}
	for i, unit := range durationUnits {
		value := durationValues[i]
		if int64(d)/value > 0 {
			return fmt.Sprintf("%.4f%s", d/float64(value), unit)
		}
	}
	return fmt.Sprintf("%.4f%s", d/float64(1), "ns")
}

func formatLatency(latency float64) string {
	return formatDuration(latency)
}

func formatLatencyI64(latency int64) string {
	return formatDuration(float64(latency))
}

func (s *statistics) print() {
	fmt.Printf("%10d %13v", s.Count, formatLatency(s.LatencyAVG))
	for _, v := range s.LatencyPercentiles {
		fmt.Printf(" %11v", formatLatencyI64(v))
	}
	fmt.Printf("\n")
}

func (s *fullStatistics) print() {
	fmt.Printf("Interval stats: (interval = %v) \n", s.Interval)
	printStatsHeader()

	for _, stats := range s.Detailed {
		stats.print()
	}

	fmt.Println("----------")

	s.Summary.print()
}

func (s *fullStatistics) GetReqCount() int64 {
	return int64(s.Summary.Count)
}
