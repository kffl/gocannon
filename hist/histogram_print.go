package hist

import "fmt"

func (r *summary) print() {
	fmt.Println("|------------------------LATENCY (μs)-----------------------|")
	fmt.Println("          AVG         P50         P75         P90         P99")
	fmt.Printf("%13.2f", r.latencyAvg)
	for _, v := range r.latencyPercentiles {
		fmt.Printf(" %11v", v)
	}
	fmt.Printf("\n")
}
