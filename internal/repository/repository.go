package repository

import (
	"github.com/Di-nis/shortener-url/internal/constants"
	"context"
)

type WriteCloser interface {
	WriteURL(URLData) error
	SaveToFile(URLData) error
	Close() error
}

type ReadCloser interface {
	ReadURL() (*URLData, error)
	LoadFromFile() ([]URLData, error)
	Close() error
}

type Storage struct {
	Producer WriteCloser
	Consumer ReadCloser
}

type URLData struct {
	URLShort    string `json:"url_short"`
	URLOriginal string `json:"url_original"`
}

// Repo - структура базы данных.
type Repo struct {
	URLOriginalAndShort []URLData
	FileStoragePath     string
	Storage             *Storage
}

// NewRepo - создание структуры Repo.
func NewRepo(fileStoragePath string, storage *Storage) *Repo {
	return &Repo{
		URLOriginalAndShort: make([]URLData, 0),
		FileStoragePath:     fileStoragePath,
		Storage:             storage,
	}
}

// Create - сохранение URL в базу данных.
func (repo *Repo) Create(ctx context.Context, urlOriginal, urlShort string) error {
	for _, urlData := range repo.URLOriginalAndShort {
		if urlData.URLOriginal == urlOriginal {
			return constants.ErrorURLAlreadyExist
		}
	}

	urlData := URLData{
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
func (repo *Repo) Get(ctx context.Context, urlShort string) (string, error) {
	for _, urlData := range repo.URLOriginalAndShort {
		if urlData.URLShort == urlShort {
			return urlData.URLOriginal, nil
		}
	}
	return "", constants.ErrorURLNotExist
}
