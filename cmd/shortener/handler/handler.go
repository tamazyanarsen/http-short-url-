package handler

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"http-short-url/cmd/shortener/config"
	"http-short-url/cmd/shortener/data"
	file_handler "http-short-url/cmd/shortener/fileHandler"
	"http-short-url/cmd/shortener/logger"
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

func readFile(cons *file_handler.Consumer) {
	sugarLogger.Infoln("START READ FILE")
	fileData, fileErr := cons.ReadEvent()
	if fileErr != nil {
		sugarLogger.Infoln(fileErr, "КОНЕЦ ФАЙЛА")
		return
	}
	sugarLogger.Infoln("WRITE TO STORE", urlStore)
	urlStore.Write(fileData.ShortURL, fileData.OriginalURL)
	readFile(cons)
}

func InitHandler() error {
	logger.InitLogger()
	sugarLogger = logger.Logger
	sugarLogger.Infoln("INIT STORE", urlStore)

	if *config.Config["f"] == "" {
		urlStore = new(data.URLStore)
	} else {
		cons, consErr := file_handler.NewConsumer(*config.Config["f"])
		if consErr != nil {
			sugarLogger.Errorln(consErr.Error())
			return consErr
		}
		urlStore = new(data.FileStore)
		sugarLogger.Infoln("call readFile()")
		readFile(cons)
	}
	return nil
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

func WithLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseHandler := responseInfo{
			ResponseWriter: w,
		}
		sugarLogger.Infoln("request url", r.URL.Path)
		sugarLogger.Infoln("request method", r.Method)

		startRequestTime := time.Now()
		handler.ServeHTTP(&responseHandler, r)
		sugarLogger.Infoln("duration:", time.Since(startRequestTime))
		sugarLogger.Infoln("\n----------------------------------------------------------\n\n")
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
	defer w.Writer.Close()
	return w.Writer.Write(b)
}

func GzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newHandler := w
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			sugarLogger.Infoln("Accept-Encoding contains gzip:", r.Header.Get("Accept-Encoding"))
			newHandler = newGzipWriter(w)
			// defer newHandler.(*gzipWriter).Writer.Close()
		}

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

		next.ServeHTTP(newHandler, r)
		//next(w, r)
	})
}

func GetShort(w http.ResponseWriter, r *http.Request) {
	sugarLogger.Info("START GetShort")
	shortURL := chi.URLParam(r, "short")
	// println("shorturl", shortURL, len(urls), urls[shortURL])
	url, err := urlStore.Read(regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(shortURL, ""))
	sugarLogger.Infoln("original url from store to header.Location", url)
	w.Header().Add("content-type", "text/plain")
	if err == nil {
		w.Header().Add("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func PostURL(w http.ResponseWriter, r *http.Request) {
	sugarLogger.Info("START PostURL")
	body, bodyErr := io.ReadAll(r.Body)
	sugarLogger.Infoln("get body", string(body))
	if bodyErr != nil {
		handleError(bodyErr, w)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	shortURL, addr, writeErr := shortName(body)
	if writeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sugarLogger.Infoln("shortURL", shortURL)
	sugarLogger.Infoln("addr", addr)
	// if writeToFile(shortURL, body) != nil {
	// 	return
	// }
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(addr + shortURL))
}

func handleError(err error, w http.ResponseWriter) {
	sugarLogger.Errorln(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func shortURL(originalURL []byte) string {
	short := base64.StdEncoding.EncodeToString(originalURL)[:]
	return regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(short, "")
}

func shortName(originalURL []byte) (string, string, error) {
	shortURL := shortURL(originalURL)
	if err := urlStore.Write(shortURL, string(originalURL)); err != nil {
		return "", "", err
	}
	addr := *config.Config["b"]
	if addr[len(addr)-1:] != "/" {
		addr += "/"
	}
	return shortURL, addr, nil
}

func PostJSON(w http.ResponseWriter, r *http.Request) {
	sugarLogger.Info("START PostJSON")
	var body struct {
		URL string `json:"url"`
	}
	var resp struct {
		Result string `json:"result"`
	}

	reqBody, bodyErr := io.ReadAll(r.Body)
	if bodyErr != nil {
		handleError(bodyErr, w)
		return
	}
	sugarLogger.Infoln("reqBody:", string(reqBody))

	if jsonErr := json.Unmarshal(reqBody, &body); jsonErr != nil {
		handleError(jsonErr, w)
		return
	}
	sugarLogger.Infoln("body.URL", body.URL)
	shortURL, addr, writeErrStore := shortName([]byte(body.URL))
	if writeErrStore != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.Result = addr + shortURL

	response, respErr := json.Marshal(resp)
	if respErr != nil {
		handleError(respErr, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
	}
	// if writeToFile(shortURL, []byte(body.URL)) != nil {
	// 	return
	// }
	w.WriteHeader(http.StatusCreated)

	_, writeErr := w.Write([]byte(strings.TrimSuffix(string(response), "\n")))
	if writeErr != nil {
		sugarLogger.Errorln(writeErr.Error())
	}
}
