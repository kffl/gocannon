package hist

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/kffl/gocannon/rescodes"
	"gopkg.in/yaml.v2"
)

type histogram []int64

type summary struct {
	ReqCount           int64
	ReqPerSec          float64
	LatencyAvg         float64
	LatencyPercentiles []int64
}

type requestHist struct {
	data      histogram
	bins      int64
	didNotFit int64
	resCodes  *rescodes.Rescodes
	results   summary
}

func NewRequestHist(timeout time.Duration) requestHist {
	bins := int64(float64(timeout) / 1000.0)
	r := requestHist{}
	r.data = make(histogram, bins)
	r.bins = bins
	r.resCodes = rescodes.NewRescodes()
	return r
}

func (h *requestHist) RecordResponse(conn int, code int, start int64, end int64) {
	if code == 0 {
		h.resCodes.RecordRequestThreadSafe(0)
	} else {
		index := int64(float64(end-start)/1000.0 + 0.5)
		if index < h.bins {
			atomic.AddInt64(&((*h).data[index]), 1)
			h.resCodes.RecordRequestThreadSafe(code)
		} else {
			atomic.AddInt64(&h.didNotFit, 1)
		}
	}
}

func (h *requestHist) CalculateStats(
	start int64,
	stop int64,
	interval time.Duration,
	outputFile string,
) error {
	percentiles, count := h.data.calculatePercentilesAndCount([]float64{50., 75., 90., 99.})
	avg := h.data.calculateAvg()
	reqPerSec := float64(count) / float64((stop-start)/int64(time.Second))

	h.results = summary{
		count,
		reqPerSec,
		avg,
		percentiles,
	}

	if outputFile != "" {
		return h.saveRawData(outputFile)
	}

	return nil
}

func (h *requestHist) PrintReport(format string) {
	if format == "default" {
		h.results.print()
		if h.didNotFit > 0 {
			fmt.Fprintf(
				os.Stderr,
				"WARNING: some recorded responses (%d) did not fit in the histogram potentially skewing the resulting stats. Consider increasing timeout duration.\n",
				h.didNotFit,
			)
		}
		h.resCodes.PrintRescodes()
	} else {
		obj := struct {
			Report    *summary
			DidNotFit int64
			ResCodes  map[int]int64
		}{
			&h.results,
			h.didNotFit,
			h.resCodes.AsMap(),
		}
		var output []byte
		if format == "json" {
			output, _ = json.MarshalIndent(obj, "", " ")
		}
		if format == "yaml" {
			output, _ = yaml.Marshal(obj)
		}
		fmt.Printf("%s", output)
	}
}

func (h *requestHist) saveRawData(fileName string) error {
	f, err := os.Create(fileName)

	if err != nil {
		return errors.New("error creating output file")
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	err = h.writeHistData(writer)

	if err != nil {
		return err
	}

	err = writer.Flush()

	return err
}

func (h *requestHist) GetReqCount() int64 {
	return h.results.ReqCount
}

func (h *requestHist) GetReqPerSec() float64 {
	return h.results.ReqPerSec
}

func (h *requestHist) GetLatencyAvg() float64 {
	return h.results.LatencyAvg * 1000.0
}

func (h *requestHist) GetLatencyPercentiles() []int64 {
	asNanoseconds := make([]int64, len(h.results.LatencyPercentiles))
	for i, p := range h.results.LatencyPercentiles {
		asNanoseconds[i] = p * 1000
	}
	return asNanoseconds
}
