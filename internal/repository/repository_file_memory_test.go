package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/mocks"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/golang/mock/gomock"
)

// setupRepoFileMemory - тестирование .
func setupRepoFileMemory(storage *Storage) *RepoFileMemory {
	repo := NewRepoFileMemory(storage)
	repo.URLs = append(repo.URLs, testURLsFull...)
	return repo
}

func TestRepoFileMemory_Ping(t *testing.T) {
	tests := []struct {
		name string
		want error
	}{
		{
			name: "тест 1",
			want: constants.ErrorMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConsumer := mocks.NewMockReadCloser(ctrl)
			mockProducer := mocks.NewMockWriteCloser(ctrl)

			storage := &Storage{
				Consumer: mockConsumer,
				Producer: mockProducer,
			}

			repo := setupRepoFileMemory(storage)
			if got := repo.Ping(context.Background()); got != tt.want {
				t.Errorf("TestRepoFile_Ping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoFileMemory_InsertBatch(t *testing.T) {
	tests := []struct {
		name string
		urls []models.URLBase
		mock func(producer *mocks.MockWriteCloser)
		want error
	}{
		{
			name: "тест 1",
			urls: []models.URLBase{testURLFull1},
			mock: func(producer *mocks.MockWriteCloser) {},
			want: constants.ErrorURLAlreadyExist,
		},
		{
			name: "тест 2",
			urls: []models.URLBase{testURLFull3},
			mock: func(producer *mocks.MockWriteCloser) {
				producer.EXPECT().Write(testURLFull3).Return(nil)
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConsumer := mocks.NewMockReadCloser(ctrl)
			mockProducer := mocks.NewMockWriteCloser(ctrl)

			tt.mock(mockProducer)

			storage := &Storage{
				Consumer: mockConsumer,
				Producer: mockProducer,
			}

			repo := setupRepoFileMemory(storage)
			if got := repo.InsertBatch(context.Background(), tt.urls); got != tt.want {
				t.Errorf("TestRepoFileMemory_InsertBatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoFileMemory_InsertOrdinary(t *testing.T) {
	tests := []struct {
		name string
		url  models.URLBase
		want error
		mock func(producer *mocks.MockWriteCloser)
	}{
		{
			name: "тест 1",
			url:  testURLFull1,
			want: constants.ErrorURLAlreadyExist,
			mock: func(producer *mocks.MockWriteCloser) {},
		},
		{
			name: "тест 2",
			url:  testURLFull3,
			want: nil,
			mock: func(producer *mocks.MockWriteCloser) {
				producer.EXPECT().Write(testURLFull3).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConsumer := mocks.NewMockReadCloser(ctrl)
			mockProducer := mocks.NewMockWriteCloser(ctrl)

			tt.mock(mockProducer)

			storage := &Storage{
				Consumer: mockConsumer,
				Producer: mockProducer,
			}

			repo := setupRepoFileMemory(storage)
			if got := repo.InsertOrdinary(context.Background(), tt.url); got != tt.want {
				t.Errorf("TestRepoFileMemory_InsertOrdinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepoFileMemory_SelectOriginal(t *testing.T) {
	tests := []struct {
		name     string
		shortURL string
		want     string
		wantErr  error
	}{
		{
			name:     "тест 1",
			shortURL: urlAlias4,
			want:     "",
			wantErr:  constants.ErrorURLAlreadyDeleted,
		},
		{
			name:     "тест 2",
			shortURL: urlAlias1,
			want:     url1,
			wantErr:  nil,
		},
		{
			name:     "тест 3",
			shortURL: urlAlias3,
			want:     "",
			wantErr:  constants.ErrorURLNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConsumer := mocks.NewMockReadCloser(ctrl)
			mockProducer := mocks.NewMockWriteCloser(ctrl)

			storage := &Storage{
				Consumer: mockConsumer,
				Producer: mockProducer,
			}

			repo := setupRepoFileMemory(storage)
			got, gotErr := repo.SelectOriginal(context.Background(), tt.shortURL)
			if got != tt.want || gotErr != tt.wantErr {
				t.Errorf("TestRepoFileMemory_SelectOriginal() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestRepoFileMemory_SelectShort(t *testing.T) {
	tests := []struct {
		name        string
		originalURL string
		want        string
		wantErr     error
	}{
		{
			name:        "тест 1",
			originalURL: url1,
			want:        urlAlias1,
			wantErr:     nil,
		},
		{
			name:        "тест 2",
			originalURL: url3,
			want:        "",
			wantErr:     constants.ErrorURLNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConsumer := mocks.NewMockReadCloser(ctrl)
			mockProducer := mocks.NewMockWriteCloser(ctrl)

			storage := &Storage{
				Consumer: mockConsumer,
				Producer: mockProducer,
			}

			repo := setupRepoFileMemory(storage)
			got, gotErr := repo.SelectShort(context.Background(), tt.originalURL)
			if got != tt.want || gotErr != tt.wantErr {
				t.Errorf("TestRepoFileMemory_SelectShort() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestRepoFileMemory_SelectAll(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		want    []models.URLBase
		wantErr error
	}{
		{
			name:    "тест 1",
			userID:  UUID,
			want:    testURLsShort,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConsumer := mocks.NewMockReadCloser(ctrl)
			mockProducer := mocks.NewMockWriteCloser(ctrl)

			storage := &Storage{
				Consumer: mockConsumer,
				Producer: mockProducer,
			}

			repo := setupRepoFileMemory(storage)
			got, gotErr := repo.SelectAll(context.Background(), tt.userID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestRepoFileMemory_SelectAll() = %v, want %v", got, tt.want)
			}

			if gotErr != tt.wantErr {
				t.Errorf("TestRepoFileMemory_SelectAll() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestRepoFileMemory_Delete(t *testing.T) {
	tests := []struct {
		name string
		urls []models.URLBase
		want error
	}{
		{
			name: "тест 1",
			urls: []models.URLBase{testURLFull1},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConsumer := mocks.NewMockReadCloser(ctrl)
			mockProducer := mocks.NewMockWriteCloser(ctrl)

			storage := &Storage{
				Consumer: mockConsumer,
				Producer: mockProducer,
			}

			repo := setupRepoFileMemory(storage)
			if got := repo.Delete(context.Background(), tt.urls); got != tt.want {
				t.Errorf("TestRepoFileMemory_Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
