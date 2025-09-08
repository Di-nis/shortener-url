package repository

import (
	"context"

	"github.com/Di-nis/shortener-url/internal/constants"
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

// RepoFile - структура базы данных.
type RepoFile struct {
	URLOriginalAndShort []URLData
	FileStoragePath     string
	Storage             *Storage
}

// NewRepoFile - создание структуры Repo.
func NewRepoFile(fileStoragePath string, storage *Storage) *RepoFile {
	return &RepoFile{
		URLOriginalAndShort: make([]URLData, 0),
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
func (repo *RepoFile) Get(ctx context.Context, urlShort string) (string, error) {
	for _, urlData := range repo.URLOriginalAndShort {
		if urlData.URLShort == urlShort {
			return urlData.URLOriginal, nil
		}
	}
	return "", constants.ErrorURLNotExist
}
