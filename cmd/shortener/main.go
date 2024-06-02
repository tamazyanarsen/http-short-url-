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
	r.Get("/{short}", handler.WithLog(handler.GzipHandler(handler.GetShort)))
	r.Post("/", handler.WithLog(handler.GzipHandler(handler.PostURL)))
	r.Post("/api/shorten", handler.WithLog(handler.GzipHandler(handler.PostJSON)))
	err := http.ListenAndServe(*config.Config["a"], r)
	if err != nil {
		log.Fatal(err)
	}
}
