package reqlog

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/kffl/gocannon/rescodes"
	"gopkg.in/yaml.v2"
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
	resCodes       *rescodes.Rescodes
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
	outputFile string,
) error {
	r.startTimestamp = start

	// save raw request data before flattening the request log
	err := r.saveRawData(outputFile)

	r.resCodes = rescodes.NewRescodes()
	r.saveResCodes()

	reqsFlat := r.reqLog.flatten()
	reqsFlat.sort()

	results := reqsFlat.calculateStats(start, stop, interval)
	r.results = &results

	// keep the flattened, sorted request log
	r.flatLog = &reqsFlat
	// while releasing the connection-paritioned one to the GC
	*(r.reqLog) = nil

	return err
}

func (r *requestLogCollector) PrintReport(format string) {
	if format == "default" {
		r.results.print()
		r.resCodes.PrintRescodes()
	} else {
		obj := struct {
			Report   *fullStatistics
			ResCodes map[int]int64
		}{
			r.results,
			r.resCodes.AsMap(),
		}
		var output []byte
		if format == "json" {
			output, _ = json.MarshalIndent(obj, "", "  ")
		}
		if format == "yaml" {
			output, _ = yaml.Marshal(obj)
		}
		fmt.Printf("%s", output)
	}
}

func (r *requestLogCollector) saveRawData(outputFile string) error {
	if outputFile != "" {
		return r.reqLog.saveCSV(r.startTimestamp, outputFile)
	}
	return nil

}

func (r *requestLogCollector) GetReqCount() int64 {
	return int64(r.results.Summary.Count)
}

func (r *requestLogCollector) GetReqPerSec() float64 {
	return r.results.Summary.ReqPerSec
}

func (r *requestLogCollector) GetLatencyAvg() float64 {
	return r.results.Summary.LatencyAVG
}

func (r *requestLogCollector) saveResCodes() {
	for _, connLog := range *r.reqLog {
		for _, req := range connLog {
			r.resCodes.RecordRequest(req.code)
		}
	}
}

func (reqs requestLog) flatten() flatRequestLog {
	flattened := make(flatRequestLog, 0, 50000)

	for i := 0; i < len(reqs); i++ {
		for j := 0; j < len(reqs[i]); j++ {
			if reqs[i][j].code != 0 {
				flattened = append(flattened, reqs[i][j])
			}
		}
	}

	return flattened
}

func (reqs flatRequestLog) sort() {
	sort.Slice(reqs, func(x, y int) bool {
		return reqs[x].end < reqs[y].end
	})
}
