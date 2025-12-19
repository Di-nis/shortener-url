package handlerbenchmark_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/handler"
	"github.com/Di-nis/shortener-url/internal/mocks"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
)

var (
	urlOriginal1 = "https://www.khl.ru/"
	urlShort1    = "lJJpJV7h"

	bodyJSONBatch = `[{"correlation_id":"1","original_url":"https://www.khl.ru/"},{"correlation_id":"2","original_url":"https://www.dynamo.ru/"}]`
	bodyJSON1     = `{"url":"https://www.khl.ru/"}`
	bodyText1     = `https://maximum.ru/`
)

func getBenchmarkMocks(ctrl *gomock.Controller) *mocks.MockURLUseCase {
	mock := mocks.NewMockURLUseCase(ctrl)

	mock.EXPECT().Ping(gomock.Any()).Return(nil).AnyTimes()
	mock.EXPECT().CreateURLOrdinary(gomock.Any(), gomock.Any()).Return(models.URLBase{}, nil).AnyTimes()
	mock.EXPECT().CreateURLBatch(gomock.Any(), gomock.AssignableToTypeOf([]models.URLBase{})).Return([]models.URLBase{}, nil).AnyTimes()
	mock.EXPECT().GetOriginalURL(gomock.Any(), urlShort1).Return(urlOriginal1, nil).AnyTimes()
	mock.EXPECT().GetAllURLs(gomock.Any(), gomock.AssignableToTypeOf("")).Return([]models.URLBase{}, nil).AnyTimes()
	mock.EXPECT().DeleteURLs(gomock.Any(), gomock.AssignableToTypeOf([]models.URLBase{})).Return(nil).AnyTimes()
	return mock
}

func BenchmarkHandler(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	cfg := config.NewConfig()
	cfg.Load()

	useCase := getBenchmarkMocks(ctrl)
	handler := handler.NewСontroller(useCase, cfg)
	router := handler.SetupRouter()
	server := httptest.NewServer(router)

	client := resty.New()

	b.Run("createURLShortJSONBatch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := client.R()
			req.Body = bodyJSONBatch
			req.Method = http.MethodPost
			req.URL = server.URL + "/api/shorten/batch"

			req.Send()
		}
	})

	b.Run("createURLShortJSON", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := client.R()
			req.Body = bodyJSON1
			req.Method = http.MethodPost
			req.URL = server.URL + "/api/shorten"
			req.Send()
		}
	})

	b.Run("createURLShortText", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := client.R()
			req.Body = bodyText1
			req.Method = http.MethodPost
			req.URL = server.URL
			req.Send()
		}
	})

	b.Run("getAllURLs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := client.R()
			req.Body = bodyText1
			req.Method = http.MethodGet
			req.URL = server.URL + "/api/user/urls"
			req.Send()
		}
	})

	b.Run("getURLOriginal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := client.R()
			req.Body = bodyText1
			req.Method = http.MethodGet
			req.URL = server.URL + "/" + urlShort1
			req.Send()
		}
	})

	b.Run("pingDB", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := client.R()
			req.Method = http.MethodGet
			req.URL = server.URL + "/ping"

			req.Send()
		}

	})

	b.Run("deleteURLs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := client.R()
			req.Body = bodyText1
			req.Method = http.MethodDelete
			req.URL = server.URL + "/api/user/urls"
			req.Send()
		}
	})
}
