package service

import (
	"errors"
)

var (
	OriginalAndShotArray = map[string]string{
		"EwHXdJfB":   "https://practicum.yandex.ru/",
		"WUonuiKQTU": "https://www.sports.ru/",
	}
)

// CreateURLShort - создание короткого адреса URL.
func CreateURLShort(urlOriginalIn string) (string, error) {
	for urlShort, urlOriginal := range OriginalAndShotArray {
		if urlOriginal == urlOriginalIn {
			return urlShort, nil
		}
	}
	return "", errors.New("internal/service/service.go: error create short url")
}

// GetURLOriginal - получение оригинального адреса URL.
func GetURLOriginal(urlShort string) (string, error) {
	urlOriginal, ok := OriginalAndShotArray[urlShort]
	if !ok {
		return "", errors.New("internal/service/service.go: no data")
	}
	return urlOriginal, nil
}
