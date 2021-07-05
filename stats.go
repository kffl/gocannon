package main

import (
	"fmt"
	"sort"
	"time"
)

var percentiles = []float64{50, 75, 90, 99}

type requestLatencies []int64

type stats struct {
	count              int
	latencyAVG         float64
	latencyPercentiles []int64
}

type intervalStats []stats

func calculateStats(
	sortedReqs flattenedRequests,
	start int64,
	stop int64,
	intervalDuration int64,
) (summary stats, detailed intervalStats) {
	summaryStats := calculateIntervalStats(sortedReqs)

	var detailedStats intervalStats

	for intervalStart := start; intervalStart < stop; intervalStart += intervalDuration {
		startIndex := sort.Search(len(sortedReqs), func(i int) bool {
			return sortedReqs[i].end >= intervalStart
		})

		endIndex := sort.Search(len(sortedReqs), func(i int) bool {
			return sortedReqs[i].end >= intervalStart+intervalDuration
		})

		detailedStats = append(
			detailedStats,
			calculateIntervalStats(sortedReqs[startIndex:endIndex]),
		)
	}

	return summaryStats, detailedStats
}

func calculateIntervalStats(reqs flattenedRequests) stats {
	latencies := make(requestLatencies, 0, len(reqs))

	for i := 0; i < len(reqs); i++ {
		latencies = append(latencies, reqs[i].end-reqs[i].start)
	}

	var r stats

	latencies.sort()

	for _, p := range percentiles {
		r.latencyPercentiles = append(r.latencyPercentiles, latencies.calculatePercentile(p))
	}

	r.count = len(latencies)
	r.latencyAVG = latencies.calculateAVG()

	return r
}

func (s stats) printHeader() {
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

func (s stats) print() {
	fmt.Printf("%10d %13v", s.count, formatLatency(s.latencyAVG))
	for _, v := range s.latencyPercentiles {
		fmt.Printf(" %11v", formatLatencyI64(v))
	}
	fmt.Printf("\n")
}

func (latencies requestLatencies) sort() {
	sort.Slice(latencies, func(x, y int) bool {
		return latencies[x] < latencies[y]
	})
}

func (latencies requestLatencies) calculatePercentile(percent float64) int64 {
	if len(latencies) == 0 {
		return -1
	}

	index := int64((percent / 100) * float64(len(latencies)))

	return latencies[index]
}

func (latencies requestLatencies) calculateAVG() float64 {
	if len(latencies) == 0 {
		return -1
	}

	var sum int64 = 0

	for _, l := range latencies {
		sum += l
	}

	return float64(sum) / float64(len(latencies))
}
