package main

import (
	"flag"
	"http-short-url/cmd/shortener/config"
	"http-short-url/cmd/shortener/handler"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	flag.Parse()
	println(*config.Config["a"], *config.Config["b"])
	r := chi.NewRouter()
	r.Get("/{short}", handler.WithLog(handler.GetShort))
	// r.Get("/{short}", handler.GetShort)
	// r.Post("/", handler.PostURL)
	r.Post("/", handler.WithLog(handler.PostURL))
	err := http.ListenAndServe(*config.Config["a"], r)
	if err != nil {
		log.Fatal(err)
	}
}
