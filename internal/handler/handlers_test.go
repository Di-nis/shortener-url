package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/repository"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/Di-nis/shortener-url/internal/storage"
	"github.com/Di-nis/shortener-url/internal/usecase"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initHandler() (http.Handler, error) {
	cfg := config.NewConfig()
	cfg.Load()

	consumer, err := storage.NewConsumer(cfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to init consumer: %w", err)
	}

	producer, err := storage.NewProducer(cfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to init producer: %w", err)
	}

	storage := &repository.Storage{
		Consumer: consumer,
		Producer: producer,
	}

	repo := repository.NewRepoFileMemory(storage)
	svc := service.NewService()

	urlUseCase := usecase.NewURLUseCase(repo, svc)
	controller := NewСontroller(urlUseCase, cfg)
	router := controller.SetupRouter()
	return router, nil
}

func setEnv() {
	var err error

	err = os.Setenv("SERVER_ADDRESS", "localhost:8080")
	if err != nil {
		log.Fatalf("set env SERVER_ADDRESS failed: %v", err)
	}

	err = os.Setenv("BASE_URL", "http://localhost:8080")
	if err != nil {
		log.Fatalf("set env BASE_URL failed: %v", err)
	}

	err = os.Setenv("FILE_STORAGE_PATH", "../../database_test.log")
	if err != nil {
		log.Fatalf("set env FILE_STORAGE_PATH failed: %v", err)
	}
}

var testServer *httptest.Server

func TestMain(m *testing.M) {
	setEnv()

	handler, err := initHandler()
	if err != nil {
		log.Fatalf("init handler failed: %v", err)
	}
	testServer = httptest.NewServer(handler)
	defer testServer.Close()

	os.Exit(m.Run())
}

func TestController_CreateURLShortJSONBatch(t *testing.T) {
	type want struct {
		statusCode  int
		body        string
		contentType string
	}

	tests := []struct {
		name        string
		body        string
		method      string
		contentType string
		want        want
	}{
		{
			name:        "POST, тест 1",
			body:        `[{"correlation_id": "1","original_url":"sberbank.ru"},{"correlation_id":"2","original_url":"dzen.ru"}]`,
			method:      http.MethodPost,
			contentType: "text/plain",
			want: want{
				statusCode:  http.StatusCreated,
				body:        `[{"short_url":"http://localhost:8080/ghDg2efU","correlation_id":"1"},{"short_url":"http://localhost:8080/j0z83CVB","correlation_id":"2"}]`,
				contentType: "application/json",
			},
		},
		{
			name:        "GET, тест 2",
			body:        `[{"correlation_id":"1","original_url":""sport-express.ru""}]`,
			method:      http.MethodGet,
			contentType: "text/plain",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				body:        "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.URL = testServer.URL + "/api/shorten/batch"
			req.Body = tt.body
			req.SetHeaders(map[string]string{
				"Content-Type": tt.contentType,
			})

			resp, err := req.Send()
			require.NoError(t, err, "error making HTTP request", tt.body)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, string(resp.Body()))
			}
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
			}
		})
	}
}

func TestController_CreateURLFromText(t *testing.T) {
	type want struct {
		statusCode  int
		body        string
		contentType string
	}

	tests := []struct {
		name            string
		body            string
		method          string
		contentType     string
		contentEncoding string
		acceptEncoding  string
		want            want
	}{
		{
			name:            "POST, короткий URL сформирован",
			body:            "https://practicum.yandex.ru",
			method:          http.MethodPost,
			contentType:     "text/plain",
			contentEncoding: "",
			acceptEncoding:  "",
			want: want{
				statusCode:  http.StatusCreated,
				body:        "http://localhost:8080/bTKNZu94",
				contentType: "text/plain",
			},
		},
		{
			name:            "GET, метод не соответствует требованиям",
			body:            "https://practicum.yandex.ru",
			method:          http.MethodGet,
			contentType:     "text/plain",
			contentEncoding: "",
			acceptEncoding:  "",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name:            "POST, запрос не содержит url",
			body:            "",
			method:          http.MethodPost,
			contentType:     "text/plain",
			contentEncoding: "",
			acceptEncoding:  "",
			want: want{
				statusCode:  http.StatusBadRequest,
				body:        fmt.Sprintf("%s%s", constants.EmptyBodyError, "\n"),
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.URL = testServer.URL
			req.Body = tt.body
			req.SetHeaders(map[string]string{
				"Content-Type":     tt.contentType,
				"Content-Encoding": tt.contentEncoding,
				"Accept-Encoding":  tt.acceptEncoding,
			})

			resp, err := req.Send()
			require.NoError(t, err, "error making HTTP request", tt.body)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, string(resp.Body()))
			}
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
			}
		})
	}
}

