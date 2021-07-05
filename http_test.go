package main

import (
	"fmt"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Oooops...")
	})

	http.HandleFunc("/timeout", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 200)
		fmt.Fprintf(w, "Too late...")
	})

	go func() {
		http.ListenAndServe(":3000", nil)
	}()

	c, _ := newHTTPClient("http://localhost:3000/", timeout, 10)

	codeOk, _, _ := performRequest(c, "http://localhost:3000/")
	codeISE, _, _ := performRequest(c, "http://localhost:3000/error")
	codeTimeout, _, _ := performRequest(c, "http://localhost:3000/timeout")

	assert.Equal(t, 200, codeOk)
	assert.Equal(t, http.StatusInternalServerError, codeISE)
	assert.Equal(t, -1, codeTimeout)
}
