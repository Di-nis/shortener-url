package service

import (
	"errors"
)


const urlShortMock = "hjkFSsrTWdf"
var (
	OriginalAndShotArray = map[string]string{
		"EwHXdJfB": "https://practicum.yandex.ru/",
	}
)

// CreateURLShort - создание короткого адреса URL.
func CreateURLShort(urlOriginalIn string) string {
	for urlShort, urlOriginal := range OriginalAndShotArray {
		if urlOriginal == urlOriginalIn {
			return urlShort
		}
	}
	OriginalAndShotArray[urlShortMock] = urlOriginalIn
	return urlShortMock
}

// GetURLOriginal - получение оригинального адреса URL.
func GetURLOriginal(urlShort string) (string, error) {
	urlOriginal, ok := OriginalAndShotArray[urlShort]
	if !ok {
		return "", errors.New("internal/service/service.go: no data")
	}
	return urlOriginal, nil
}
