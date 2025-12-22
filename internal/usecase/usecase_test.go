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

func TestURLUseCase_Ping(t *testing.T) {
	tests := []struct {
		name string
		mock func(*mocks.MockURLRepository)
		want error
	}{
		{
			name: "ping - ok",
			mock: func(mockRepo *mocks.MockURLRepository) {
				mockRepo.EXPECT().Ping(gomock.Any()).Return(nil)
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockURLRepository(ctrl)
		tt.mock(mockRepo)

		service := service.NewService()
		useCase := NewURLUseCase(mockRepo, service)

		if got := useCase.Ping(context.Background()); got != tt.want {
			t.Errorf("GetOrderInfo() = %v, want %v", got, tt.want)
		}
	}
}

func TestURLUseCase_CreateURLOrdinary(t *testing.T) {
	tests := []struct {
		name    string
		urlIn   any
		mock    func(*mocks.MockURLRepository)
		want    models.URLBase
		wantErr error
	}{
		{
			name:  "создание короткого URL, кейс 1",
			urlIn: urlIn1,
			mock: func(mockRepo *mocks.MockURLRepository) {
				mockRepo.EXPECT().InsertOrdinary(gomock.Any(), urlOut1).Return(nil)
			},
			want:    urlOut1,
			wantErr: nil,
		},
		{
			name:  "создание короткого URL, кейс 2",
			urlIn: urlIn2,
			mock: func(mockRepo *mocks.MockURLRepository) {
				mockRepo.EXPECT().InsertOrdinary(gomock.Any(), urlOut2).Return(constants.ErrorURLAlreadyExist)
				mockRepo.EXPECT().SelectShort(gomock.Any(), urlOriginal2).Return(urlShort2, nil)
			},
			want:    urlOut2,
			wantErr: constants.ErrorURLAlreadyExist,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockURLRepository(ctrl)
		tt.mock(mockRepo)

		service := service.NewService()
		useCase := NewURLUseCase(mockRepo, service)

		got, gotErr := useCase.CreateURLOrdinary(context.Background(), tt.urlIn)
		if got != tt.want {
			t.Errorf("CreateURLOrdinary() = %v, want %v", got, tt.want)
		}
		if gotErr != tt.wantErr {
			t.Errorf("CreateURLOrdinary() = %v, wantErr %v", gotErr, tt.wantErr)
		}
	}
}

func TestURLUseCase_CreateURLBatch(t *testing.T) {
	tests := []struct {
		name    string
		urls    []models.URLBase
		mock    func(*mocks.MockURLRepository)
		want    []models.URLBase
		wantErr error
	}{
		{
			name: "создание коротких URL (batch), кейс 1",
			urls: urlsIn,
			mock: func(mockRepo *mocks.MockURLRepository) {
				mockRepo.EXPECT().InsertBatch(gomock.Any(), urlsOut).Return(nil)
			},
			want:    urlsOut,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockURLRepository(ctrl)
		tt.mock(mockRepo)

		service := service.NewService()
		useCase := NewURLUseCase(mockRepo, service)

		got, gotErr := useCase.CreateURLBatch(context.Background(), tt.urls)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CreateURLBatch() = %v, want %v", got, tt.want)
		}
		if gotErr != tt.wantErr {
			t.Errorf("CreateURLBatch() = %v, wantErr %v", gotErr, tt.wantErr)
		}
	}
}

func TestURLUseCase_GetOriginalURL(t *testing.T) {
	tests := []struct {
		name     string
		shortURL string
		mock     func(*mocks.MockURLRepository)
		want     string
		wantErr  error
	}{
		{
			name:     "получение оригинального URL, кейс 1",
			shortURL: urlShort1,
			mock: func(mockRepo *mocks.MockURLRepository) {
				mockRepo.EXPECT().SelectOriginal(gomock.Any(), urlShort1).Return(urlOriginal1, nil)
			},
			want:    urlOriginal1,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockURLRepository(ctrl)
		tt.mock(mockRepo)

		service := service.NewService()
		useCase := NewURLUseCase(mockRepo, service)

		got, gotErr := useCase.GetOriginalURL(context.Background(), tt.shortURL)
		if got != tt.want {
			t.Errorf("GetOriginalURL() = %v, want %v", got, tt.want)
		}

		if gotErr != tt.wantErr {
			t.Errorf("GetOriginalURL() = %v, wantErr %v", gotErr, tt.wantErr)
		}
	}
}

func TestURLUseCase_GetAllURLs(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		mock    func(*mocks.MockURLRepository)
		want    []models.URLBase
		wantErr error
	}{
		{
			name:   "получение всех URL пользователя, кейс 1",
			userID: UUID,
			mock: func(mockRepo *mocks.MockURLRepository) {
				mockRepo.EXPECT().SelectAll(gomock.Any(), UUID).Return(urlsOut, nil)
			},
			want:    urlsOut,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockURLRepository(ctrl)
		tt.mock(mockRepo)

		service := service.NewService()
		useCase := NewURLUseCase(mockRepo, service)

		got, gotErr := useCase.GetAllURLs(context.Background(), tt.userID)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("GetAllURLs() = %v, want %v", got, tt.want)
		}

		if gotErr != tt.wantErr {
			t.Errorf("GetAllURLs() = %v, wantErr %v", gotErr, tt.wantErr)
		}
	}
}

func TestURLUseCase_DeleteURLs(t *testing.T) {
	tests := []struct {
		name    string
		urls    []models.URLBase
		mock    func(*mocks.MockURLRepository)
		wantErr error
	}{
		{
			name: "удаление URL, кейс 1",
			urls: urlsOut,
			mock: func(mockRepo *mocks.MockURLRepository) {
				mockRepo.EXPECT().Delete(gomock.Any(), urlsOut).Return(nil)
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockURLRepository(ctrl)
		tt.mock(mockRepo)

		service := service.NewService()
		useCase := NewURLUseCase(mockRepo, service)

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
