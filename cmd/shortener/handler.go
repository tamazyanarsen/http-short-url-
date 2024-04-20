package main

import (
	"encoding/base64"
	"github.com/go-chi/chi"
	"http-short-url/cmd/shortener/config"
	"io"
	"net/http"
	"regexp"
)

var urls = make(map[string]string)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		shortURL := chi.URLParam(r, "short")
		println("shorturl", shortURL, len(urls), urls[shortURL])
		w.Header().Add("content-type", "text/plain")
		w.Header().Add("Location",
			urls[regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(shortURL, "")])
		w.WriteHeader(http.StatusTemporaryRedirect)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		println(string(body), err)
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		short := base64.StdEncoding.EncodeToString(body)[:8]
		shortURL := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(short, "")
		urls[shortURL] = string(body)
		println("save", shortURL, "to map; result:", urls[shortURL])
		addr := *config.Config["b"]
		if addr[len(addr)-1:] != "/" {
			addr += "/"
		}
		w.Write([]byte(addr + shortURL))
	default:
		http.Error(w, "not allowed method", http.StatusMethodNotAllowed)
	}
}
