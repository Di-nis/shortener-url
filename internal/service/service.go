package service

import (
	"errors"
)


const urlShortMock = "EwHXdJfB"
var (
	OriginalAndShotArray = map[string]string{
		urlShortMock: "https://practicum.yandex.ru/",
	}
)

// CreateURLShort - создание короткого адреса URL.
func CreateURLShort(urlOriginalIn string) string {
	OriginalAndShotArray[urlOriginalIn] = urlShortMock
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
