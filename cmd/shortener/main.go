package main

import (
	"flag"
	"http-short-url/internal/config"
	"http-short-url/internal/handler"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	config.InitConfig()
	println("---------------CONFIG before flag.parse-------------------", *config.Config["a"], *config.Config["b"], *config.Config["f"])
	flag.Parse()
	println("---------------CONFIG after flag.parse-------------------", *config.Config["a"], *config.Config["b"], *config.Config["f"])
	if err := handler.InitHandler(); err != nil {
		log.Fatal(err)
	}
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
