package hist

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRequestHist(t *testing.T) {
	h := NewRequestHist(time.Millisecond)

	assert.Len(
		t,
		h.data,
		1000,
		"should have length equal to number of microseconds in the specified timeout",
	)
}

func TestRecordResponse(t *testing.T) {
	h := NewRequestHist(time.Microsecond * 200)

	h.RecordResponse(0, 200, 10, 100010)
	h.RecordResponse(0, 200, 30, 100030)
	h.RecordResponse(0, 200, 30000, 150000)

	assert.Equal(t, int64(2), h.data[100], "there should be two hits registered @100μs")
	assert.Equal(t, int64(1), h.data[120], "there should be one hit registered @120μs")
	assert.Zero(t, h.data[0])
	assert.Zero(t, h.data[10])
}

func TestCalculateStats(t *testing.T) {
	h := NewRequestHist(time.Microsecond * 200)

	h.RecordResponse(0, 200, 10, 100010)     // 100μs
	h.RecordResponse(0, 200, 30, 100030)     // 100μs
	h.RecordResponse(0, 200, 30000, 150000)  // 120μs
	h.RecordResponse(0, 200, 100000, 150000) // 50μs
	h.RecordResponse(0, 200, 1000, 2005)     // 1μs

	h.CalculateStats(10, int64(time.Second)*10+10, time.Millisecond)

	assert.Equal(t, int64(5), h.GetReqCount())
	assert.Equal(t, 0.5, h.GetReqPerSec())
}

func TestCalculateAvg(t *testing.T) {
	data := histogram{0, 0, 0, 1, 5, 8, 1, 2, 0, 3}

	assert.Equal(t, 5.5, data.calculateAvg())
}

func TestCalculatePercentilesAndCount(t *testing.T) {
	data := histogram{0, 0, 0, 1, 5, 8, 1, 2, 0, 3}

	percentiles, count := data.calculatePercentilesAndCount([]float64{50, 75, 90, 99})

	assert.Equal(t, int64(20), count)
	assert.Len(t, percentiles, 4)
	assert.ElementsMatch(t, percentiles, []int64{5, 6, 9, 9})
}

func TestWriteHistData(t *testing.T) {
	h := NewRequestHist(time.Microsecond * 14)
	h.data = histogram{0, 0, 0, 1, 5, 8, 1, 2, 0, 0, 0, 0, 0, 0}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	h.writeHistData(w)

	w.Flush()

	assert.Equal(t, "0\n0\n0\n1\n5\n8\n1\n2\n", b.String())
}
