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
	config.InitConfig()
	println("---------------CONFIG-------------------", *config.Config["a"], *config.Config["b"], *config.Config["f"])
	flag.Parse()
	handler.InitHandler()
	r := chi.NewRouter()
	r.Use(handler.GzipHandler)
	r.Use(handler.WithLog)
	r.Get("/{short}", handler.GetShort)
	r.Post("/", handler.PostURL)
	r.Post("/api/shorten", handler.PostJSON)
	err := http.ListenAndServe(*config.Config["a"], r)
	if err != nil {
		log.Fatal(err)
	}
}
