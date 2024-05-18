package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type args struct {
	w            *httptest.ResponseRecorder
	method       string
	url          string
	body         io.Reader
	want         string
	expectedCode int
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name string
		args args
	}{
		{name: "post test", args: args{
			method:       http.MethodPost,
			url:          "/",
			body:         bytes.NewReader([]byte("https://practicum.yandex.ru")),
			want:         "http://localhost:8080/aHR0cHM6",
			expectedCode: http.StatusCreated,
		}},
		{name: "post json test", args: args{
			method:       http.MethodPost,
			url:          "/api/shorten",
			body:         bytes.NewReader([]byte("{\"url\": \"https://practicum.yandex.ru\"}")),
			want:         "{\"result\":\"http://localhost:8080/aHR0cHM6\"}",
			expectedCode: http.StatusCreated,
		}},
		{name: "get test", args: args{
			method:       http.MethodGet,
			url:          "/aHR0cHM6",
			body:         nil,
			want:         "",
			expectedCode: http.StatusOK,
		}},
	}
	r := chi.NewRouter()
	r.Get("/{short}", GetShort)
	r.Post("/", PostURL)
	r.Post("/api/shorten", PostJSON)

	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := resty.New().R()
			client.URL = ts.URL + tt.args.url
			client.Method = tt.args.method
			client.SetBody(tt.args.body)
			res, err := client.Send()
			if err != nil {
				t.Fatal("req err")
			}
			println(res.StatusCode(), res.Header().Get("Location"))

			data := ""
			if tt.args.method == http.MethodGet {
				data = res.Header().Get("Location")
			}
			if tt.args.method == http.MethodPost {
				data = string(res.Body())
			}
			require.Equal(t, tt.args.expectedCode, res.StatusCode(), "Неверный код ответа")
			assert.Equal(t, tt.args.want, data)
		})
	}
}
