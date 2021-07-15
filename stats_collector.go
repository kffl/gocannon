package main

import (
	"fmt"
	"time"

	"github.com/kffl/gocannon/hist"
	"github.com/kffl/gocannon/reqlog"
)

type statsCollector interface {
	RecordResponse(conn int, code int, start int64, end int64)
	CalculateStats(start int64, stop int64, interval time.Duration, fileName string) error
	PrintReport()
	GetReqCount() int64
	GetReqPerSec() float64
	GetLatencyAvg() float64
}

func newStatsCollector(
	mode string,
	conns int,
	preallocate int,
	timeout time.Duration,
) (statsCollector, error) {
	switch mode {
	case "reqlog":
		r := reqlog.NewRequestLog(conns, preallocate)
		return r, nil
	case "hist":

		r := hist.NewRequestHist(timeout)
		return &r, nil
	}

	return nil, fmt.Errorf("wrong mode '%s'", mode)
}
