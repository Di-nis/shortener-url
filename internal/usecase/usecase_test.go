package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Di-nis/shortener-url/internal/mocks"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/golang/mock/gomock"
)

var (
	baseURL = "http://localhost:8080/"

	UUID         = "01KA3YRQCWTNAJEGR5Z30PH6VT"
	urlOriginal1 = "https://www.khl.ru/"
	urlShort1    = "lJJpJV7h"
	urlOriginal2 = "https://www.dynamo.ru/"
	urlShort2    = "kiFL71uv"

	urlIn1 = models.URL{
		UUID:        UUID,
		Original:    urlOriginal1,
		DeletedFlag: false,
	}
	urlIn2 = models.URL{
		UUID:        UUID,
		Original:    urlOriginal2,
		DeletedFlag: false,
	}

	url1 = models.URL{
		UUID:        UUID,
		Original:    urlOriginal1,
		Short:       urlShort1,
		DeletedFlag: false,
	}
	url2 = models.URL{
		UUID:        UUID,
		Original:    urlOriginal2,
		Short:       urlShort2,
		DeletedFlag: false,
	}

	urlsIn = []models.URL{
		{
			UUID:        UUID,
			Original:    urlOriginal1,
			DeletedFlag: false,
		},
		// {
		// 	UUID:        UUID,
		// 	Original:    urlOriginal2,
		// 	DeletedFlag: false,
		// },
	}

	urls = []models.URL{
		{
			UUID:        UUID,
			Original:    urlOriginal1,
			Short:       urlShort1,
			DeletedFlag: false,
		},
		// {
		// 	UUID:        UUID,
		// 	Original:    urlOriginal2,
		// 	Short:       urlShort2,
		// 	DeletedFlag: false,
		// },
	}
	err1 = &pgconn.PgError{
		Code: "23505",
	}
)

// TestService - формирование мок.
func TestService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := getMocks(ctrl)

	service := service.NewService()
	useCase := NewURLUseCase(mockRepository, service)

	t.Run("01_testPing", func(t *testing.T) {
		testPing(t, useCase)
	})
	t.Run("02_testCreateURLOrdinary", func(t *testing.T) {
		testCreateURLOrdinary(t, useCase)
	})
	t.Run("03_TestCreateURLBatch", func(t *testing.T) {
		testCreateURLBatch(t, useCase)
	})
	t.Run("04_testGetOriginalURL", func(t *testing.T) {
		testGetOriginalURL(t, useCase)
	})
	t.Run("05_testGetAllURLs", func(t *testing.T) {
		testGetAllURLs(t, useCase)
	})
	t.Run("06_testDeleteURLs", func(t *testing.T) {
		testDeleteURLs(t, useCase)
	})
}

func getMocks(ctrl *gomock.Controller) *mocks.MockURLRepository {
	mockRepository := mocks.NewMockURLRepository(ctrl)

	mockRepository.EXPECT().Ping(gomock.Any()).Return(nil)
	mockRepository.EXPECT().InsertOrdinary(gomock.Any(), url1).Return(nil)
	mockRepository.EXPECT().InsertOrdinary(gomock.Any(), url2).Return(err1)
	mockRepository.EXPECT().InsertBatch(gomock.Any(), urls).Return(nil)
	// mockRepository.EXPECT().SelectShort(gomock.Any(), urlOriginal1).Return(urlShort1, nil)
	mockRepository.EXPECT().SelectShort(gomock.Any(), urlOriginal2).Return(urlShort2, nil)
	mockRepository.EXPECT().SelectOriginal(gomock.Any(), urlShort1).Return(urlOriginal1, nil)
	mockRepository.EXPECT().SelectAll(gomock.Any(), UUID).Return(urls, nil)
	mockRepository.EXPECT().Delete(gomock.Any(), urls).Return(nil)
	return mockRepository
}

func testPing(t *testing.T, useCase *URLUseCase) {
	tests := []struct {
		name string
		want error
	}{
		{
			name: "ping - ok",
			want: nil,
		},
	}
	for _, tt := range tests {
		if got := useCase.Ping(context.Background()); got != tt.want {
			t.Errorf("GetOrderInfo() = %v, want %v", got, tt.want)
		}
	}
}

func testCreateURLOrdinary(t *testing.T, useCase *URLUseCase) {
	tests := []struct {
		name    string
		urlIn   any
		baseURL string
		want    models.URL
		wantErr error
	}{
		{
			name:    "создание короткого URL, кейс 1",
			urlIn:   urlIn1,
			baseURL: baseURL,
			want:    url1,
			wantErr: nil,
		},
		{
			name:    "создание короткого URL, кейс 2",
			urlIn:   urlIn2,
			baseURL: baseURL,
			want:    url2,
			wantErr: err1,
		},
	}
	for _, tt := range tests {
		got, gotErr := useCase.CreateURLOrdinary(context.Background(), tt.urlIn, tt.baseURL)
		if got != tt.want {
			t.Errorf("CreateURLOrdinary() = %v, want %v", got, tt.want)
		}
		if gotErr != tt.wantErr {
			t.Errorf("CreateURLOrdinary() = %v, wantErr %v", gotErr, tt.wantErr)
		}
	}
}

func testCreateURLBatch(t *testing.T, useCase *URLUseCase) {
	tests := []struct {
		name    string
		urls    []models.URL
		baseURL string
		want    []models.URL
		wantErr error
	}{
		{
			name:    "создание коротких URL (batch), кейс 1",
			urls:    urlsIn,
			baseURL: baseURL,
			want:    urls,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		got, gotErr := useCase.CreateURLBatch(context.Background(), tt.urls, tt.baseURL)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CreateURLBatch() = %v, want %v", got, tt.want)
		}
		if gotErr != tt.wantErr {
			t.Errorf("CreateURLBatch() = %v, wantErr %v", gotErr, tt.wantErr)
		}
	}
}

func testGetOriginalURL(t *testing.T, useCase *URLUseCase) {
	tests := []struct {
		name     string
		shortURL string
		want     string
		wantErr  error
	}{
		{
			name:     "получение оригинального URL, кейс 1",
			shortURL: urlShort1,
			want:     urlOriginal1,
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		got, gotErr := useCase.GetOriginalURL(context.Background(), tt.shortURL)
		if got != tt.want {
			t.Errorf("GetOriginalURL() = %v, want %v", got, tt.want)
		}

		if gotErr != tt.wantErr {
			t.Errorf("GetOriginalURL() = %v, wantErr %v", gotErr, tt.wantErr)
		}
	}
}

func testGetAllURLs(t *testing.T, useCase *URLUseCase) {
	tests := []struct {
		name    string
		userID  string
		want    []models.URL
		wantErr error
	}{
		{
			name:    "получение всех URL пользователя, кейс 1",
			userID:  UUID,
			want:    urls,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		got, gotErr := useCase.GetAllURLs(context.Background(), tt.userID)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("GetAllURLs() = %v, want %v", got, tt.want)
		}

		if gotErr != tt.wantErr {
			t.Errorf("GetAllURLs() = %v, wantErr %v", gotErr, tt.wantErr)
		}
	}
}

func testDeleteURLs(t *testing.T, useCase *URLUseCase) {
	tests := []struct {
		name    string
		urls    []models.URL
		wantErr error
	}{
		{
			name:    "удаление URL, кейс 1",
			urls:    urls,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		gotErr := useCase.DeleteURLs(context.Background(), tt.urls)
		if gotErr != tt.wantErr {
			t.Errorf("DeleteURLs() = %v, wantErr %v", gotErr, tt.wantErr)
		}
	}
}
