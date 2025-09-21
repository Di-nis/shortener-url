package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	
	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/repository"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/Di-nis/shortener-url/internal/usecase"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// const testJWTToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTg0OTUwMjMsIlVzZXJJRCI6IjAxSzVQWDdOWjgwTldNTkVIUFJaMjNTVzJQIn0.ifklfGg01ZFVrrDLVoSN3cMc0fCoAMZzC6aTUvUhj04"

func initHandler() (http.Handler, error) {
	envTest, _ := godotenv.Read("../../.env.test")

	fileStoragePath := envTest["FILE_STORAGE_PATH"]
	serverAddress := envTest["SERVER_ADDRESS"]
	baseURL := envTest["BASE_URL"]

	cfg := &config.Config{
		ServerAddress:   serverAddress,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
	}

	consumer, err := repository.NewConsumer(cfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to init consumer: %w", err)
	}

	producer, err := repository.NewProducer(cfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to init producer: %w", err)
	}

	storage := &repository.Storage{
		Consumer: consumer,
		Producer: producer,
	}

	repo := repository.NewRepoFile(cfg.FileStoragePath, storage)
	svc := service.NewService()

	urlUseCase := usecase.NewURLUseCase(repo, svc)
	controller := NewСontroller(urlUseCase, cfg)
	router := controller.CreateRouter()
	return router, nil
}

func TestCreateAndGetURL(t *testing.T) {
	handler, _ := initHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	t.Run("TestCreateURLFromText", func(t *testing.T) {
		testCreateURLFromText(t, server)
	})

	t.Run("TestCreateURLFromJSON", func(t *testing.T) {
		testCreateURLFromJSON(t, server)
	})

	t.Run("TestGetURL", func(t *testing.T) {
		testGetURL(t, server)
	})

	t.Run("TestGetAllURLs", func(t *testing.T) {
		testGetAllURLs(t, server)
	})
}

func testCreateURLFromText(t *testing.T, server *httptest.Server) {
	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name            string
		body            string
		method          string
		contentType     string
		contentEncoding string
		acceptEncoding  string
		// authorization   string
		want            want
	}{
		{
			name:            "POST, короткий URL сформирован",
			body:            "https://practicum.yandex.ru",
			method:          http.MethodPost,
			contentType:     "text/plain",
			contentEncoding: "",
			acceptEncoding:  "",
			// authorization:   constants.JWTTokenTestUser,
			want: want{
				statusCode:  http.StatusCreated,
				response:    "http://localhost:8080/bTKNZu94",
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
			// authorization:   constants.JWTTokenTestUser,
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
			// authorization:   constants.JWTTokenTestUser,
			want: want{
				statusCode:  http.StatusBadRequest,
				response:    "Не удалось прочитать тело запроса\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tt.method
			req.URL = server.URL
			req.Body = tt.body
			req.SetHeaders(map[string]string{
				"Content-Type":     tt.contentType,
				"Content-Encoding": tt.contentEncoding,
				"Accept-Encoding":  tt.acceptEncoding,
				// "Authorization":    tt.authorization,
			})

			resp, err := req.Send()
			require.NoError(t, err, "error making HTTP request", tt.body)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, string(resp.Body()))
			}
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
			}
		})
	}
}

func testCreateURLFromJSON(t *testing.T, server *httptest.Server) {
	type want struct {
		statusCode      int
		response        string
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
				response:    "Не удалось прочитать тело запроса\n",
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
				response:    "json: cannot unmarshal number into Go struct field",
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
				response:        `{"result":"http://localhost:8080/4BeKySvE"}`,
				contentType:     "application/json",
				contentEncoding: "gzip",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			if tt.want.response != "" {
				checkBody := strings.Contains(string(resp.Body()), tt.want.response)
				assert.True(t, checkBody)
			}

			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
			}
		})
	}
}

func testGetURL(t *testing.T, server *httptest.Server) {
	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name     string
		shortURL string
		method   string
		want     want
	}{
		{
			name:     "GET, адрес - существующий в БД адрес, кейс 1",
			shortURL: "bTKNZu94",
			method:   http.MethodGet,
			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				response:    "https://practicum.yandex.ru",
				contentType: "text/plain",
			},
		},
		{
			name:     "GET, адрес - существующий в БД адрес, кейс 2",
			shortURL: "4BeKySvE",
			method:   http.MethodGet,
			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				response:    "https://www.sports.ru",
				contentType: "text/plain",
			},
		},
		{
			name:     "POST, адрес - существующий в БД адрес",
			shortURL: "bTKNZu94",
			method:   http.MethodPost,
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name:     "GET, адрес в БД не найден",
			shortURL: "nvjkrhsf",
			method:   http.MethodGet,
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).R()
			req.Method = tt.method
			req.URL = server.URL + "/" + tt.shortURL

			resp, err := req.Send()
			if err != nil {
				assert.True(t, strings.Contains(err.Error(), "auto redirect is disabled"))
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, string(resp.Header().Get("Location")))
			}
			if tt.want.contentType != "" {
				assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
			}
		})
	}
}

func testGetAllURLs(t *testing.T, server *httptest.Server) {
	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name          string
		authorization string
		method        string
		want          want
	}{
		{
			name:          "GET, кейс 1",
			// authorization: constants.JWTTokenTestUser,
			method:        http.MethodGet,
			want: want{
				statusCode:  http.StatusTemporaryRedirect,
				response:    `{"short_url":"http://localhost:8080/bTKNZu94","original_url":"https://practicum.yandex.ru}`,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).R()
			req.Method = tt.method
			// req.Header.Set("Authorization", tt.authorization)
			req.URL = server.URL + "/api/user/urls"

			// resp, err := req.Send()
			// if err != nil {
			// 	assert.True(t, strings.Contains(err.Error(), "auto redirect is disabled"))
			// }

			// assert.Equal(t, tt.want.statusCode, resp.StatusCode())
			// if tt.want.response != "" {
			// 	assert.Equal(t, tt.want.response, string(resp.Header().Get("Location")))
			// }
			// if tt.want.contentType != "" {
			// 	assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
			// }
		})
	}
}
