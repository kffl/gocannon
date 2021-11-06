package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// spawn an HTTP server that will be used as target for other tests

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	http.HandleFunc("/postonly", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Wrong method")
		} else {
			fmt.Fprintf(w, "Ok")
		}
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

	exitVal := m.Run()
	os.Exit(exitVal)
}
