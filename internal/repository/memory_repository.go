package repository

import (
	"context"
	"slices"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
)

// RepoMemory - структура базы данных.
type RepoMemory struct {
	URLOriginalAndShort []models.URL
}

// NewRepoMemory - создание структуры Repo.
func NewRepoMemory() *RepoMemory {
	return &RepoMemory{
		URLOriginalAndShort: make([]models.URL, 0),
	}
}

// Create - сохранение URL в базу данных.
func (repo *RepoMemory) Create(ctx context.Context, urls []models.URL) error {
	for _, url := range urls {
		if slices.Contains(repo.URLOriginalAndShort, url) {
			return constants.ErrorURLAlreadyExist
		}
		repo.URLOriginalAndShort = append(repo.URLOriginalAndShort, url)
	}
	return nil
}

// Get - получение оригинального URL из базы данных.
func (repo *RepoMemory) Get(ctx context.Context, urlShort string) (string, error) {
	for _, urlData := range repo.URLOriginalAndShort {
		if urlData.Short == urlShort {
			return urlData.Original, nil
		}
	}
	return "", constants.ErrorURLNotExist
}
