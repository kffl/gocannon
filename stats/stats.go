package stats

import (
	"sort"
	"time"
)

var percentiles = []float64{50, 75, 90, 99}

type requestLatencies []int64

type statistics struct {
	count              int
	latencyAVG         float64
	latencyPercentiles []int64
}

type intervalStatistics []statistics

type fullStatistics struct {
	summary   statistics
	detailed  intervalStatistics
	interval  time.Duration
	reqCount  int
	reqPerSec float64
}

func calculateStats(
	sortedReqs flattenedRequests,
	start int64,
	stop int64,
	intervalDuration time.Duration,
) fullStatistics {
	summaryStats := calculateIntervalStats(sortedReqs)

	var detailedStats intervalStatistics

	for intervalStart := start; intervalStart < stop; intervalStart += int64(intervalDuration) {
		startIndex := sort.Search(len(sortedReqs), func(i int) bool {
			return sortedReqs[i].end >= intervalStart
		})

		endIndex := sort.Search(len(sortedReqs), func(i int) bool {
			return sortedReqs[i].end >= intervalStart+int64(intervalDuration)
		})

		detailedStats = append(
			detailedStats,
			calculateIntervalStats(sortedReqs[startIndex:endIndex]),
		)
	}

	reqCount := len(sortedReqs)
	reqPerSec := float64(reqCount) / float64((stop-start)/int64(time.Second))

	return fullStatistics{summaryStats, detailedStats, intervalDuration, reqCount, reqPerSec}
}

func calculateIntervalStats(reqs flattenedRequests) statistics {
	latencies := make(requestLatencies, 0, len(reqs))

	for i := 0; i < len(reqs); i++ {
		latencies = append(latencies, reqs[i].end-reqs[i].start)
	}

	var r statistics

	latencies.sort()

	for _, p := range percentiles {
		r.latencyPercentiles = append(r.latencyPercentiles, latencies.calculatePercentile(p))
	}

	r.count = len(latencies)
	r.latencyAVG = latencies.calculateAVG()

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
