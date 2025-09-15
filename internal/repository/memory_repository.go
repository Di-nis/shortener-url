package repository

import (
	"context"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
)

// RepoMemory - структура базы данных.
type RepoMemory struct {
	OriginalAndShortURL []models.URL
}

// NewRepoMemory - создание структуры Repo.
func NewRepoMemory() *RepoMemory {
	return &RepoMemory{
		OriginalAndShortURL: make([]models.URL, 0),
	}
}

// CreateBatch - сохранение URL в базу данных.
func (repo *RepoMemory) CreateBatch(ctx context.Context, urls []models.URL) error {
	for _, url := range urls {
		for _, urlDB := range repo.OriginalAndShortURL {
			if urlDB.Original == url.Original {
				return constants.ErrorURLAlreadyExist
			}
		}

		repo.OriginalAndShortURL = append(repo.OriginalAndShortURL, url)
	}
	return nil
}

// CreateOrdinaty - сохранение URL в базу данных.
func (repo *RepoMemory) CreateOrdinaty(ctx context.Context, urlIn models.URL) error {
	for _, urlDB := range repo.OriginalAndShortURL {
		if urlDB.Original == urlIn.Original {
			return constants.ErrorURLAlreadyExist
		}
	}

	repo.OriginalAndShortURL = append(repo.OriginalAndShortURL, urlIn)
	return nil

}

// GetOriginalURL - получение оригинального URL из базы данных.
func (repo *RepoMemory) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	for _, urlData := range repo.OriginalAndShortURL {
		if urlData.Short == shortURL {
			return urlData.Original, nil
		}
	}
	return "", constants.ErrorURLNotExist
}

// GetShortURL - получение оригинального URL из базы данных.
func (repo *RepoMemory) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	for _, urlData := range repo.OriginalAndShortURL {
		if urlData.Original == originalURL {
			return urlData.Short, nil
		}
	}
	return "", constants.ErrorURLNotExist
}
