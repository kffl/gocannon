package hist

import "sort"

func (h *histogram) calculateAvg() float64 {
	var sum int64 = 0
	var count int64 = 0

	for i, el := range *h {
		sum += int64(el) * int64(i)
		count += el
	}

	return float64(sum) / float64(count)
}

func (h *histogram) calculatePercentilesAndCount(percents []float64) ([]int64, int64) {
	cumulativeHist := make(histogram, len(*h))
	results := make([]int64, len(percents))

	var count int64 = 0

	for i, v := range *h {
		count += v
		cumulativeHist[i] = count
	}

	for i, percent := range percents {
		target := int64(percent / 100.0 * float64(count))

		index := sort.Search(len(cumulativeHist), func(i int) bool {
			return cumulativeHist[i] >= target
		})

		results[i] = int64(index)
	}

	return results, count
}
