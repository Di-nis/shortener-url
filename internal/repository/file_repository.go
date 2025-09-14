package repository

import (
	"context"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
)

type WriteCloser interface {
	WriteURL(models.URL) error
	SaveToFile(models.URL) error
	Close() error
}

type ReadCloser interface {
	ReadURL() (*models.URL, error)
	LoadFromFile() ([]models.URL, error)
	Close() error
}

type Storage struct {
	Producer WriteCloser
	Consumer ReadCloser
}

// RepoFile - структура базы данных.
type RepoFile struct {
	OriginalAndShortURL []models.URL
	FileStoragePath     string
	Storage             *Storage
}

// NewRepoFile - создание структуры Repo.
func NewRepoFile(fileStoragePath string, storage *Storage) *RepoFile {
	return &RepoFile{
		OriginalAndShortURL: make([]models.URL, 0),
		FileStoragePath:     fileStoragePath,
		Storage:             storage,
	}
}

// CreateBatch - сохранение URL в базу данных.
func (repo *RepoFile) CreateBatch(ctx context.Context, urls []models.URL) error {
	for _, url := range urls {
		for _, urlDB := range repo.OriginalAndShortURL {
			if urlDB.Original == url.Original {
				return constants.ErrorURLAlreadyExist
			}
		}

		repo.OriginalAndShortURL = append(repo.OriginalAndShortURL, url)

		err := repo.Storage.Producer.SaveToFile(url)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateOrdinaty - сохранение URL в базу данных.
func (repo *RepoFile) CreateOrdinaty(ctx context.Context, url models.URL) error {
	for _, urlDB := range repo.OriginalAndShortURL {
		if urlDB.Original == url.Original {
			return constants.ErrorURLAlreadyExist
		}
	}

	repo.OriginalAndShortURL = append(repo.OriginalAndShortURL, url)

	err := repo.Storage.Producer.SaveToFile(url)
		if err != nil {
			return err
		}
	return nil

}

// GetOriginalURL - получение оригинального URL из базы данных.
func (repo *RepoFile) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	for _, urlData := range repo.OriginalAndShortURL {
		if urlData.Short == shortURL {
			return urlData.Original, nil
		}
	}
	return "", constants.ErrorURLNotExist
}

// GetShortURL - получение оригинального URL из базы данных.
func (repo *RepoFile) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	for _, urlData := range repo.OriginalAndShortURL {
		if urlData.Original == originalURL {
			return urlData.Short, nil
		}
	}
	return "", constants.ErrorURLNotExist
}
