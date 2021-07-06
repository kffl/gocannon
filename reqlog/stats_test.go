package reqlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePercentile(t *testing.T) {
	r := requestLatencies{10, 20, 30}
	r2 := requestLatencies{10, 20}

	assert.Equal(t, int64(10), r.calculatePercentile(14.0))
	assert.Equal(t, int64(30), r.calculatePercentile(99.0))
	assert.Equal(t, int64(20), r.calculatePercentile(50.0))
	assert.Equal(t, int64(20), r2.calculatePercentile(50.0))
}

func TestCalculateAVG(t *testing.T) {
	r := requestLatencies{10, 20, 30, 50}
	r2 := requestLatencies{}

	assert.Equal(t, float64(27.5), r.calculateAVG())
	assert.Equal(t, float64(-1), r2.calculateAVG())
}

func TestSort(t *testing.T) {
	r := requestLatencies{}
	r2 := requestLatencies{5, 10, 20, 3, 5}

	r.sort()
	r2.sort()

	assert.Len(t, r, 0, "should be a zero-length slice")
	assert.IsNonDecreasing(t, r2, "should be in non-decreasing order")
}

func TestCalculateIntervalStatsEmpty(t *testing.T) {
	reqs := flatRequestLog{}

	stats := reqs.calculateIntervalStats()

	assert.Equal(t, 0, stats.count)
	assert.Equal(t, float64(-1), stats.latencyAVG)
	assert.ElementsMatch(t, stats.latencyPercentiles, []int64{-1, -1, -1, -1})
}

func TestCalculateIntervalStatsPopulated(t *testing.T) {
	reqs := flatRequestLog{
		{200, 123, 223},
		{200, 123, 223},
		{200, 234, 235},
		{200, 234, 235},
		{200, 534, 535},
	}

	stats := reqs.calculateIntervalStats()

	assert.Equal(t, 5, stats.count)
	assert.Equal(t, float64(40.6), stats.latencyAVG)
	assert.ElementsMatch(t, stats.latencyPercentiles, []int64{1, 100, 100, 100})
}

func TestCalculateStats(t *testing.T) {
	reqs := flatRequestLog{
		{200, 123, 223},
		{200, 123, 223},
		{200, 234, 235},
		{200, 234, 235},
		{200, 234, 535},
	}

	full := reqs.calculateStats(100, 600, 100)
	summary := full.summary
	detailed := full.detailed

	assert.Equal(t, 5, summary.count)
	assert.Len(t, detailed, 5, "5 intervals should fit in the specified range")
	assert.Equal(
		t,
		detailed[4].count,
		1,
		"request should be assigned to an interval based on the response timestamp",
	)
}
