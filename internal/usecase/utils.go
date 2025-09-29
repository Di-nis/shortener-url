package usecase

import (
	// "context"
	// "fmt"
	"context"
	"sync"
	// "time"

	"github.com/Di-nis/shortener-url/internal/models"
)

// generator функция из предыдущего примера, делает то же, что и делала
func (urlUseCase *URLUseCase) generator(doneCh chan struct{}, input []models.URL) chan []models.URL {
    inputCh := make(chan []models.URL)
	batchSize := 1000
	temp := []models.URL{}

    go func() {
        defer close(inputCh)

		for i := 0; i < len(input); i += batchSize {
			end := i + batchSize
			end = min(end, len(input))
			temp = input[i:end]

			// inputCh <- temp
		// }

			// for _, data := range input {
			select {
			case <-doneCh:
				return
			case inputCh <- temp:
				}
			}
			// }
    }()

    return inputCh
}

// multiply функция из предыдущего примера, делает то же, что и делала
// func multiply(doneCh chan struct{}, inputCh chan int) chan int {
//     multiplyRes := make(chan int)

//     go func() {
//         defer close(multiplyRes)

//         for data := range inputCh {
//             result := data * 2

//             select {
//             case <-doneCh:
//                 return
//             case multiplyRes <- result:
//             }
//         }
//     }()
//     return multiplyRes
// }

// add функция из предыдущего примера, делает то же, что и делала
func (urlUseCase *URLUseCase) add(ctx context.Context, doneCh chan struct{}, inputCh chan []models.URL) chan error {
    addRes := make(chan error)
	// batchSize := 1000

    go func() {
        defer close(addRes)

        for data := range inputCh {
            // замедлим вычисление, как будто функция add требует больше вычислительных ресурсов
            // time.Sleep(time.Second)

            // тут вызов запроса в БД
			// for i := 0; i < len(urls); i += batchSize {
			// 	end := i + batchSize
			// 	end = min(end, len(urls))
			// 	temp = urls[i:end]
			// }

			result := urlUseCase.Repo.DeleteURL(ctx, data)
            // result := data + 1

            select {
            case <-doneCh:
                return
            case addRes <- result:
            }
        }
    }()
    return addRes
}


// fanOut принимает канал данных, порождает 10 горутин
func (urlUseCase *URLUseCase) fanOut(ctx context.Context, doneCh chan struct{}, inputCh chan []models.URL) []chan error {
    // количество горутин add
    numWorkers := 10
    // каналы, в которые отправляются результаты
    channels := make([]chan error, numWorkers)


    for i := 0; i < numWorkers; i++ {
        // получаем канал из горутины add
        addResultCh := urlUseCase.add(ctx, doneCh, inputCh)
        // отправляем его в слайс каналов
        channels[i] = addResultCh
    }

    // возвращаем слайс каналов
    return channels
}

// fanIn объединяет несколько каналов resultChs в один.
func (urlUseCase *URLUseCase) fanIn(doneCh chan struct{}, resultChs ...chan error) chan error {
    // конечный выходной канал в который отправляем данные из всех каналов из слайса, назовём его результирующим
    finalCh := make(chan error)

    // понадобится для ожидания всех горутин
    var wg sync.WaitGroup

    // перебираем все входящие каналы
    for _, ch := range resultChs {
        // в горутину передавать переменную цикла нельзя, поэтому делаем так 
        chClosure := ch

        // инкрементируем счётчик горутин, которые нужно подождать
        wg.Add(1)

        go func() {
            // откладываем сообщение о том, что горутина завершилась
            defer wg.Done()

            // получаем данные из канала
            for data := range chClosure {
                select {
                // выходим из горутины, если канал закрылся
                case <-doneCh:
                    return
                // если не закрылся, отправляем данные в конечный выходной канал
                case finalCh <- data:
                }
            }
        }()
    }

    go func() {
        // ждём завершения всех горутин
        wg.Wait()
        // когда все горутины завершились, закрываем результирующий канал
        close(finalCh)
    }()

    // возвращаем результирующий канал
    return finalCh
}

// func main() {
//     // слайс данных
//     input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

//     // сигнальный канал для завершения горутин
//     doneCh := make(chan struct{})
//     // закрываем его при завершении программы
//     defer close(doneCh)

//     // канал с данными
//     inputCh := generator(doneCh, input)

//     // получаем слайс каналов из 10 рабочих add
//     channels := fanOut(doneCh, inputCh)

//     // а теперь объединяем десять каналов в один
//     addResultCh := fanIn(doneCh, channels...)

