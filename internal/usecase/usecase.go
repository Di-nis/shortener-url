package usecase

import (
	"context"
	// "fmt"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/Di-nis/shortener-url/internal/service"
)

// URLRepository - интерфейс для базы данных.
type URLRepository interface {
	Create(context.Context, []models.URL) error
	// Create(context.Context, string, string) error
	Get(context.Context, string) (string, error)
}

// URLUseCase - структура создания короткого и получение оригинального url.
type URLUseCase struct {
	Repo    URLRepository
	Service *service.Service
}

// NewURLUseCase - создание структуры URLUseCase.
func NewURLUseCase(repo URLRepository, service *service.Service) *URLUseCase {
	return &URLUseCase{
		Repo:    repo,
		Service: service,
	}
}

func MyFunc(urlsIn any) []models.URL {
	url1, ok1 := urlsIn.([]models.URL)
	if ok1 {
		return url1
	}

	url2, ok2 := urlsIn.(models.URL)
	if ok2 {
		return []models.URL{url2}
	}
	url3, ok3 := urlsIn.(models.URLCopyOne)
	if ok3 {
		urlsOut := models.URL(url3)
		return []models.URL{urlsOut}
	}
	return nil
}

// CreateURL - создание короткого URL и его запись в базу данных.
func (urlUserCase *URLUseCase) CreateURL(ctx context.Context, urlsIn any) ([]models.URL, error) {
	var (
		urls    = MyFunc(urlsIn)
		idxTemp int
	)

	for idx, url := range urls {
		urls[idx].Short = urlUserCase.Service.ShortHash(url.Original, constants.HashLength)

		if idx%1000 == 0 || idx == len(urls)-1 {
			urlsTemp := urls[idxTemp : idx+1]
			idxTemp = idx

			err := urlUserCase.Repo.Create(ctx, urlsTemp)
			if err != nil {
				return nil, err
			}
		}
	}
	return urls, nil
}

// GetURL - получение оригинального URL.
func (urlUserCase *URLUseCase) GetURL(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := urlUserCase.Repo.Get(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}
