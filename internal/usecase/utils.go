package usecase

import (
	"context"
	"sync"

	"github.com/Di-nis/shortener-url/internal/models"
)

// generator функция из предыдущего примера, делает то же, что и делала
func (urlUseCase *URLUseCase) generator(ctx context.Context, input [][]models.URL) chan []models.URL {
	inputCh := make(chan []models.URL)

	go func() {
		defer close(inputCh)
		for _, data := range input {
			select {
			case <-ctx.Done():
				return
			case inputCh <- data:
			}
		}
	}()

	return inputCh
}

// add функция из предыдущего примера, делает то же, что и делала
func (urlUseCase *URLUseCase) add(ctx context.Context, inputCh chan []models.URL) chan error {
	addRes := make(chan error)

	go func() {
		defer close(addRes)

		for data := range inputCh {
			result := urlUseCase.Repo.DeleteURL(ctx, data)

			select {
			case <-ctx.Done():
				return
			case addRes <- result:
			}
		}
	}()
	return addRes
}

// fanOut принимает канал данных, порождает 10 горутин
func (urlUseCase *URLUseCase) fanOut(ctx context.Context, inputCh chan []models.URL) []chan error {
	numWorkers := 50
	channels := make([]chan error, numWorkers)
	wg := sync.WaitGroup{}

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(inputCh chan []models.URL, wg *sync.WaitGroup) {
			defer wg.Done()
			addResultCh := urlUseCase.add(ctx, inputCh)
			channels[i] = addResultCh
		}(inputCh, &wg)
	}

	wg.Wait()

	return channels
}

// fanIn объединяет несколько каналов resultChs в один.
func (urlUseCase *URLUseCase) fanIn(ctx context.Context, resultChs ...chan error) chan error {
	finalCh := make(chan error)

	var wg sync.WaitGroup

	for _, ch := range resultChs {
		chClosure := ch
		wg.Add(1)

		go func() {
			defer wg.Done()

			for data := range chClosure {
				select {
				case <-ctx.Done():
					return
				case finalCh <- data:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}
