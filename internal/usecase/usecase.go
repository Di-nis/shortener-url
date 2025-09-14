package usecase

import (
	"context"
	"errors"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/Di-nis/shortener-url/internal/service"
)

// URLRepository - интерфейс для базы данных.
type URLRepository interface {
	CreateBatch(context.Context, []models.URL) error
	CreateOrdinaty(context.Context, models.URL) error
	GetOriginalURL(context.Context, string) (string, error)
	GetShortURL(context.Context, string) (string, error)
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

// convertToSingleType - приведение к единому типу данных.
func convertToSingleType(urlIn any) models.URL {
	url1, ok1 := urlIn.(models.URL)
	if ok1 {
		return url1
	}
	url2, ok2 := urlIn.(models.URLCopyOne)
	if ok2 {
		return models.URL(url2)
	}
	return models.URL{}
}

// CreateURLOrdinary - создание короткого URL и его запись в базу данных.
func (urlUserCase *URLUseCase) CreateURLOrdinary(ctx context.Context, urlIn any) (models.URL, error) {
	urlOrdinary := convertToSingleType(urlIn)
	urlOrdinary.Short = urlUserCase.Service.ShortHash(urlOrdinary.Original, constants.HashLength)

	err := urlUserCase.Repo.CreateOrdinaty(ctx, urlOrdinary)

	if err != nil && errors.As(err, &constants.PgErr) {
		switch constants.PgErr.Code {
        case "23505":
            urlOrdinary.Short, _ = urlUserCase.Repo.GetShortURL(ctx, urlOrdinary.Original)
        }
		return urlOrdinary, err
    } else if err != nil && errors.Is(err, constants.ErrorURLAlreadyExist) {
		urlOrdinary.Short, _ = urlUserCase.Repo.GetShortURL(ctx, urlOrdinary.Original)
		return urlOrdinary, err
	}
	return urlOrdinary, nil
}

// CreateURLBatch - создание короткого URL и его запись в базу данных.
func (urlUserCase *URLUseCase) CreateURLBatch(ctx context.Context, urls []models.URL) ([]models.URL, error) {
	var idxTemp int

	for idx, url := range urls {
		urls[idx].Short = urlUserCase.Service.ShortHash(url.Original, constants.HashLength)

		if idx%1000 == 0 || idx == len(urls)-1 {
			urlsTemp := urls[idxTemp : idx+1]
			idxTemp = idx + 1

			err := urlUserCase.Repo.CreateBatch(ctx, urlsTemp)
			if err != nil {
				return nil, err
			}
		}
	}
	return urls, nil
}

// GetURL - получение оригинального URL.
func (urlUserCase *URLUseCase) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := urlUserCase.Repo.GetOriginalURL(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}
