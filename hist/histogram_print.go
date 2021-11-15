package hist

import "fmt"

func (r *summary) print() {
	fmt.Println("|------------------------LATENCY (μs)-----------------------|")
	fmt.Println("          AVG         P50         P75         P90         P99")
	fmt.Printf("%13.2f", r.LatencyAvg)
	for _, v := range r.LatencyPercentiles {
		fmt.Printf(" %11v", v)
	}
	fmt.Printf("\n")
}
