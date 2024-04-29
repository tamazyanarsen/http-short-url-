package app

import (
	"encoding/base64"
	"http-short-url/cmd/shortener/config"
	"io"
	"net/http"
	"regexp"
	"sync"

	"github.com/go-chi/chi"
)

// var urls = make(map[string]string)

var urlStore sync.Map

func GetShort(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "short")
	// println("shorturl", shortURL, len(urls), urls[shortURL])
	url, ok := urlStore.Load(regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(shortURL, ""))
	if ok {
		w.Header().Add("content-type", "text/plain")
		w.Header().Add("Location", url.(string))
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func PostURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	println(string(body), err)
	w.Header().Add("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	short := base64.StdEncoding.EncodeToString(body)[:8]
	shortURL := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(short, "")
	// urls[shortURL] = string(body)
	urlStore.Store(shortURL, string(body))
	// println("save", shortURL, "to map; result:", urls[shortURL])
	addr := *config.Config["b"]
	if addr[len(addr)-1:] != "/" {
		addr += "/"
	}
	w.Write([]byte(addr + shortURL))
}
