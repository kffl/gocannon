package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/kffl/gocannon/common"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPClientWrongURL(t *testing.T) {

	_, err1 := newHTTPClient("XYZthisisawrongurl123", time.Millisecond*200, 10, true, false)
	_, err2 := newHTTPClient("ldap://something/", time.Millisecond*200, 10, true, false)
	_, err3 := newHTTPClient("http://localhost/", time.Millisecond*200, 10, true, false)

	assert.ErrorIs(t, err1, ErrWrongTarget, "target URL should be detected as invalid")
	assert.ErrorIs(t, err2, ErrUnsupportedProtocol, "target URL should be detected as unsupported (other than http and https)")
	assert.ErrorIs(t, err3, ErrMissingPort, "target URL without port should cause an error")
}

func TestNewHTTPClientCorrectUrl(t *testing.T) {
	timeout := time.Millisecond * 200
	maxConnections := 123

	c, err := newHTTPClient("http://localhost:3000/", timeout, maxConnections, true, true)
	c2, err2 := newHTTPClient("https://localhost:443/", timeout, maxConnections, false, false)

	assert.Nil(t, err, "correct http target")
	assert.Equal(t, "localhost:3000", c.Addr)
	assert.Equal(t, maxConnections, c.MaxConns)
	assert.Equal(t, timeout, c.ReadTimeout)
	assert.Equal(t, timeout, c.WriteTimeout)
	assert.Equal(t, false, c.IsTLS)
	assert.Equal(t, true, c.TLSConfig.InsecureSkipVerify)

	assert.Nil(t, err2, "correct https target")
	assert.Equal(t, "localhost:443", c2.Addr)
	assert.Equal(t, maxConnections, c2.MaxConns)
	assert.Equal(t, timeout, c2.ReadTimeout)
	assert.Equal(t, timeout, c2.WriteTimeout)
	assert.Equal(t, true, c2.IsTLS)
	assert.Equal(t, false, c2.TLSConfig.InsecureSkipVerify)
}

func TestPerformRequest(t *testing.T) {
	timeout := time.Millisecond * 100

	c, _ := newHTTPClient("http://localhost:3000/", timeout, 10, true, false)
	r := common.RequestHeaders{}
	customHeader := common.RequestHeaders{common.RequestHeader{Key: "Custom-Header", Value: "gocannon"}}

	codeOk, _, _ := performRequest(c, "http://localhost:3000/", "GET", []byte(""), r)
	testBody := common.RawRequestBody([]byte("testbody"))
	codePost, _, _ := performRequest(c, "http://localhost:3000/postonly", "POST", testBody, r)
	codeISE, _, _ := performRequest(c, "http://localhost:3000/error", "GET", []byte(""), r)
	codeTimeout, _, _ := performRequest(c, "http://localhost:3000/timeout", "GET", []byte(""), r)
	codeMissingHeader, _, _ := performRequest(c, "http://localhost:3000/customheader", "GET", []byte(""), r)
	codeWithHeader, _, _ := performRequest(c, "http://localhost:3000/customheader", "GET", []byte(""), customHeader)

	assert.Equal(t, 200, codeOk)
	assert.Equal(t, 200, codePost)
	assert.Equal(t, http.StatusInternalServerError, codeISE)
	assert.Equal(t, 0, codeTimeout)
	assert.Equal(t, http.StatusBadRequest, codeMissingHeader)
	assert.Equal(t, 200, codeWithHeader)
}

func TestPerformRequestHTTPS(t *testing.T) {
	timeout := time.Second * 3

	c, _ := newHTTPClient("https://dev.kuffel.io:443/", timeout, 1, false, true)
	r := common.RequestHeaders{}

	codeOk, _, _ := performRequest(c, "https://dev.kuffel.io:443/", "GET", []byte(""), r)

	assert.Equal(t, 200, codeOk)
}

func TestPerformRequestHTTPSInvalidCert(t *testing.T) {
	timeout := time.Second * 3
	targetBadCert := "https://self-signed.badssl.com:443/"

	trustingClient, _ := newHTTPClient(targetBadCert, timeout, 1, true, false)
	regularClient, _ := newHTTPClient(targetBadCert, timeout, 1, false, false)
	r := common.RequestHeaders{}

	codeTrusting, _, _ := performRequest(trustingClient, targetBadCert, "GET", []byte(""), r)
	codeRegular, _, _ := performRequest(regularClient, targetBadCert, "GET", []byte(""), r)

	assert.Equal(t, 200, codeTrusting)
	assert.Equal(t, 0, codeRegular)
}
