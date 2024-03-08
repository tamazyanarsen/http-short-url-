package main

import (
	"io"
	"net/http"
	"strings"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		shortURL := strings.Split(r.URL.Path, "/")[1]
		println(shortURL)

		w.Header().Add("content-type", "text/plain")
		w.Header().Add("Location", "sdfsdfdsf")
		w.WriteHeader(307)
		w.Write([]byte(""))
		return
	}
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		println(string(body), err)
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(201)
		w.Write([]byte("sdfsdfsdf"))
		return
	}
	http.Error(w, "now allowed method", http.StatusMethodNotAllowed)
}
