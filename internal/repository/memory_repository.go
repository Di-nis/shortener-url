package repository

import (
	"context"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
)

// RepoMemory - структура базы данных.
type RepoMemory struct {
	OriginalAndShortUrl []models.URL
}

// NewRepoMemory - создание структуры Repo.
func NewRepoMemory() *RepoMemory {
	return &RepoMemory{
		OriginalAndShortUrl: make([]models.URL, 0),
	}
}

// CreateBatch - сохранение URL в базу данных.
func (repo *RepoMemory) CreateBatch(ctx context.Context, urls []models.URL) error {
	for _, url := range urls {
		for _, urlDB := range repo.OriginalAndShortUrl {
			if urlDB.Original == url.Original {
				return constants.ErrorURLAlreadyExist
			}
		}

		repo.OriginalAndShortUrl = append(repo.OriginalAndShortUrl, url)
	}
	return nil
}

// CreateOrdinaty - сохранение URL в базу данных.
func (repo *RepoMemory) CreateOrdinaty(ctx context.Context, urlIn models.URL) error {
	for _, urlDB := range repo.OriginalAndShortUrl {
		if urlDB.Original == urlIn.Original {
			return constants.ErrorURLAlreadyExist
		}
	}

	repo.OriginalAndShortUrl = append(repo.OriginalAndShortUrl, urlIn)
	return nil

}

// GetOriginalURL - получение оригинального URL из базы данных.
func (repo *RepoMemory) GetOriginalURL(ctx context.Context, shortUrl string) (string, error) {
	for _, urlData := range repo.OriginalAndShortUrl {
		if urlData.Short == shortUrl {
			return urlData.Original, nil
		}
	}
	return "", constants.ErrorURLNotExist
}

// GetShortURL - получение оригинального URL из базы данных.
func (repo *RepoMemory) GetShortURL(ctx context.Context, originalUrl string) (string, error) {
	for _, urlData := range repo.OriginalAndShortUrl {
		if urlData.Original == originalUrl {
			return urlData.Short, nil
		}
	}
	return "", constants.ErrorURLNotExist
}
