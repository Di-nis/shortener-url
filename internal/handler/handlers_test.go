package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createShortURL(t *testing.T) {
	router := CreateRouter()
	server := httptest.NewServer(router)

	defer server.Close()

	type want struct {
		statusCode  int
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
				statusCode:  http.StatusCreated,
				response:    `http://localhost:8080/EwHXdJfB`,
				contentType: "text/plain",
			},
		},
		{
			name:   "Test_createShortURL, запрос 2",
			body:   `https://practicum.yandex.ru/`,
			method: http.MethodGet,
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				response:    ``,
				contentType: "",
			},
		},
		{
			name:   "Test_createShortURL, запрос 3",
			body:   ``,
			method: http.MethodPost,
			want: want{
				statusCode:  http.StatusBadRequest,
				response:    "Не удалось прочитать тело запроса\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.URL = server.URL
			req.Body = tt.body

			resp, err := req.Send()

			require.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			assert.Equal(t, tt.want.response, string(resp.Body()))
			assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
		})
	}
}

func Test_getOriginalURL(t *testing.T) {
	router := CreateRouter()
	server := httptest.NewServer(router)

	defer server.Close()

	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name     string
		shortUrl string
		method   string
		want     want
	}{
		{
			name:     "Test_getOriginalURL, метод - Get, адрес - существующий в БД адрес",
			shortUrl: "EwHXdJfB",
			method:   http.MethodGet,
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    `https://practicum.yandex.ru/`,
				contentType: "text/plain",
			},
		},
		{
			name:     "Test_getOriginalURL, метод - Post, адрес - существующий в БД адрес",
			shortUrl: "EwHXdJfB",
			method:   http.MethodPost,
			want: want{
				code:        http.StatusMethodNotAllowed,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "Test_getOriginalURL, метод - Post, адрес в БД не найден",
			shortUrl: "nvjkrhsfdvn",
			method:   http.MethodGet,
			want: want{
				code:        http.StatusNotFound,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).R()
			req.Method = tt.method
			req.URL = server.URL + "/" + tt.shortUrl

			resp, err := req.Send()
			if err != nil {
				assert.True(t, strings.Contains(err.Error(), "auto redirect is disabled"))
			}

			assert.Equal(t, tt.want.code, resp.StatusCode())
			assert.Equal(t, tt.want.response, string(resp.Header().Get("Location")))
			assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))

		})
	}
}
