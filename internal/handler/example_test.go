package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/mocks"
)

func getExampleMocks(ctrl *gomock.Controller) *mocks.MockURLUseCase {
	mock := mocks.NewMockURLUseCase(ctrl)

	mock.EXPECT().Ping(gomock.Any()).Return(nil).AnyTimes()
	mock.EXPECT().CreateURLOrdinary(gomock.Any(), urlIn3).Return(urlOut3, nil).AnyTimes()
	mock.EXPECT().CreateURLOrdinary(gomock.Any(), urlIn3).Return(urlOut3, nil).AnyTimes()
	mock.EXPECT().CreateURLOrdinary(gomock.Any(), urlIn4).Return(urlOut4, nil).AnyTimes()
	mock.EXPECT().CreateURLBatch(gomock.Any(), urlsIn1).Return(urlsOut1, nil).AnyTimes()
	mock.EXPECT().GetOriginalURL(gomock.Any(), urlShort1).Return(urlOriginal1, nil).AnyTimes()
	mock.EXPECT().GetAllURLs(gomock.Any(), UUID).Return(urlsOut2, nil).AnyTimes()
	mock.EXPECT().DeleteURLs(gomock.Any(), urlsIn2).Return(nil).AnyTimes()
	return mock
}

func setupServer() *httptest.Server {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	cfg := &config.Config{
		UseMockAuth:     true,
		ServerAddress:   "localhost:8081",
		BaseURL:         "http://localhost:8081",
		FileStoragePath: "../../database_test.log",
	}

	useCase := getExampleMocks(ctrl)
	handler := NewСontroller(useCase, cfg)
	router := handler.SetupRouter()
	return httptest.NewServer(router)
}

func ExampleController_CreateURLShortJSONBatch() {
	srv := setupServer()

	req := resty.New().R()
	req.Body = bodyJSONBatch
	req.URL = srv.URL + "/api/shorten/batch"
	req.Method = http.MethodPost

	ctx := context.WithValue(req.Context(), constants.UserIDKey, UUID)
	req.SetContext(ctx)

	req.SetHeaders(map[string]string{
		"Content-Type": "application/json",
	})

	resp, _ := req.Send()

	fmt.Println(string(resp.Body()))
	fmt.Println(resp.StatusCode())
	fmt.Println(resp.Header().Get("Content-Type"))

	// Output:
	// [{"short_url":"http://localhost:8081/lJJpJV7h","correlation_id":"1"},{"short_url":"http://localhost:8081/kiFL71uv","correlation_id":"2"}]
	// 201
	// application/json
}

func ExampleController_createURLShortJSON() {
	srv := setupServer()

	req := resty.New().R()
	req.Body = bodyJSON1
	req.URL = srv.URL + "/api/shorten"
	req.Method = http.MethodPost

	ctx := context.WithValue(req.Context(), constants.UserIDKey, UUID)
	req.SetContext(ctx)

	req.SetHeaders(map[string]string{
		"Content-Type": "application/json",
	})

	resp, _ := req.Send()

	fmt.Println(string(resp.Body()))
	fmt.Println(resp.StatusCode())
	fmt.Println(resp.Header().Get("Content-Type"))

	// Output:
	// {"result":"http://localhost:8081/lJJpJV7h"}
	// 201
	// application/json
}

func ExampleController_createURLShortText() {
	srv := setupServer()

	req := resty.New().R()
	req.Body = bodyText2
	req.URL = srv.URL + "/"
	req.Method = http.MethodPost

	ctx := context.WithValue(req.Context(), constants.UserIDKey, UUID)
	req.SetContext(ctx)

	req.SetHeaders(map[string]string{
		"Content-Type": "text/plain",
	})

	resp, _ := req.Send()

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.Header().Get("Content-Type"))
	// TODO проверить на запросе
	// fmt.Println(resp.Header().Get("Content-Length"))

	// Output:
	// 201
	// text/plain
}

func ExampleController_getAllURLs() {
	srv := setupServer()

	req := resty.New().R()
	req.URL = srv.URL + "/api/user/urls"
	req.Method = http.MethodGet

	ctx := context.WithValue(req.Context(), constants.UserIDKey, UUID)
	req.SetContext(ctx)

	resp, _ := req.Send()

	fmt.Println(string(resp.Body()))
	fmt.Println(resp.StatusCode())
	fmt.Println(resp.Header().Get("Content-Type"))

	// Output:
	// [{"short_url":"http://localhost:8081/lJJpJV7h","original_url":"https://www.khl.ru/"}]
	// 200
	// application/json
}

func ExampleController_getURLOriginal() {
	srv := setupServer()

	req := resty.New().SetRedirectPolicy(resty.NoRedirectPolicy()).R()
	req.URL = srv.URL + "/" + urlShort1
	req.Method = http.MethodGet

	ctx := context.WithValue(req.Context(), constants.UserIDKey, UUID)
	req.SetContext(ctx)

	resp, _ := req.Send()

	fmt.Println(resp.StatusCode())
	fmt.Println(resp.Header().Get("Location"))

	// Output:
	// 307
	// https://www.khl.ru/

}

func ExampleController_deleteURLs() {
	srv := setupServer()

	req := resty.New().R()
	req.URL = srv.URL + "/api/user/urls"
	req.Body = bodyJSON2
	req.Method = http.MethodDelete

	ctx := context.WithValue(req.Context(), constants.UserIDKey, UUID)
	req.SetContext(ctx)

	resp, _ := req.Send()

	fmt.Println(resp.StatusCode())

	// Output:
	// 202
}

func ExampleController_pingDB() {
	srv := setupServer()

	req := resty.New().R()
	req.URL = srv.URL + "/ping"
	req.Method = http.MethodGet

	resp, _ := req.Send()

	fmt.Println(resp.StatusCode())

	// Output:
	// 200
}
