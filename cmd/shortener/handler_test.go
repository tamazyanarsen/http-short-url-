package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainHandler(t *testing.T) {
	type args struct {
		//w http.ResponseWriter
		w            *httptest.ResponseRecorder
		r            *http.Request
		want         string
		expectedCode int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "get test", args: args{
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodGet, "/aHR0cHM6", nil),
			want:         "https://practicum.yandex.ru/",
			expectedCode: http.StatusTemporaryRedirect,
		}},
		{name: "post test", args: args{
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("https://practicum.yandex.ru/"))),
			want:         "http://localhost:8080/aHR0cHM6",
			expectedCode: http.StatusCreated,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MainHandler(tt.args.w, tt.args.r)
			data := ""
			if tt.args.r.Method == http.MethodGet {
				data = tt.args.w.Header().Get("Location")
			}
			if tt.args.r.Method == http.MethodPost {
				data = tt.args.w.Body.String()
			}
			require.Equal(t, tt.args.expectedCode, tt.args.w.Code, "Неверный код ответа")
			assert.Equal(t, tt.args.want, data)
		})
	}
}
