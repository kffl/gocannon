package main

import (
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	ErrWrongTarget         = errors.New("wrong target URL")
	ErrUnsupportedProtocol = errors.New("unsupported target protocol")
	ErrMissingPort         = errors.New("missing target port")
)

func newHTTPClient(
	target string,
	timeout time.Duration,
	connections int,
	trustAll bool,
) (*fasthttp.HostClient, error) {
	u, err := url.ParseRequestURI(target)
	if err != nil {
		return nil, ErrWrongTarget
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, ErrUnsupportedProtocol
	}
	tokenizedHost := strings.Split(u.Host, ":")
	port := tokenizedHost[len(tokenizedHost)-1]
	if _, err := strconv.Atoi(port); err != nil {
		return nil, ErrMissingPort
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
		IsTLS: u.Scheme == "https",
		TLSConfig: &tls.Config{
			InsecureSkipVerify: trustAll,
		},
	}
	return c, nil
}

func performRequest(c *fasthttp.HostClient, target string, method string, body []byte, headers requestHeaders) (
	code int, start int64, end int64,
) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	if strings.HasPrefix(target, "https") {
		req.URI().SetScheme("https")
	} else {
		req.URI().SetScheme("http")
	}
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
