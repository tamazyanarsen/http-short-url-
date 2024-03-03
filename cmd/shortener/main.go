package main

import (
	"io"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			shortUrl := strings.Split(r.URL.Path, "/")[1]
			println(shortUrl)

			w.Header().Add("content-type", "text/plain")
			w.Header().Add("Location", "sdfsdfdsf")
			w.WriteHeader(301)
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
		return
	})
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
