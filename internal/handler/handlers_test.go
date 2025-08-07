package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bytes"
	"io"
)

func Test_createShortURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		body   string
		method string
		want   want
	}{
		{
			name:   "Test_createShortURL, запрос 1",
			body:   `https://practicum.yandex.ru/`,
			method: http.MethodPost,
			want: want{
				code:        http.StatusCreated,
				response:    `http://localhost:8080/EwHXdJfB`,
				contentType: "text/plain",
			},
		},
		{
			name:   "Test_createShortURL, запрос 2",
			body:   `https://practicum.yandex.ru/`,
			method: http.MethodGet,
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    ``,
				contentType: "",
			},
		},
		{
			name:   "Test_createShortURL, запрос 3",
			body:   ``,
			method: http.MethodPost,
			want: want{
				code:        http.StatusBadRequest,
				response:    "Не удалось прочитать тело запроса\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(tt.body)
			request := httptest.NewRequest(tt.method, "/", bytes.NewBuffer(body))

			w := httptest.NewRecorder()
			createShortURL(w, request)

			res := w.Result()
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_getOriginalURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		url  string
		method string
		want want
	}{
		{
			name: "Test_getOriginalURL, метод - Get, существующий в БД адрес",
			url:  "EwHXdJfB",
			method: http.MethodGet,
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    `https://practicum.yandex.ru/`,
				contentType: "text/plain",
			},
		},
		{
			name: "Test_getOriginalURL, метод - Post, существующий в БД адрес",
			url:  "EwHXdJfB",
			method: http.MethodPost,
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "",
			},
		},
		{
			name: "Test_getOriginalURL, метод - Post, адрес в БД не найден",
			url:  "nvjkrhsfdvn",
			method: http.MethodGet,
			want: want{
				code:        http.StatusNotFound,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/" + tt.url
			request := httptest.NewRequest(tt.method, url, nil)

			w := httptest.NewRecorder()
			getOriginalURL(w, request)

			res := w.Result()
			headerLocation := res.Header.Get("Location")
			defer res.Body.Close()

			assert.Equal(t, tt.want.response, string(headerLocation))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.code, res.StatusCode)

		})
	}
}
