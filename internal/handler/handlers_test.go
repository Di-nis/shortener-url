package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Di-nis/shortener-url/internal/compress"
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
	handler := compress.GzipMiddleware(router)

	server := httptest.NewServer(handler)

	defer server.Close()

	type want struct {
		statusCode      int
		response        string
		contentType     string
		contentEncoding string
	}

	testsCreate := []struct {
		name            string
		body            string
		method          string
		contentType     string
		contentEncoding string
		acceptEncoding  string
		want            want
	}{
		{
			name:            "TestCreateURLFromText, метод - POST, короткий URL сформирован",
			body:            "https://practicum.yandex.ru",
			method:          http.MethodPost,
			contentType:     "text/plain",
			contentEncoding: "",
			acceptEncoding:  "",
			want: want{
				statusCode:  http.StatusCreated,
				response:    "http://localhost:8080/bTKNZu94",
				contentType: "text/plain",
			},
		},
		{
			name:        "TestCreateURLFromText, метод - GET, метод не соответствует требованиям",
			body:        "https://practicum.yandex.ru",
			method:      http.MethodGet,
			contentType: "text/plain",
			contentEncoding: "",
			acceptEncoding:  "",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				response:    "",
				contentType: "",
			},
		},
		{
			name:        "TestCreateURLFromText, метод - POST, запрос не содержит url",
			body:        "",
			method:      http.MethodPost,
			contentType: "",
			contentEncoding: "",
			acceptEncoding:  "",
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
			req.SetHeaders(map[string]string{
				"Content-Type":     tt.contentType,
				"Content-Encoding": tt.contentEncoding,
				"Accept-Encoding":  tt.acceptEncoding,
			})

			resp, err := req.Send()

			require.NoError(t, err, "error making HTTP request", tt.body)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			assert.Equal(t, tt.want.response, string(resp.Body()))
			assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
		})
	}

	testsCreateJSON := []struct {
		name            string
		body            string
		method          string
		contentType     string
		contentEncoding string
		acceptEncoding  string
		want            want
	}{
		{
			name:        "TestCreateURLFromJSON, метод - GET, метод не соответствует требованиям",
			body:        `{"url": "https://www.sports.ru"}`,
			method:      http.MethodGet,
			contentType: "application/json",
			contentEncoding: "",
			acceptEncoding:  "",
			want: want{
				statusCode:      http.StatusMethodNotAllowed,
				response:        "",
				contentType:     "",
				contentEncoding: "",
			},
		},
		{
			name:        "TestCreateURLFromJSON, метод - POST, тело запроса - пустое",
			body:        "",
			method:      http.MethodPost,
			contentType: "application/json",
			contentEncoding: "",
			acceptEncoding:  "",
			want: want{
				statusCode:      http.StatusBadRequest,
				response:        "Не удалось прочитать тело запроса\n",
				contentType:     "text/plain; charset=utf-8",
				contentEncoding: "",
			},
		},
		{
			name:        "TestCreateURLFromJSON, метод - POST, данные не валидны",
			body:        `{"url": 555}`,
			method:      http.MethodPost,
			contentType: "application/json",
			contentEncoding: "",
			acceptEncoding:  "",
			want: want{
				statusCode:      http.StatusBadRequest,
				response:        "json: cannot unmarshal number into Go struct field Request.url of type string\n",
				contentType:     "text/plain; charset=utf-8",
				contentEncoding: "",
			},
		},
		{
			name:            "TestCreateURLFromJSON, метод - POST, короткий URL сформирован",
			body:            `{"url": "https://www.sports.ru"}`,
			method:          http.MethodPost,
			contentType:     "application/json",
			contentEncoding: "",
			acceptEncoding:  "",
			want: want{
				statusCode:      http.StatusCreated,
				response:        `{"result":"http://localhost:8080/4BeKySvE"}`,
				contentType:     "application/json",
				contentEncoding: "gzip",
			},
		},
	}

	for _, tt := range testsCreateJSON {
		t.Run(tt.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.URL = server.URL + "/api/shorten"
			req.Header.Set("Content-Type", tt.contentType)
			req.Header.Set("Content-Encoding", tt.contentEncoding)
			req.Header.Set("Accept-Encoding", tt.acceptEncoding)

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
			name:     "TestGetURL, метод - Get, адрес - существующий в БД адрес, кейс 1",
			shortURL: "bTKNZu94",
			method:   http.MethodGet,
			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				response:    "https://practicum.yandex.ru",
				contentType: "text/plain",
			},
		},
		{
			name:     "TestGetURL, метод - Get, адрес - существующий в БД адрес, кейс 2",
			shortURL: "4BeKySvE",
			method:   http.MethodGet,
			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				response:    "https://www.sports.ru",
				contentType: "text/plain",
			},
		},
		{
			name:     "TestGetURL, метод - Post, адрес - существующий в БД адрес",
			shortURL: "bTKNZu94",
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
