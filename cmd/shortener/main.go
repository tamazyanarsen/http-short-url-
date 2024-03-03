package main

import (
	"io"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}
		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			println(string(body), err)
			w.Header().Add("content-type", "text/plain")
			w.WriteHeader(201)
			w.Write([]byte("sdfsdfsdf"))
		}
	})
	mux.HandleFunc("/{url}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
			return
		}
		for k, v := range r.URL.Query() {
			println(k, v)
		}
		w.Header().Add("content-type", "text/plain")
		w.Header().Add("Location", "sdfsdfdsf")
		w.WriteHeader(301)
		w.Write([]byte(""))
	})
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
