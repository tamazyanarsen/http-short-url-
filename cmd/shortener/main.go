package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Get("/{someText}", MainHandler)
	r.Post("/", MainHandler)
	http.ListenAndServe(":8080", r)
}