//     // передаём тот один канал в следующий этап обработки
//     resultCh := multiply(doneCh, addResultCh)

//     // выводим результаты расчетов из канала
//     for res := range resultCh {
//         fmt.Println(res)
//     }
// }




// func (urlUseCase *URLUseCase) worker (ctx context.Context, inputCh <-chan []models.URL, results chan<- error) {
// 	// var m sync.RWMutex
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			// log.Printf("Worker %d завершает работу", id)
// 			return
// 		case job, ok := <-inputCh:
// 			if !ok {
// 				return
// 			}

// 			// m.Lock()
// 			err := urlUseCase.Repo.DeleteURL(ctx, job)
// 			fmt.Println(err)
// 			// results <- err
// 			// m.Unlock()
// 		}
// 	}
// }

// // generator функция из предыдущего примера, делает то же, что и делала
// func (urlUseCase *URLUseCase) generator(urls []models.URL, numJobs int) chan []models.URL {
// 	batchSize := 1000
// 	inputCh := make(chan []models.URL, numJobs)
// 	temp := make([]models.URL, batchSize)

// 	go func() {
// 		defer close(inputCh)

// 		for i := 0; i < len(urls); i += batchSize {
// 			end := i + batchSize
// 			end = min(end, len(urls))
// 			temp = urls[i:end]

// 			inputCh <- temp
// 		}


// 		// select {
// 		// case <-ctx.Done():
// 		// 	return
// 		// case inputCh <- temp:
// 		// }
// 	}()

// 	return inputCh
// }

// multiply функция из предыдущего примера, делает то же, что и делала
// func multiply(doneCh chan struct{}, inputCh chan int) chan int {
//     multiplyRes := make(chan int)

//     go func() {
//         defer close(multiplyRes)

//         for data := range inputCh {
//             result := data * 2

//             select {
//             case <-doneCh:
//                 return
//             case multiplyRes <- result:
//             }
//         }
//     }()
//     return multiplyRes
// }

// add функция из предыдущего примера, делает то же, что и делала
// func add(doneCh chan struct{}, inputCh chan int) chan int {
// 	addRes := make(chan int)

// 	go func() {
// 		defer close(addRes)

// 		for data := range inputCh {
// 			// замедлим вычисление, как будто функция add требует больше вычислительных ресурсов
// 			time.Sleep(time.Second)

// 			result := data + 1

// 			select {
// 			case <-doneCh:
// 				return
// 			case addRes <- result:
// 			}
// 		}
// 	}()
// 	return addRes
// }

// fanOut принимает канал данных, порождает 10 горутин
// func fanOut(doneCh chan struct{}, inputCh chan models.URL) []chan int {
//     // количество горутин add
//     numWorkers := 10
//     // каналы, в которые отправляются результаты
//     channels := make([]chan int, numWorkers)

//     for i := 0; i < numWorkers; i++ {
//         // получаем канал из горутины add
//         addResultCh := add(doneCh, inputCh)
//         // отправляем его в слайс каналов
//         channels[i] = addResultCh
//     }

//     // возвращаем слайс каналов
//     return channels
// }

// fanIn объединяет несколько каналов resultChs в один.
// func fanIn(doneCh chan struct{}, resultChs ...chan int) chan int {
// 	// конечный выходной канал в который отправляем данные из всех каналов из слайса, назовём его результирующим
// 	finalCh := make(chan int)

// 	// понадобится для ожидания всех горутин
// 	var wg sync.WaitGroup

// 	// перебираем все входящие каналы
// 	for _, ch := range resultChs {
// 		// в горутину передавать переменную цикла нельзя, поэтому делаем так
// 		chClosure := ch

// 		// инкрементируем счётчик горутин, которые нужно подождать
// 		wg.Add(1)

// 		go func() {
// 			// откладываем сообщение о том, что горутина завершилась
// 			defer wg.Done()

// 			// получаем данные из канала
// 			for data := range chClosure {
// 				select {
// 				// выходим из горутины, если канал закрылся
// 				case <-doneCh:
// 					return
// 				// если не закрылся, отправляем данные в конечный выходной канал
// 				case finalCh <- data:
// 				}
// 			}
// 		}()
// 	}

// 	go func() {
// 		// ждём завершения всех горутин
// 		wg.Wait()
// 		// когда все горутины завершились, закрываем результирующий канал
// 		close(finalCh)
// 	}()

// 	// возвращаем результирующий канал
// 	return finalCh
// }
