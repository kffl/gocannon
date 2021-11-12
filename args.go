package main

import (
	"os"

	"github.com/kffl/gocannon/common"
	"gopkg.in/alecthomas/kingpin.v2"
)

func parseRequestBody(s kingpin.Settings) *common.RawRequestBody {
	r := &common.RawRequestBody{}
	s.SetValue((*common.RawRequestBody)(r))
	return r
}

func parseRequestHeaders(s kingpin.Settings) *common.RequestHeaders {
	r := &common.RequestHeaders{}
	s.SetValue((*common.RequestHeaders)(r))
	return r
}

var app = kingpin.New("gocannon", "Performance-focused HTTP benchmarking tool")

var config = common.Config{
	Duration: app.Flag("duration", "Load test duration").
		Short('d').
		Default("10s").
		Duration(),
	Connections: app.Flag("connections", "Maximum number of concurrent connections").
		Short('c').
		Default("50").
		Int(),
	Timeout: app.Flag("timeout", "HTTP client timeout").
		Short('t').
		Default("200ms").
		Duration(),
	Mode: app.Flag("mode", "Statistics collection mode: reqlog (logs each request) or hist (stores histogram of completed requests latencies)").
		Default("reqlog").
		Short('m').
		String(),
	OutputFile: app.Flag("output", "File to save the request log in CSV format (reqlog mode) or a text file with raw histogram data (hist mode)").
		PlaceHolder("file.csv").
		Short('o').
		String(),
	Interval: app.Flag("interval", "Interval for statistics calculation (reqlog mode)").
		Default("250ms").
		Short('i').
		Duration(),
	Preallocate: app.Flag("preallocate", "Number of requests in req log to preallocate memory for per connection (reqlog mode)").
		Default("1000").
		Int(),
	Method:   app.Flag("method", "The HTTP request method (GET, POST, PUT, PATCH or DELETE)").Default("GET").Enum("GET", "POST", "PUT", "PATCH", "DELETE"),
	Body:     parseRequestBody(app.Flag("body", "HTTP request body").Short('b').PlaceHolder("\"{data...\"")),
	Headers:  parseRequestHeaders(kingpin.Flag("header", "HTTP request header(s). You can set more than one header by repeating this flag.").Short('h').PlaceHolder("\"k:v\"")),
	TrustAll: app.Flag("trust-all", "Omit SSL certificate validation").Bool(),
	Plugin:   app.Flag("plugin", "Plugin to run Gocannon with").PlaceHolder("/to/p.so").ExistingFile(),
	Target:   app.Arg("target", "HTTP target URL with port (i.e. http://localhost:80/test or https://host:443/x)").Required().String(),
}

func parseArgs() error {
	app.Version("1.0.0")
	_, err := app.Parse(os.Args[1:])
	return err
}
