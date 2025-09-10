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
	URLOriginalAndShort []models.URL
	FileStoragePath     string
	Storage             *Storage
}

// NewRepoFile - создание структуры Repo.
func NewRepoFile(fileStoragePath string, storage *Storage) *RepoFile {
	return &RepoFile{
		URLOriginalAndShort: make([]models.URL, 0),
		FileStoragePath:     fileStoragePath,
		Storage:             storage,
	}
}

// Create - сохранение URL в базу данных.
func (repo *RepoFile) Create(ctx context.Context, urlOriginal, urlShort string) error {
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

	err := repo.Storage.Producer.SaveToFile(urlData)
	if err != nil {
		return err
	}
	return nil
}

// Get - получение оригинального URL из базы данных.
func (repo *RepoFile) Get(ctx context.Context, urlShort string) (string, error) {
	for _, urlData := range repo.URLOriginalAndShort {
		if urlData.URLShort == urlShort {
			return urlData.URLOriginal, nil
		}
	}
	return "", constants.ErrorURLNotExist
}
