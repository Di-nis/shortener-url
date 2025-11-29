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
func setupRepoFileMemory(t *testing.T) *RepoFileMemory {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	consumer := mocks.NewMockReadCloser(ctrl)
	producer := mocks.NewMockWriteCloser(ctrl)

	storage := &Storage{
		Consumer: consumer,
		Producer: producer,
	}
	repo := NewRepoFileMemory(storage)
	repo.URLs = append(repo.URLs, urlsOut1...)
	repo.URLs = append(repo.URLs, urlOut4)
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
			repo := setupRepoFileMemory(t)
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
		want error
	}{
		{
			name: "тест 1",
			urls: []models.URLBase{urlOut1},
			want: constants.ErrorURLAlreadyExist,
		},
		{
			name: "тест 2",
			urls: []models.URLBase{urlOut3},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupRepoFileMemory(t)
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
	}{
		{
			name: "тест 1",
			url:  urlOut1,
			want: constants.ErrorURLAlreadyExist,
		},
		{
			name: "тест 2",
			url:  urlOut3,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupRepoFileMemory(t)
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
			shortURL: urlShort4,
			want:     "",
			wantErr:  constants.ErrorURLAlreadyDeleted,
		},
		{
			name:     "тест 2",
			shortURL: urlShort1,
			want:     urlOriginal1,
			wantErr:  nil,
		},
		{
			name:     "тест 3",
			shortURL: urlShort3,
			want:     "",
			wantErr:  constants.ErrorURLNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupRepoFileMemory(t)
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
			originalURL: urlOriginal1,
			want:        urlShort1,
			wantErr:     nil,
		},
		{
			name:        "тест 2",
			originalURL: urlOriginal3,
			want:        "",
			wantErr:     constants.ErrorURLNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupRepoFileMemory(t)
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
			want:    urlsOut3,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupRepoFileMemory(t)
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
			urls: []models.URLBase{urlOut1},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupRepoFileMemory(t)
			if got := repo.Delete(context.Background(), tt.urls); got != tt.want {
				t.Errorf("TestRepoFileMemory_Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
