package usecase

import (
	"github.com/Di-nis/shortener-url/internal/repository"
	"github.com/Di-nis/shortener-url/internal/service"
)

const hashLength = 8

type URLRepository interface {
	Create(string, string) error
	Get(string) (string, error)
}

// URLUseCase - структура создания короткого и получение оригинального url.
type URLUseCase struct {
	Repo    URLRepository
	Service *service.Service
}

// NewURLUseCase - создание структуры URLUseCase.
func NewURLUseCase(repo *repository.Repo, service *service.Service) *URLUseCase {
	return &URLUseCase{
		Repo:    repo,
		Service: service,
	}
}

// CreateURL - создание короткого URL и его запись в базу данных.
func (urlUserCase *URLUseCase) CreateURL(originalURL string) (string, error) {
	shortURL := urlUserCase.Service.ShortHash(originalURL, hashLength)

	err := urlUserCase.Repo.Create(originalURL, shortURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

// GetURL - получение оригинального URL.
func (urlUserCase *URLUseCase) GetURL(shortURL string) (string, error) {
	originalURL, err := urlUserCase.Repo.Get(shortURL)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}
