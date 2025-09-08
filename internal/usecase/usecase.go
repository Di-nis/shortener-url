package usecase

import (
	"context"
	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/service"
)

// URLRepository - интерфейс для базы данных.
type URLRepository interface {
	Create(context.Context, string, string) error
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

// CreateURL - создание короткого URL и его запись в базу данных.
func (urlUserCase *URLUseCase) CreateURL(ctx context.Context, originalURL string) (string, error) {
	shortURL := urlUserCase.Service.ShortHash(originalURL, constants.HashLength)

	err := urlUserCase.Repo.Create(ctx, originalURL, shortURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

// GetURL - получение оригинального URL.
func (urlUserCase *URLUseCase) GetURL(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := urlUserCase.Repo.Get(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}
