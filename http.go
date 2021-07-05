package main

import (
	"errors"
	"net"
	"net/url"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	ErrWrongTarget         = errors.New("wrong target URL")
	ErrUnsupportedProtocol = errors.New("unsupported target protocol")
)

func newHTTPClient(
	target string,
	timeout time.Duration,
	connections int,
) (*fasthttp.HostClient, error) {
	c := new(fasthttp.HostClient)
	u, err := url.ParseRequestURI(target)
	if err != nil {
		return nil, ErrWrongTarget
	}
	if u.Scheme != "http" {
		return nil, ErrUnsupportedProtocol
	}
	c = &fasthttp.HostClient{
		Addr:                          u.Host,
		MaxConns:                      int(connections),
		ReadTimeout:                   timeout,
		WriteTimeout:                  timeout,
		DisableHeaderNamesNormalizing: true,
		Dial: func(addr string) (net.Conn, error) {
			return fasthttp.DialTimeout(addr, timeout)
		},
	}
	return c, nil
}

func performRequest(c *fasthttp.HostClient, target string) (
	code int, start int64, end int64,
) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.URI().SetScheme("http")
	req.Header.SetMethod("GET")
	req.SetRequestURI(target)

	start = makeTimestamp()
	err := c.Do(req, resp)
	if err != nil {
		code = -1
	} else {
		code = resp.StatusCode()
	}
	end = makeTimestamp()

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)

	return
}
