package main

import (
	"io"
	"net/http"
	"strings"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		shortURL := strings.Split(r.URL.Path, "/")[1]
		w.Header().Add("content-type", "text/plain")
		w.Header().Add("Location", shortURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		println(string(body), err)
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("test-short-url"))
	default:
		http.Error(w, "not allowed method", http.StatusMethodNotAllowed)
	}
}
