package handler

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"http-short-url/cmd/shortener/config"
	"http-short-url/cmd/shortener/data"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

var urlStore data.Store

var sugarLogger zap.SugaredLogger

// var logger, err = zap.NewDevelopment()

func init() {
	if logger, err := zap.NewDevelopment(); err == nil {
		sugarLogger = *logger.Sugar()
	}
	urlStore = new(data.URLStore)
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

type gzipWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	sugarLogger.Infoln("new gzipWriter")
	return &gzipWriter{
		ResponseWriter: w,
		Writer:         gzip.NewWriter(w),
	}
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	sugarLogger.Infoln("response Content-Type", w.Header().Get("Content-Type"))
	if !(strings.Contains(w.Header().Get("Content-Type"), "application/json") ||
		strings.Contains(w.Header().Get("Content-Type"), "text/html")) {
		sugarLogger.Infoln("call ResponseWriter:", string(b))
		return w.ResponseWriter.Write(b)
	}
	sugarLogger.Infoln("call gzip write:", string(b))
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Accept-Encoding", "gzip")
	return w.Writer.Write(b)
}

func GzipHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newHandler := w
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			sugarLogger.Infoln("Accept-Encoding NOT contains gzip:", r.Header.Get("Accept-Encoding"))
			newHandler = newGzipWriter(w)
			defer newHandler.(*gzipWriter).Writer.Close()
			//next(w, r)
			//return
		}

		sugarLogger.Infoln("Accept-Encoding contains gzip:", r.Header.Get("Accept-Encoding"))

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			sugarLogger.Infoln("request Content-Encoding", r.Header.Get("Content-Encoding"))
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			r.Body = gz
		}

		next(newHandler, r)
		//next(w, r)
	})
}

func GetShort(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "short")
	// println("shorturl", shortURL, len(urls), urls[shortURL])
	url, ok := urlStore.Read(regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(shortURL, ""))
	if ok {
		w.Header().Add("content-type", "text/plain")
		w.Header().Add("Location", url)
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
	shortURL, addr := shortName(body)
	println("shortURL:", shortURL, "addr:", addr)
	w.Write([]byte(addr + shortURL))
}

func shortName(originalURL []byte) (string, string) {
	short := base64.StdEncoding.EncodeToString(originalURL)[:8]
	shortURL := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(short, "")
	urlStore.Write(shortURL, string(originalURL))
	addr := *config.Config["b"]
	if addr[len(addr)-1:] != "/" {
		addr += "/"
	}
	return shortURL, addr
}

func PostJSON(w http.ResponseWriter, r *http.Request) {
	sugarLogger.Info("call PostJSON")
	var body struct {
		URL string `json:"url"`
	}
	var resp struct {
		Result string `json:"result"`
	}
	if reqBody, err := io.ReadAll(r.Body); err == nil {
		sugarLogger.Infoln("reqBody:", string(reqBody))
		if err := json.Unmarshal(reqBody, &body); err == nil {
			sugarLogger.Infoln("body.URL", body.URL)
			shortURL, addr := shortName([]byte(body.URL))
			resp.Result = addr + shortURL
			if response, err := json.Marshal(resp); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				write, err := w.Write(response)
				if err != nil {
					println(err.Error())
				} else {
					println("write:", write)
				}
			} else {
				println(err.Error())
			}
		} else {
			println(err.Error())
		}
	} else {
		println("ошибка парсинга body, ReadAll", err.Error())
	}
}
