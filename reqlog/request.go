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

func NewRequests(n int, preallocate int) requestLog {
	reqs := make(requestLog, n)
	for i := 0; i < n; i++ {
		reqs[i] = make([]request, 0, preallocate)
	}

	return reqs
}

func (reqs *requestLog) RecordResponse(conn int, code int, start int64, end int64) {
	(*reqs)[conn] = append((*reqs)[conn], request{code, start, end})
}

func (reqs *requestLog) CalculateStats(
	start int64,
	stop int64,
	interval time.Duration,
	outputFile string,
) (fullStatistics, error) {
	reqsFlat := reqs.flatten()
	reqsFlat.sort()

	s := calculateStats(reqsFlat, start, stop, interval)

	if outputFile != "" {
		return s, reqsFlat.saveCSV(start, outputFile)
	}

	return s, nil
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
