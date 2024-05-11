package app

import (
	"encoding/base64"
	"http-short-url/cmd/shortener/config"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// var urls = make(map[string]string)

var urlStore sync.Map

var sugarLogger zap.SugaredLogger

// var logger, err = zap.NewDevelopment()

func init() {
	if logger, err := zap.NewDevelopment(); err == nil {
		sugarLogger = *logger.Sugar()
	}
}

type responseInfo struct {
	http.ResponseWriter
}

func (r *responseInfo) Write(b []byte) (int, error) {
	sugarLogger.Infoln("response size:", len(b))
	return r.ResponseWriter.Write(b)
}

func (r *responseInfo) WriteHeader(s int) {
	sugarLogger.Infoln("response header:", s)
	r.ResponseWriter.WriteHeader(s)
}

func WithLog(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseHandler := responseInfo{
			ResponseWriter: w,
		}
		sugarLogger.Infoln("request url", r.URL.Path)
		sugarLogger.Infoln("request method", r.Method)

		startRequestTime := time.Now()
		handler(&responseHandler, r)
		sugarLogger.Infoln("duration:", time.Since(startRequestTime))
	})
}

func GetShort(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "short")
	// println("shorturl", shortURL, len(urls), urls[shortURL])
	url, ok := urlStore.Load(regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(shortURL, ""))
	if ok {
		w.Header().Add("content-type", "text/plain")
		w.Header().Add("Location", url.(string))
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
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
