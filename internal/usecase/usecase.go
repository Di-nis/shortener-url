package usecase

import (
	"context"
	"errors"
	"fmt"
	// "sync"

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

// func (urlUseCase *URLUseCase) worker(ctx context.Context, jobs <-chan Job, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case job, ok := <-jobs:
// 			if !ok {
// 				return
// 			}
// 			_ = urlUseCase.Repo.DeleteURL(ctx, job.urls)

// 		}
// 	}
// }
// type Job struct {
// 	urls []models.URL
// }

func (urlUseCase *URLUseCase) DeleteURLs(ctx context.Context, urls []models.URL) error {
	// ctx, cancel := context.WithCancel(ctx)
	// defer cancel()


	doneCh := make(chan struct{})
    // закрываем его при завершении программы
    defer close(doneCh)

    // канал с данными
    inputCh := urlUseCase.generator(doneCh, urls)

    // получаем слайс каналов из 10 рабочих add
    channels := urlUseCase.fanOut(ctx, doneCh, inputCh)

    // а теперь объединяем десять каналов в один
    resultCh := urlUseCase.fanIn(doneCh, channels...)

    // передаём тот один канал в следующий этап обработки
    // resultCh := multiply(doneCh, addResultCh)

    // выводим результаты расчетов из канала
    for res := range resultCh {
        fmt.Println(res)
    }
	return nil


	// const numJobs = 1000
	// // создаем буферизованный канал для принятия задач в воркер
	// // jobs := make(chan int, numJobs)
	// // создаем буферизованный канал для отправки результатов
	// results := make(chan error, numJobs)
	// defer close(results)
	// // doneCh := make(chan struct{})
	// // закрываем его при завершении программы
	// // defer close(doneCh)

	// // канал с данными
	// inputCh := urlUseCase.generator(urls, numJobs)

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for w := 1; w <= 3; w++ {
	// 		urlUseCase.worker(ctx, inputCh, results)
	// 	}
	// }()

	// // defer close(inputCh)

	// var resultWg sync.WaitGroup
	// errs := make([]error, 0)
	// resultWg.Add(1)
	// go func() {
	// 	defer resultWg.Done()
	// 	for r := range results {
	// 		errs = append(errs, r)
	// 	}
	// }()

	// go func() {
	// 	// Ждём воркеров
	// 	wg.Wait()
	// 	fmt.Println("Привет 1")
	// 	// Закрываем канал результатов и ждём агрегатор
	// 	resultWg.Wait()

	// 	fmt.Println("Привет 2")
	// }()
	// close(results)
	// fmt.Println("Привет 3")

	// // for res := range results {
    // //     fmt.Println(res)
    // // }

	// return errs[0]
}
