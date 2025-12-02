package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/mocks"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/golang/mock/gomock"
)

var (
	UUID         = "01KA3YRQCWTNAJEGR5Z30PH6VT"
	urlOriginal1 = "https://www.khl.ru/"
	urlShort1    = "lJJpJV7h"
	urlOriginal2 = "https://www.dynamo.ru/"
	urlShort2    = "kiFL71uv"
	urlOriginal3 = "https://chatgpt.com/"
	urlShort3    = "826drChJ"

	urlIn1 = models.URLBase{
		UUID:        UUID,
		Original:    urlOriginal1,
		DeletedFlag: false,
	}
	urlIn2 = models.URLBase{
		UUID:        UUID,
		Original:    urlOriginal2,
		DeletedFlag: false,
	}

	urlOut1 = models.URLBase{
		UUID:        UUID,
		Original:    urlOriginal1,
		Short:       urlShort1,
		DeletedFlag: false,
	}
	urlOut2 = models.URLBase{
		UUID:        UUID,
		Original:    urlOriginal2,
		Short:       urlShort2,
		DeletedFlag: false,
	}

	urlsIn = []models.URLBase{
		{
			UUID:        UUID,
			Original:    urlOriginal1,
			DeletedFlag: false,
		},
	}

	urlsOut = []models.URLBase{
		{
			UUID:        UUID,
			Original:    urlOriginal1,
			Short:       urlShort1,
			DeletedFlag: false,
		},
	}
)

// TestService - тестирование сервиса.
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

// getMocks - формирование мок.
func getMocks(ctrl *gomock.Controller) *mocks.MockURLRepository {
	mockRepository := mocks.NewMockURLRepository(ctrl)

	mockRepository.EXPECT().Ping(gomock.Any()).Return(nil)
	mockRepository.EXPECT().InsertOrdinary(gomock.Any(), urlOut1).Return(nil)
	mockRepository.EXPECT().InsertOrdinary(gomock.Any(), urlOut2).Return(constants.ErrorURLAlreadyExist)
	mockRepository.EXPECT().InsertBatch(gomock.Any(), urlsOut).Return(nil)
	mockRepository.EXPECT().SelectShort(gomock.Any(), urlOriginal2).Return(urlShort2, nil)
	mockRepository.EXPECT().SelectOriginal(gomock.Any(), urlShort1).Return(urlOriginal1, nil)
	mockRepository.EXPECT().SelectAll(gomock.Any(), UUID).Return(urlsOut, nil)
	mockRepository.EXPECT().Delete(gomock.Any(), urlsOut).Return(nil)
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
		want    models.URLBase
		wantErr error
	}{
		{
			name:    "создание короткого URL, кейс 1",
			urlIn:   urlIn1,
			want:    urlOut1,
			wantErr: nil,
		},
		{
			name:    "создание короткого URL, кейс 2",
			urlIn:   urlIn2,
			want:    urlOut2,
			wantErr: constants.ErrorURLAlreadyExist,
		},
	}
	for _, tt := range tests {
		got, gotErr := useCase.CreateURLOrdinary(context.Background(), tt.urlIn)
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
		urls    []models.URLBase
		want    []models.URLBase
		wantErr error
	}{
		{
			name:    "создание коротких URL (batch), кейс 1",
			urls:    urlsIn,
			want:    urlsOut,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		got, gotErr := useCase.CreateURLBatch(context.Background(), tt.urls)
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
		want    []models.URLBase
		wantErr error
	}{
		{
			name:    "получение всех URL пользователя, кейс 1",
			userID:  UUID,
			want:    urlsOut,
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
		urls    []models.URLBase
		wantErr error
	}{
		{
			name:    "удаление URL, кейс 1",
			urls:    urlsOut,
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

func BenchmarkService(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepository := getBenchmarkMocks(ctrl)

	service := service.NewService()
	useCase := NewURLUseCase(mockRepository, service)

	b.Run("Ping", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			useCase.Ping(ctx)
		}
	})
	b.Run("CreateURLOrdinary", func(b *testing.B) {
		tests := []struct {
			name string
			url  any
		}{
			{
				name: "создание короткого URL, кейс 1",
				url:  urlIn1,
			},
			{
				name: "создание короткого URL, кейс 2",
				url:  urlIn2,
			},
		}

		for _, tt := range tests {
			for i := 0; i < b.N; i++ {
				useCase.CreateURLOrdinary(ctx, tt.url)
			}
		}

	})
	b.Run("CreateURLBatch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			useCase.CreateURLBatch(ctx, urlsOut)
		}
	})
	b.Run("GetOriginalURL", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			useCase.GetOriginalURL(ctx, urlShort1)
		}
	})
	b.Run("GetAllURLs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			useCase.GetAllURLs(ctx, UUID)
		}
	})
	b.Run("DeleteURLs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			useCase.DeleteURLs(ctx, urlsOut)
		}
	})
}

func getBenchmarkMocks(ctrl *gomock.Controller) *mocks.MockURLRepository {
	mockRepository := mocks.NewMockURLRepository(ctrl)

	mockRepository.EXPECT().Ping(gomock.Any()).Return(nil).AnyTimes()
	mockRepository.EXPECT().InsertOrdinary(gomock.Any(), urlOut1).Return(nil).AnyTimes()
	mockRepository.EXPECT().InsertOrdinary(gomock.Any(), urlOut2).Return(constants.ErrorURLAlreadyExist).AnyTimes()
	mockRepository.EXPECT().InsertBatch(gomock.Any(), urlsOut).Return(nil).AnyTimes().AnyTimes()
	mockRepository.EXPECT().SelectShort(gomock.Any(), urlOriginal2).Return(urlShort2, nil).AnyTimes()
	mockRepository.EXPECT().SelectShort(gomock.Any(), urlOriginal3).Return(urlShort3, nil).AnyTimes()
	mockRepository.EXPECT().SelectOriginal(gomock.Any(), urlShort1).Return(urlOriginal1, nil).AnyTimes()
	mockRepository.EXPECT().SelectAll(gomock.Any(), UUID).Return(urlsOut, nil).AnyTimes().AnyTimes()
	mockRepository.EXPECT().Delete(gomock.Any(), urlsOut).Return(nil).AnyTimes().AnyTimes()
	return mockRepository
}
