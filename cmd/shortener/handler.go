package main

import (
	"github.com/go-chi/chi"
	"io"
	"net/http"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		shortURL := chi.URLParam(r, "short")
		println("shorturl", shortURL)
		w.Header().Add("content-type", "text/plain")
		w.Header().Add("Location", shortURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(""))
		return
	}
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		println(string(body), err)
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("test-short-url"))
		return
	}
	http.Error(w, "now allowed method", http.StatusMethodNotAllowed)
}