func TestController_CreateURLFromJSON(t *testing.T) {
	type want struct {
		statusCode      int
		body            string
		contentType     string
		contentEncoding string
	}

	tests := []struct {
		name            string
		body            string
		method          string
		contentType     string
		contentEncoding string
		acceptEncoding  string
		want            want
	}{
		{
			name:        "GET, метод не соответствует требованиям",
			body:        `{"https://www.sports.ru"}`,
			method:      http.MethodGet,
			contentType: "application/json",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name:        "POST, тело запроса - пустое",
			body:        "",
			method:      http.MethodPost,
			contentType: "application/json",
			want: want{
				statusCode:  http.StatusBadRequest,
				body:        constants.EmptyBodyError,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "POST, данные не валидны",
			body:        `{"url": 555}`,
			method:      http.MethodPost,
			contentType: "application/json",
			want: want{
				statusCode:  http.StatusBadRequest,
				body:        constants.InvalidJSONError,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "POST, короткий URL сформирован",
			body:        `{"url": "https://www.sports.ru"}`,
			method:      http.MethodPost,
			contentType: "application/json",
			want: want{
				statusCode:      http.StatusCreated,
				body:            `{"result":"http://localhost:8080/4BeKySvE"}`,
				contentType:     "application/json",
				contentEncoding: "gzip",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.URL = testServer.URL + "/api/shorten"
			req.Header.Set("Content-Type",

				tt.contentType)
			req.Header.Set("Content-Encoding", tt.contentEncoding)
			req.Header.Set("Accept-Encoding", tt.acceptEncoding)
			req.Body = tt.body

			resp, err := req.Send()
			require.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.want.statusCode, resp.StatusCode())

			if tt.want.body != "" {
				checkBody := strings.Contains(string(resp.Body()), tt.want.body)
				assert.True(t, checkBody)
			}

			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
			}
		})
	}
}

func TestController_GetURL(t *testing.T) {
	var cookies []*http.Cookie

	t.Run("Предварительное создание данных", func(t *testing.T) {
		reqPre := resty.New().R()
		reqPre.Method = http.MethodPost
		reqPre.Body = "https://www.sports.ru"
		reqPre.URL = testServer.URL

		respPre, err := reqPre.Send()
		if err != nil {
			assert.True(t, strings.Contains(err.Error(), "auto redirect is disabled"))
		}
		cookies = respPre.Cookies()
	})

	type want struct {
		statusCode  int
		body        string
		contentType string
	}

	tests := []struct {
		name     string
		shortURL string
		method   string
		cookies  []*http.Cookie
		want     want
	}{
		{
			name:     "GET, адрес - существующий в БД адрес, кейс 1",
			shortURL: "4BeKySvE",
			method:   http.MethodGet,
			cookies:  cookies,
			want: want{
				statusCode:  http.StatusOK,
				body:        "https://www.sports.ru",
				contentType: "text/html; charset=UTF-8",
			},
		},
		{
			name:     "POST, адрес - существующий в БД адрес",
			shortURL: "bTKNZu94",
			method:   http.MethodPost,
			cookies:  cookies,
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name:     "GET, адрес в БД не найден",
			shortURL: "nvjkrhsf",
			method:   http.MethodGet,
			cookies:  cookies,
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.Cookies = tt.cookies
			req.URL = testServer.URL + "/" + tt.shortURL

			resp, err := req.Send()
			if err != nil {
				assert.True(t, strings.Contains(err.Error(), "auto redirect is disabled"))
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			if tt.want.body != "" {
				assert.Contains(t, string(resp.Body()), tt.want.body)
			}
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
			}
		})
	}
}

