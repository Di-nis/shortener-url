package app

import (
	"fmt"
	"net/http"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/handler"
	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/repository"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/Di-nis/shortener-url/internal/usecase"

	"go.uber.org/zap"
)

// initConfigAndLogger - инициализация конфигурации и логгера.
func initConfigAndLogger() (*config.Config, error) {
	config := &config.Config{}
	config.Parse()

	var err error
	if err = logger.Initialize(config.LogLevel); err != nil {
		return nil, err
	}
	logger.Log.Info("Запуска сервера", zap.String("address", config.ServerAddress))
	return config, nil
}

// initRepoPostgres - инициализация репозитория для работы с PostgreSQL.
func initRepoPostgres(cfg *config.Config) (*repository.RepoPostgres, error) {
	repo, err := repository.NewRepoPostgres(cfg.DataBaseDSN)
	if err != nil {
		return nil, err
	}

	// выполнение миграций
	err = repo.Migrations()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// initRepoFile - инициализация репозитория для работы с файлом.
func initRepoFile(cfg *config.Config) (*repository.RepoFile, error) {
	consumer, err := repository.NewConsumer(cfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации consumer: %w", err)
	}

	producer, err := repository.NewProducer(cfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации producer: %w", err)
	}

	storage := &repository.Storage{
		Consumer: consumer,
		Producer: producer,
	}

	repo := repository.NewRepoFile(cfg.FileStoragePath, storage)
	repo.OriginalAndShortURL, err = storage.Consumer.LoadFromFile()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки данных из файла-хранилища: %w", err)
	}

	return repo, nil
}

// initStorage - инициализация хранилища данных.
func initStorage(cfg *config.Config) (usecase.URLRepository, error) {
	if cfg.DataBaseDSN != "" {
		return initRepoPostgres(cfg)
	}
	return initRepoFile(cfg)
}

// setupRouter - настройка маршрутизатора.
func setupRouter(cfg *config.Config, repo usecase.URLRepository, svc *service.Service) http.Handler {
	urlUseCase := usecase.NewURLUseCase(repo, svc)

	controller := handler.NewСontroller(urlUseCase, cfg)
	return controller.CreateRouter()
}
