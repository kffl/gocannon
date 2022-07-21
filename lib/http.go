package lib

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kffl/gocannon/common"
	"github.com/valyala/fasthttp"
)

var (
	ErrWrongTarget         = errors.New("wrong target URL")
	ErrUnsupportedProtocol = errors.New("unsupported target protocol")
	ErrMissingPort         = errors.New("missing target port")
)

func parseTarget(target string) (scheme string, host string, err error) {
	u, err := url.ParseRequestURI(target)
	if err != nil {
		err = ErrWrongTarget
		return
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		err = ErrUnsupportedProtocol
		return
	}
	tokenizedHost := strings.Split(u.Host, ":")
	port := tokenizedHost[len(tokenizedHost)-1]
	if _, err = strconv.Atoi(port); err != nil {
		err = ErrMissingPort
		return
	}
	return u.Scheme, u.Host, nil
}

func dialHost(host string, timeout time.Duration) error {
	conn, err := fasthttp.DialTimeout(host, timeout)
	if err != nil {
		return fmt.Errorf("dialing host failed (%w)", err)
	}
	conn.Close()
	return nil
}

func newHTTPClient(
	target string,
	timeout time.Duration,
	connections int,
	trustAll bool,
	checkHost bool,
) (*fasthttp.HostClient, error) {
	scheme, host, err := parseTarget(target)
	if err != nil {
		return nil, err
	}
	c := &fasthttp.HostClient{
		Addr:                          host,
		MaxConns:                      int(connections),
		ReadTimeout:                   timeout,
		WriteTimeout:                  timeout,
		DisableHeaderNamesNormalizing: true,
		Dial: func(addr string) (net.Conn, error) {
			return fasthttp.DialTimeout(addr, timeout)
		},
		IsTLS: scheme == "https",
		TLSConfig: &tls.Config{
			InsecureSkipVerify: trustAll,
		},
	}
	if checkHost {
		err = dialHost(host, timeout)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func performRequest(c *fasthttp.HostClient, target string, method string, body []byte, headers common.RequestHeaders) (
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
		req.Header.Add(h.Key, h.Value)
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
