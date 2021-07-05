package main

import "sort"

type request struct {
	code  int
	start int64
	end   int64
}

type requests [][]request

type flattenedRequests []request

func newRequests(n int) requests {
	reqs := make(requests, n)
	for i := 0; i < n; i++ {
		reqs[i] = make([]request, 0, 1000)
	}

	return reqs
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
