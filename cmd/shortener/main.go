package main

import (
	"http-short-url/cmd/shortener/config"
	"http-short-url/cmd/shortener/handler"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	config.InitConfig()
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
