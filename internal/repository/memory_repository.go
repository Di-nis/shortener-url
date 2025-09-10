package repository

import (
	"context"

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
func (repo *RepoMemory) Create(ctx context.Context, urlOriginal, urlShort string) error {
	for _, urlData := range repo.URLOriginalAndShort {
		if urlData.URLOriginal == urlOriginal {
			return constants.ErrorURLAlreadyExist
		}
	}

	urlData := models.URL{
		URLShort:    urlShort,
		URLOriginal: urlOriginal,
	}

	repo.URLOriginalAndShort = append(repo.URLOriginalAndShort, urlData)
	return nil
}

// Get - получение оригинального URL из базы данных.
func (repo *RepoMemory) Get(ctx context.Context, urlShort string) (string, error) {
	for _, urlData := range repo.URLOriginalAndShort {
		if urlData.URLShort == urlShort {
			return urlData.URLOriginal, nil
		}
	}
	return "", constants.ErrorURLNotExist
}
