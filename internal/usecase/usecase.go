// Package usecase реализовывает связь между слоями handler, service и repository приложения.
package usecase

import (
	"context"
	"errors"
	"sync"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/Di-nis/shortener-url/internal/service"
)

// URLRepository - интерфейс для базы данных.
type URLRepository interface {
	Ping(context.Context) error
	InsertBatch(context.Context, []models.URLBase) error
	InsertOrdinary(context.Context, models.URLBase) error
	SelectOriginal(context.Context, string) (string, error)
	SelectShort(context.Context, string) (string, error)
	SelectAll(context.Context, string) ([]models.URLBase, error)
	Delete(context.Context, []models.URLBase) error
	GetCountURLs(context.Context) (int, error)
	GetCountUsers(context.Context) (int, error)
	Close() error
}

// convertToSingleType - приведение к единому типу данных.
func convertToSingleType(urlIn any) models.URLBase {
	url1, ok1 := urlIn.(models.URLBase)
	if ok1 {
		return url1
	}
	url2, ok2 := urlIn.(models.URLJSON)
	if ok2 {
		return models.URLBase(url2)
	}
	return models.URLBase{}
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
func (urlUseCase *URLUseCase) Ping(ctx context.Context) error {
	return urlUseCase.Repo.Ping(ctx)
}

// CreateURLOrdinary - создание короткого URL и его запись в базу данных.
func (urlUseCase *URLUseCase) CreateURLOrdinary(ctx context.Context, urlIn any) (models.URLBase, error) {
	urlOrdinary := convertToSingleType(urlIn)
	urlOrdinary.Short = urlUseCase.Service.ShortHash(urlOrdinary.Original, constants.HashLength)

	err := urlUseCase.Repo.InsertOrdinary(ctx, urlOrdinary)

	if err == nil {
		return urlOrdinary, nil
	}

	if errors.Is(err, constants.ErrorURLAlreadyExist) {
		urlOrdinary.Short, _ = urlUseCase.Repo.SelectShort(ctx, urlOrdinary.Original)
		return urlOrdinary, err
	} else {
		return urlOrdinary, err
	}
}

// CreateURLBatch - создание короткого URL и его запись в базу данных.
func (urlUseCase *URLUseCase) CreateURLBatch(ctx context.Context, urls []models.URLBase) ([]models.URLBase, error) {
	var idxTemp int

	for idx, url := range urls {
		urls[idx].Short = urlUseCase.Service.ShortHash(url.Original, constants.HashLength)

		if idx%1000 == 0 || idx == len(urls)-1 {
			urlsTemp := urls[idxTemp : idx+1]
			idxTemp = idx + 1

			err := urlUseCase.Repo.InsertBatch(ctx, urlsTemp)
			if err != nil {
				return nil, err
			}
		}
	}
	return urls, nil
}

// GetOriginalURL - получение оригинального URL.
func (urlUseCase *URLUseCase) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := urlUseCase.Repo.SelectOriginal(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

// GetAllURLs - получение всех когда-либо сокращенных пользователем URL.
func (urlUseCase *URLUseCase) GetAllURLs(ctx context.Context, userID string) ([]models.URLBase, error) {
	urls, err := urlUseCase.Repo.SelectAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// generator - генерирует сообщения в канал.
func (urlUseCase *URLUseCase) generator(ctx context.Context, urls []models.URLBase, inChan chan models.URLBase) {
	for _, url := range urls {
		select {
		case <-ctx.Done():
			return
		case inChan <- url:
		}
	}
}

// worker - работник.
func (urlUseCase *URLUseCase) worker(ctx context.Context, urls <-chan models.URLBase, result chan error) {
	urlsToDB := make([]models.URLBase, 0, 100)

	for {
		select {
		case <-ctx.Done():
			if len(urlsToDB) > 0 {
				result <- urlUseCase.Repo.Delete(ctx, urlsToDB)
			}
			return

		case url, ok := <-urls:
			if !ok {
				if len(urlsToDB) > 0 {
					result <- urlUseCase.Repo.Delete(ctx, urlsToDB)
				}
				return
			}
			urlsToDB = append(urlsToDB, url)
			if len(urlsToDB) >= 1 {
				result <- urlUseCase.Repo.Delete(ctx, urlsToDB)
				urlsToDB = urlsToDB[:0]
			}

		}
	}

}

// DeleteURLs - удаление сокращенных URL.
func (urlUseCase *URLUseCase) DeleteURLs(ctx context.Context, urls []models.URLBase) error {
	const numWorkers = 3
	inChan := make(chan models.URLBase, 1024)
	resultChan := make(chan error, numWorkers)
	var wg sync.WaitGroup

	go func() {
		defer close(inChan)
		urlUseCase.generator(ctx, urls, inChan)
	}()

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			urlUseCase.worker(ctx, inChan, resultChan)
		}(w)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var firstErr error
	for err := range resultChan {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr

}

// GetStats - получение статистики по записям и пользователям.
func (urlUseCase *URLUseCase) GetStats(ctx context.Context) (int, int, error) {
	countURLs, err := urlUseCase.Repo.GetCountURLs(ctx)
	if err != nil {
		return 0, 0, err
	}

	countUsers, err := urlUseCase.Repo.GetCountUsers(ctx)
	if err != nil {
		return 0, 0, err
	}

	return countURLs, countUsers, nil
}
