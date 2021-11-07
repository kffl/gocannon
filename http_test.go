package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPClientWrongURL(t *testing.T) {

	_, err1 := newHTTPClient("XYZthisisawrongurl123", time.Millisecond*200, 10)
	_, err2 := newHTTPClient("https://something/", time.Millisecond*200, 10)

	assert.ErrorIs(t, err1, ErrWrongTarget, "target URL should be detected as invalid")
	assert.ErrorIs(t, err2, ErrUnsupportedProtocol, "target URL should be detected as non-http")
}

func TestNewHTTPClientCorrectUrl(t *testing.T) {
	timeout := time.Millisecond * 200
	maxConnections := 123

	c, err := newHTTPClient("http://localhost:3000/", timeout, maxConnections)

	assert.Nil(t, err, "correct target")
	assert.Equal(t, "localhost:3000", c.Addr)
	assert.Equal(t, maxConnections, c.MaxConns)
	assert.Equal(t, timeout, c.ReadTimeout)
	assert.Equal(t, timeout, c.WriteTimeout)
}

func TestPerformRequest(t *testing.T) {
	timeout := time.Millisecond * 100

	c, _ := newHTTPClient("http://localhost:3000/", timeout, 10)
	r := requestHeaders{}
	customHeader := requestHeaders{requestHeader{"Custom-Header", "gocannon"}}

	codeOk, _, _ := performRequest(c, "http://localhost:3000/", "GET", []byte(""), r)
	testBody := rawRequestBody([]byte("testbody"))
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
