package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/repository"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/Di-nis/shortener-url/internal/usecase"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAndGetURL(t *testing.T) {
	options := new(config.Config)
	options.Parse()

	repo := repository.NewRepo()
	serv := service.NewService()
	urlUseCase := usecase.NewURLUseCase(repo, serv)
	controller := NewСontroller(urlUseCase, options)

	router := controller.CreateRouter()

	server := httptest.NewServer(router)

	defer server.Close()

	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	testsCreate := []struct {
		name     string
		body     string
		shortURL string
		method   string
		want     want
	}{
		{
			name:   "TestCreateURL, метод - POST, короткий URL сформирован",
			body:   `https://practicum.yandex.ru/`,
			method: http.MethodPost,
			want: want{
				statusCode:  http.StatusCreated,
				response:    "http://localhost:8080/5J3xKXF9",
				contentType: "text/plain",
			},
		},
		{
			name:   "TestCreateURL, метод - GET, метод не соответствует требованиям",
			body:   "https://practicum.yandex.ru/",
			method: http.MethodGet,
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				response:    ``,
				contentType: "",
			},
		},
		{
			name:   "TestCreateURL, метод - POST, запрос не содержит url",
			body:   ``,
			method: http.MethodPost,
			want: want{
				statusCode:  http.StatusBadRequest,
				response:    "Не удалось прочитать тело запроса\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range testsCreate {
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

	testsGet := []struct {
		name     string
		shortURL string
		method   string
		want     want
	}{
		{
			name:     "TestGetURL, метод - Get, адрес - существующий в БД адрес",
			shortURL: "5J3xKXF9",
			method:   http.MethodGet,
			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				response:    `https://practicum.yandex.ru/`,
				contentType: "text/plain",
			},
		},
		{
			name:     "TestGetURL, метод - Post, адрес - существующий в БД адрес",
			shortURL: "5J3xKXF9",
			method:   http.MethodPost,
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "TestGetURL, метод - Post, адрес в БД не найден",
			shortURL: "nvjkrhsf",
			method:   http.MethodGet,
			want: want{
				statusCode:  http.StatusNotFound,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range testsGet {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).R()
			req.Method = tt.method
			req.URL = server.URL + "/" + tt.shortURL

			resp, err := req.Send()
			if err != nil {
				assert.True(t, strings.Contains(err.Error(), "auto redirect is disabled"))
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			assert.Equal(t, tt.want.response, string(resp.Header().Get("Location")))
			assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))

		})
	}
}