func TestController_GetAllURLs(t *testing.T) {
	var cookies []*http.Cookie

	t.Run("Предварительное создание данных", func(t *testing.T) {
		req := resty.New().R()
		req.Method = http.MethodPost
		req.Body = "google.ru"
		req.URL = testServer.URL

		resp, err := req.Send()
		if err != nil {
			assert.True(t, strings.Contains(err.Error(), "auto redirect is disabled"))
		}
		cookies = resp.Cookies()
	})

	type want struct {
		statusCode  int
		body        string
		contentType string
	}

	tests := []struct {
		name    string
		method  string
		cookies []*http.Cookie
		want    want
	}{
		{
			name:    "testGetAllURLs, кейс 1",
			method:  http.MethodGet,
			cookies: cookies,
			want: want{
				statusCode:  http.StatusOK,
				body:        `[{"short_url":"http://localhost:8080/5S4OlfVc","original_url":"google.ru"}]`,
				contentType: "application/json",
			},
		},
		{
			name:    "testGetAllURLs, кейс 2",
			method:  http.MethodGet,
			cookies: []*http.Cookie{},
			want: want{
				statusCode:  http.StatusNoContent,
				body:        "",
				contentType: "",
			},
		},
		{
			name:    "testGetAllURLs, метод не соответствует требованиям хендлера",
			method:  http.MethodPost,
			cookies: cookies,
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				body:        "",
				contentType: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.URL = testServer.URL + "/api/user/urls"
			req.Cookies = tt.cookies

			resp, err := req.Send()
			if err != nil {
				assert.True(t, strings.Contains(err.Error(), "auto redirect is disabled"))
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			if tt.want.body != "" {
				assert.Equal(t, tt.want.body, string(resp.Body()))
			}
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
			}
		})
	}
}

func TestController_deleteURLs(t *testing.T) {
	var cookies []*http.Cookie

	t.Run("Предварительное создание данных", func(t *testing.T) {
		req := resty.New().R()
		req.Method = http.MethodPost
		req.Body = `{"correlation_id":"1","original_url":"https://maximum.ru/"},{"correlation_id":"2","original_url":"https://radioultra.ru/"}`
		req.URL = fmt.Sprintf("%s/api/shorten/batch", testServer.URL)

		resp, err := req.Send()
		if err != nil {
			assert.True(t, strings.Contains(err.Error(), "auto redirect is disabled"))
		}
		cookies = resp.Cookies()
	})

	type want struct {
		statusCode int
	}

	tests := []struct {
		name    string
		method  string
		body    string
		cookies []*http.Cookie
		want    want
	}{
		{
			name:    "testdeleteURLs, кейс 1",
			method:  http.MethodPost,
			body:    `["imGf5jQO","gVxAI0xB"]`,
			cookies: cookies,
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name:    "testdeleteURLs, кейс 2",
			method:  http.MethodDelete,
			body:    `["imGf5jQO","gVxAI0xB"]`,
			cookies: cookies,
			want: want{
				statusCode: http.StatusAccepted,
			},
		},
	}
	for _, tt := range tests {
		req := resty.New().R()
		req.Method = tt.method
		req.Body = tt.body
		req.URL = testServer.URL + "/api/user/urls"
		req.Cookies = tt.cookies

		resp, err := req.Send()
		if err != nil {
			assert.True(t, strings.Contains(err.Error(), "auto redirect is disabled"))
		}

		if tt.want.statusCode != 0 {
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
		}
	}
}

func TestController_PingDB(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name   string
		method string
		want   want
	}{
		{
			name:   "пинг базы данных",
			method: http.MethodGet,
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
	}
	for _, tt := range tests {
		req := resty.New().R()
		req.Method = tt.method
		req.URL = testServer.URL + "/ping"

		resp, err := req.Send()
		if err == nil {
			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
		}
	}
}
