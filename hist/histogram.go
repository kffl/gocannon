package hist

import (
	"bufio"
	"errors"
	"os"
	"sync/atomic"
	"time"
)

type histogram []int64

type summary struct {
	reqCount           int64
	reqPerSec          float64
	latencyAvg         float64
	latencyPercentiles []int64
}

type requestHist struct {
	data    histogram
	bins    int64
	results summary
}

func NewRequestHist(timeout time.Duration) requestHist {
	bins := int64(float64(timeout) / 1000.0)
	r := requestHist{}
	r.data = make(histogram, bins)
	r.bins = bins
	return r
}

func (h *requestHist) RecordResponse(conn int, code int, start int64, end int64) {
	index := int64(float64(end-start)/1000.0 + 0.5)
	if index < h.bins {
		atomic.AddInt64(&((*h).data[index]), 1)
	}
}

func (h *requestHist) CalculateStats(
	start int64,
	stop int64,
	interval time.Duration,
) {
	percentiles, count := h.data.calculatePercentilesAndCount([]float64{50., 75., 90., 99.})
	avg := h.data.calculateAvg()
	reqPerSec := float64(count) / float64((stop-start)/int64(time.Second))

	h.results = summary{
		count,
		reqPerSec,
		avg,
		percentiles,
	}
}

func (h *requestHist) PrintReport() {
	h.results.print()
}

func (h *requestHist) SaveRawData(fileName string) error {
	f, err := os.Create(fileName)

	if err != nil {
		return errors.New("error creating output file")
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	err = writer.Flush()

	return err
}

func (h *requestHist) GetReqCount() int64 {
	return h.results.reqCount
}

func (h *requestHist) GetReqPerSec() float64 {
	return h.results.reqPerSec
}

func (h *requestHist) GetLatencyAvg() float64 {
	return h.results.latencyAvg * 1000.0
}
