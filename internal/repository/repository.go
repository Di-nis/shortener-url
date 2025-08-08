package repository

import (
	"errors"
)

const urlShortMock = "hjkFSsrTWdf"

var (
	OriginalAndShotArray = map[string]string{
		"EwHXdJfB": "https://practicum.yandex.ru/",
	}
)

type Repo struct {}

func NewRepo() *Repo {
	return &Repo{}
}

// CreateURL - создание короткого адреса URL.
func (s *Repo) CreateURL(urlOriginalIn string) string {
	for urlShort, urlOriginal := range OriginalAndShotArray {
		if urlOriginal == urlOriginalIn {
			return urlShort
		}
	}
	OriginalAndShotArray[urlShortMock] = urlOriginalIn
	return urlShortMock
}

// GetURL - получение оригинального адреса URL.
func (s *Repo) GetURL(urlShort string) (string, error) {
	urlOriginal, ok := OriginalAndShotArray[urlShort]
	if !ok {
		return "", errors.New("internal/service/service.go: no data")
	}
	return urlOriginal, nil
}
