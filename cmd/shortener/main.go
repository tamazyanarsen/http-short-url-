package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Get("/{short}", MainHandler)
	r.Post("/", MainHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
