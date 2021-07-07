package reqlog

import (
	"sort"
	"time"
)

type request struct {
	code  int
	start int64
	end   int64
}

type requestLog [][]request

type flatRequestLog []request

type requestLogCollector struct {
	reqLog         *requestLog
	flatLog        *flatRequestLog
	results        *fullStatistics
	startTimestamp int64
}

func newRequests(n int, preallocate int) requestLog {
	reqs := make(requestLog, n)
	for i := 0; i < n; i++ {
		reqs[i] = make([]request, 0, preallocate)
	}

	return reqs
}

func NewRequestLog(n int, preallocate int) *requestLogCollector {
	reqs := newRequests(n, preallocate)

	logCollector := requestLogCollector{}

	logCollector.reqLog = &reqs

	return &logCollector
}

func (r *requestLogCollector) RecordResponse(conn int, code int, start int64, end int64) {
	(*r.reqLog)[conn] = append((*r.reqLog)[conn], request{code, start, end})
}

func (r *requestLogCollector) CalculateStats(
	start int64,
	stop int64,
	interval time.Duration,
) {
	r.startTimestamp = start
	reqsFlat := r.reqLog.flatten()
	reqsFlat.sort()

	results := reqsFlat.calculateStats(start, stop, interval)
	r.results = &results

	// keep the flattened, sorted request log in case the output is to be saved
	r.flatLog = &reqsFlat
	// while releasing the connection-paritioned one to the GC
	*(r.reqLog) = nil

}

func (r *requestLogCollector) PrintReport() {
	r.results.print()
}

func (r *requestLogCollector) SaveRawData(outputFile string) error {
	if outputFile != "" {
		return r.flatLog.saveCSV(r.startTimestamp, outputFile)
	}
	return nil

}

func (r *requestLogCollector) GetReqCount() int64 {
	return r.results.reqCount
}

func (r *requestLogCollector) GetReqPerSec() float64 {
	return r.results.reqPerSec
}

func (r *requestLogCollector) GetLatencyAvg() float64 {
	return r.results.summary.latencyAVG
}

func (reqs requestLog) flatten() flatRequestLog {
	flattened := make(flatRequestLog, 0, 50000)

	for i := 0; i < len(reqs); i++ {
		flattened = append(flattened, reqs[i]...)
	}

	return flattened
}

func (reqs flatRequestLog) sort() {
	sort.Slice(reqs, func(x, y int) bool {
		return reqs[x].end < reqs[y].end
	})
}
