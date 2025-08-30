package repository

import (
	"errors"
)

type URLData struct {
	URLShort    string `json:"url_short"`
	URLOriginal string `json:"url_original"`
}

// Repo - структура базы данных.
type Repo struct {
	URLOriginalAndShort []URLData
	FileStoragePath     string
}

// NewRepo - создание структуры Repo.
func NewRepo(fileStoragePath string) *Repo {
	return &Repo{
		URLOriginalAndShort: make([]URLData, 100),
		FileStoragePath:     fileStoragePath,
	}
}

// Create - сохранение URL в базу данных.
func (repo *Repo) Create(urlOriginal, urlShort string) error {
	for _, urlData := range repo.URLOriginalAndShort {
		if urlData.URLOriginal == urlOriginal {
			return errors.New("internal/repository/repository.go: короткий URL уже существует")
		}
	}

	urlData := URLData{
		URLShort:    urlShort,
		URLOriginal: urlOriginal,
	}

	repo.URLOriginalAndShort = append(repo.URLOriginalAndShort, urlData)

	err := SaveToFile(repo.FileStoragePath, urlData)
	if err != nil {
		return err
	}
	return nil
}

// Get - получение оригинального URL из базы данных.
func (repo *Repo) Get(urlShort string) (string, error) {
	for _, urlData := range repo.URLOriginalAndShort {
		if urlData.URLShort == urlShort {
			return urlData.URLOriginal, nil
		}
	}
	return "", errors.New("internal/repository/repository.go: данные отсутствуют")
}
