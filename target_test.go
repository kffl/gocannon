package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
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

	exitVal := m.Run()
	os.Exit(exitVal)
}
