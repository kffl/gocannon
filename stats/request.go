package stats

import (
	"sort"
	"time"
)

type request struct {
	code  int
	start int64
	end   int64
}

type requests [][]request

type flattenedRequests []request

func NewRequests(n int) requests {
	reqs := make(requests, n)
	for i := 0; i < n; i++ {
		reqs[i] = make([]request, 0, 1000)
	}

	return reqs
}

func (reqs *requests) RecordResponse(conn int, code int, start int64, end int64) {
	(*reqs)[conn] = append((*reqs)[conn], request{code, start, end})
}

func (reqs *requests) CalculateStats(
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

func (reqs requests) flatten() flattenedRequests {
	flattened := make(flattenedRequests, 0, 50000)

	for i := 0; i < len(reqs); i++ {
		flattened = append(flattened, reqs[i]...)
	}

	return flattened
}

func (reqs flattenedRequests) sort() {
	sort.Slice(reqs, func(x, y int) bool {
		return reqs[x].end < reqs[y].end
	})
}
