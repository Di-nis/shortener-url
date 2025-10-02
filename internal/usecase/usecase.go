package usecase

import (
	"context"
	"errors"
	"sync"

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
	DeleteURL(context.Context, []models.URL) error
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
func (urlUseCase *URLUseCase) Ping(ctx context.Context) error {
	return urlUseCase.Repo.Ping(ctx)
}

// CreateURLOrdinary - создание короткого URL и его запись в базу данных.
func (urlUseCase *URLUseCase) CreateURLOrdinary(ctx context.Context, urlIn any, baseURL string) (models.URL, error) {
	var PgErr *pgconn.PgError

	urlOrdinary := convertToSingleType(urlIn)
	urlOrdinary.Short = urlUseCase.Service.ShortHash(urlOrdinary.Original, constants.HashLength)

	err := urlUseCase.Repo.CreateOrdinary(ctx, urlOrdinary)

	if err == nil {
		return urlOrdinary, nil
	}

	if errors.As(err, &PgErr) {
		switch PgErr.Code {
		case "23505":
			urlOrdinary.Short, _ = urlUseCase.Repo.GetShortURL(ctx, urlOrdinary.Original)
		}
		return urlOrdinary, err
	} else if errors.Is(err, constants.ErrorURLAlreadyExist) {
		urlOrdinary.Short, _ = urlUseCase.Repo.GetShortURL(ctx, urlOrdinary.Original)
		return urlOrdinary, err
	} else {
		return urlOrdinary, err
	}
}

// CreateURLBatch - создание короткого URL и его запись в базу данных.
func (urlUseCase *URLUseCase) CreateURLBatch(ctx context.Context, urls []models.URL, baseURL string) ([]models.URL, error) {
	var idxTemp int

	for idx, url := range urls {
		urls[idx].Short = urlUseCase.Service.ShortHash(url.Original, constants.HashLength)

		if idx%1000 == 0 || idx == len(urls)-1 {
			urlsTemp := urls[idxTemp : idx+1]
			idxTemp = idx + 1

			err := urlUseCase.Repo.CreateBatch(ctx, urlsTemp)
			if err != nil {
				return nil, err
			}
		}
	}
	return urls, nil
}

// GetURL - получение оригинального URL.
func (urlUseCase *URLUseCase) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := urlUseCase.Repo.GetOriginalURL(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

// GetAllURLs - получение всех когда-либо сокращенных пользователем URL.
func (urlUseCase *URLUseCase) GetAllURLs(ctx context.Context, userID string) ([]models.URL, error) {
	urls, err := urlUseCase.Repo.GetAllURLs(ctx, userID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// generator - генерирует сообщения в канал.
func (urlUseCase *URLUseCase) generator(ctx context.Context, urls []models.URL, inChan chan models.URL) {
	for _, url := range urls {
		select {
		case <-ctx.Done():
			return
		case inChan <- url:
		}
	}
}

func (urlUseCase *URLUseCase) worker(ctx context.Context, urls <-chan models.URL, result chan error) {
	urlsToDB := make([]models.URL, 0, 100)

	for {
		select {
		case <-ctx.Done():
			if len(urlsToDB) > 0 {
				result <- urlUseCase.Repo.DeleteURL(ctx, urlsToDB)
			}
			return

		case url, ok := <-urls:
			if !ok {
				if len(urlsToDB) > 0 {
					result <- urlUseCase.Repo.DeleteURL(ctx, urlsToDB)
				}
				return
			}
			urlsToDB = append(urlsToDB, url)
			if len(urlsToDB) >= 1 {
				result <- urlUseCase.Repo.DeleteURL(ctx, urlsToDB)
				urlsToDB = urlsToDB[:0]
			}

		}
	}

}

// DeleteURLs - удаление сокращенных URL.
func (urlUseCase *URLUseCase) DeleteURLs(ctx context.Context, urls []models.URL) error {
	const numWorkers = 3
	inChan := make(chan models.URL, 1024)
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
	return nil

}
