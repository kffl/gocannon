package main

import (
	"fmt"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

type rawRequestBody []byte

func (b *rawRequestBody) Set(value string) error {
	(*b) = []byte(value)
	return nil
}

func (b *rawRequestBody) String() string {
	return fmt.Sprint(*b)
}

func (b *rawRequestBody) IsCumulative() bool {
	return false
}

func parseRequestBody(s kingpin.Settings) *rawRequestBody {
	r := &rawRequestBody{}
	s.SetValue((*rawRequestBody)(r))
	return r
}

type requestHeader struct {
	key   string
	value string
}

type requestHeaders []requestHeader

func (r *requestHeaders) Set(value string) error {
	tokenized := strings.Split(value, ":")
	if len(tokenized) != 2 {
		return fmt.Errorf("Header '%s' doesn't match 'Key:Value' format (i.e. 'Content-Type:application/json')", value)
	}
	h := requestHeader{tokenized[0], tokenized[1]}
	(*r) = append(*r, h)
	return nil
}

func (r *requestHeaders) String() string {
	return fmt.Sprint(*r)
}

func (r *requestHeaders) IsCumulative() bool {
	return true
}

func parseRequestHeaders(s kingpin.Settings) *requestHeaders {
	r := &requestHeaders{}
	s.SetValue((*requestHeaders)(r))
	return r
}

var (
	duration = kingpin.Flag("duration", "Load test duration").
			Short('d').
			Default("10s").
			Duration()
	connections = kingpin.Flag("connections", "Maximum number of concurrent connections").
			Short('c').
			Default("50").
			Int()
	timeout = kingpin.Flag("timeout", "HTTP client timeout").
		Short('t').
		Default("200ms").
		Duration()
	mode = kingpin.Flag("mode", "Statistics collection mode: reqlog (logs each request) or hist (stores histogram of completed requests latencies)").
		Default("reqlog").
		Short('m').
		String()
	outputFile = kingpin.Flag("output", "File to save the request log in CSV format (reqlog mode) or a text file with raw histogram data (hist mode)").
			PlaceHolder("file.csv").
			Short('o').
			String()
	interval = kingpin.Flag("interval", "Interval for statistics calculation (reqlog mode)").
			Default("250ms").
			Short('i').
			Duration()
	preallocate = kingpin.Flag("preallocate", "Number of requests in req log to preallocate memory for per connection (reqlog mode)").
			Default("1000").
			Int()
	method   = kingpin.Flag("method", "The HTTP request method (GET, POST, PUT, PATCH or DELETE)").Default("GET").Enum("GET", "POST", "PUT", "PATCH", "DELETE")
	body     = parseRequestBody(kingpin.Flag("body", "HTTP request body").Short('b').PlaceHolder("\"{data...\""))
	headers  = parseRequestHeaders(kingpin.Flag("header", "HTTP request header(s). You can set more than one header by repeating this flag.").Short('h').PlaceHolder("\"k:v\""))
	trustAll = kingpin.Flag("trust-all", "Omit SSL certificate validation").Bool()
	target   = kingpin.Arg("target", "HTTP target URL with port (i.e. http://localhost:80/test or https://host:443/x)").Required().String()
)

func parseArgs() {
	kingpin.Version("0.2.1")
	kingpin.Parse()
}
