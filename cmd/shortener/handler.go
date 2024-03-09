package main

import (
	"io"
	"net/http"
	"strings"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		shortURL := strings.Split(r.URL.Path, "/")[1]
		w.Header().Add("content-type", "text/plain")
		w.Header().Add("Location", shortURL)
		w.WriteHeader(307)
		return
	}
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		println(string(body), err)
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(201)
		w.Write([]byte("test-short-url"))
		return
	}
	http.Error(w, "not allowed method", http.StatusMethodNotAllowed)
}
