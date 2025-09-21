package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/Di-nis/shortener-url/internal/service"
)

// URLRepository - интерфейс для базы данных.
type URLRepository interface {
	Ping(context.Context) error
	CreateBatch(context.Context, []models.URL) error
	CreateOrdinary(context.Context, models.URL) error
	GetOriginalURL(context.Context, string) (string, error)
	GetShortURL(context.Context, string) (string, error)
	GetAllURLs(context.Context, string) ([]models.URL, error)
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

// Ping - проверка соединения с базой данных.
func (urlUserCase *URLUseCase) Ping(ctx context.Context) error {
	return urlUserCase.Repo.Ping(ctx)
}

// CreateURLOrdinary - создание короткого URL и его запись в базу данных.
func (urlUserCase *URLUseCase) CreateURLOrdinary(ctx context.Context, urlIn any, baseURL string) (models.URL, error) {
	var PgErr *pgconn.PgError

	urlOrdinary := convertToSingleType(urlIn)
	urlOrdinary.Short = urlUserCase.Service.ShortHash(urlOrdinary.Original, constants.HashLength)

	err := urlUserCase.Repo.CreateOrdinary(ctx, urlOrdinary)

	fmt.Println(err)
	if err == nil {
		return urlOrdinary, nil
	}

	if errors.As(err, &PgErr) {
		switch PgErr.Code {
		case "23505":
			urlOrdinary.Short, _ = urlUserCase.Repo.GetShortURL(ctx, urlOrdinary.Original)
		}
		return urlOrdinary, err
	} else if errors.Is(err, constants.ErrorURLAlreadyExist) {
		urlOrdinary.Short, _ = urlUserCase.Repo.GetShortURL(ctx, urlOrdinary.Original)
		return urlOrdinary, err
	} else {
		return urlOrdinary, err
	}
}

// CreateURLBatch - создание короткого URL и его запись в базу данных.
func (urlUserCase *URLUseCase) CreateURLBatch(ctx context.Context, urls []models.URL, baseURL string) ([]models.URL, error) {
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

// GetAllURLs - получение всех когда-либо сокращенных пользователем URL.
func (urlUserCase *URLUseCase) GetAllURLs(ctx context.Context, userID string) ([]models.URL, error) {
	urls, err := urlUserCase.Repo.GetAllURLs(ctx, userID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}
