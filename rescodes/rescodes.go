package rescodes

import (
	"fmt"
	"sync/atomic"
)

type Rescodes [600]int64

func NewRescodes() *Rescodes {
	var r Rescodes
	return &r
}

func (r *Rescodes) RecordRequest(code int) {
	r[code]++
}

func (r *Rescodes) RecordRequestThreadSafe(code int) {
	atomic.AddInt64(&((*r)[code]), 1)
}

func (r *Rescodes) PrintRescodes() {
	fmt.Println("Responses by HTTP status code:")
	for code := 1; code < 600; code++ {
		hits := r[code]
		if hits > 0 {
			fmt.Printf("%5d  ->%8d\n", code, hits)
		}
	}
	if r[0] > 0 {
		fmt.Printf("Requests ended with timeout/socket error: %d\n", r[0])
	}
}
