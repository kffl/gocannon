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
	u, err := url.ParseRequestURI(target)
	if err != nil {
		return nil, ErrWrongTarget
	}
	if u.Scheme != "http" {
		return nil, ErrUnsupportedProtocol
	}
	c := &fasthttp.HostClient{
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

func performRequest(c *fasthttp.HostClient, target string, method string, body []byte, headers requestHeaders) (
	code int, start int64, end int64,
) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.URI().SetScheme("http")
	req.Header.SetMethod(method)
	req.SetRequestURI(target)

	req.SetBodyRaw(body)

	for _, h := range headers {
		req.Header.Add(h.key, h.value)
	}

	start = makeTimestamp()
	err := c.Do(req, resp)
	if err != nil {
		code = 0
	} else {
		code = resp.StatusCode()
	}
	end = makeTimestamp()

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)

	return
}
