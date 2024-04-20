package main

import (
	"flag"
	"github.com/go-chi/chi"
	"http-short-url/cmd/shortener/config"
	"net/http"
)

func main() {
	flag.Parse()
	println(*config.Config["a"], *config.Config["b"])
	r := chi.NewRouter()
	r.Get("/{short}", MainHandler)
	r.Post("/", MainHandler)
	err := http.ListenAndServe(*config.Config["a"], r)
	if err != nil {
		return
	}
}
