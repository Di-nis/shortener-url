package repository

import (
	"errors"
)

// Repo - структура базы данных.
type Repo struct {
	urlOriginalAndShort map[string]string
}

// NewRepo - создание структуры Repo.
func NewRepo() *Repo {
	return &Repo{
		urlOriginalAndShort: make(map[string]string, 100),
	}
}

// Create - сохранение URL в базу данных.
func (repo *Repo) Create(urlOriginal, urlShort string) error {
	if _, ok := repo.urlOriginalAndShort[urlShort]; ok {
		return errors.New("internal/repository/repository.go: короткий URL уже существует")
	}
	repo.urlOriginalAndShort[urlShort] = urlOriginal
	return nil
}

// Get - получение оригинального URL из базы данных.
func (repo *Repo) Get(urlShort string) (string, error) {
	urlOriginal, ok := repo.urlOriginalAndShort[urlShort]
	if !ok {
		return "", errors.New("internal/repository/repository.go: данные отсутствуют")
	}
	return urlOriginal, nil
}
