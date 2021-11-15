package reqlog

import (
	"sort"
	"time"
)

var percentiles = []float64{50, 75, 90, 99}

type requestLatencies []int64

type statistics struct {
	Count              int
	LatencyAVG         float64
	LatencyPercentiles []int64
	ReqPerSec          float64
}

type intervalStatistics []statistics

type fullStatistics struct {
	Summary  statistics
	Detailed intervalStatistics
	Interval time.Duration
}

func (sortedReqs *flatRequestLog) calculateStats(
	start int64,
	stop int64,
	intervalDuration time.Duration,
) fullStatistics {
	summaryStats := sortedReqs.calculateIntervalStats(stop - start)

	var detailedStats intervalStatistics

	for intervalStart := start; intervalStart < stop; intervalStart += int64(intervalDuration) {
		startIndex := sort.Search(len(*sortedReqs), func(i int) bool {
			return (*sortedReqs)[i].end >= intervalStart
		})

		endIndex := sort.Search(len(*sortedReqs), func(i int) bool {
			return (*sortedReqs)[i].end >= intervalStart+int64(intervalDuration)
		})

		slicedRequests := (*sortedReqs)[startIndex:endIndex]

		detailedStats = append(
			detailedStats,
			slicedRequests.calculateIntervalStats(int64(intervalDuration)),
		)
	}

	return fullStatistics{summaryStats, detailedStats, intervalDuration}
}

func (sortedReqs *flatRequestLog) calculateIntervalStats(timespan int64) statistics {
	latencies := make(requestLatencies, 0, len(*sortedReqs))

	for i := 0; i < len(*sortedReqs); i++ {
		latencies = append(latencies, (*sortedReqs)[i].end-(*sortedReqs)[i].start)
	}

	var r statistics

	latencies.sort()

	for _, p := range percentiles {
		r.LatencyPercentiles = append(r.LatencyPercentiles, latencies.calculatePercentile(p))
	}

	c := len(latencies)
	r.Count = c
	r.LatencyAVG = latencies.calculateAVG()
	r.ReqPerSec = float64(c) / (float64(timespan) / float64(time.Second))

	return r
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
