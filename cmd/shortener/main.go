package main

import (
	"flag"
	"http-short-url/cmd/shortener/app"
	"http-short-url/cmd/shortener/config"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	flag.Parse()
	println(*config.Config["a"], *config.Config["b"])
	r := chi.NewRouter()
	r.Get("/{short}", app.WithLog(app.GetShort))
	// r.Get("/{short}", app.GetShort)
	// r.Post("/", app.PostURL)
	r.Post("/", app.WithLog(app.PostURL))
	err := http.ListenAndServe(*config.Config["a"], r)
	if err != nil {
		log.Fatal(err)
	}
}
