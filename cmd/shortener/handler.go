package main

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		shortURL := chi.URLParam(r, "date")
		w.Header().Add("content-type", "text/plain")
		w.Header().Add("Location", shortURL)
		w.WriteHeader(307)
		w.Write([]byte(""))
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
	http.Error(w, "now allowed method", http.StatusMethodNotAllowed)
}
