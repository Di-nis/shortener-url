package app

import (
	"fmt"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/Di-nis/shortener-url/internal/repository"
	"github.com/Di-nis/shortener-url/internal/storage"
	"github.com/Di-nis/shortener-url/internal/usecase"

	"go.uber.org/zap"
)

// initConfigAndLogger - инициализация конфигурации и логгера.
func initConfigAndLogger() (*config.Config, error) {
	cfg := config.NewConfig()
	cfg.Load()

	var err error
	if err = logger.Initialize(cfg.LogLevel); err != nil {
		return nil, err
	}
	logger.Sugar.Info("Запуска сервера", zap.String("address", cfg.ServerAddress))
	return cfg, nil
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

// InitRepoFile - инициализация репозитория для работы с файлом.
func InitRepoFile(fileStoragePath string) (*repository.RepoFileMemory, error) {
	consumer, err := storage.NewConsumer(fileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации consumer: %w", err)
	}

	producer, err := storage.NewProducer(fileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации producer: %w", err)
	}

	storage := &repository.Storage{
		Consumer: consumer,
		Producer: producer,
	}

	repo := repository.NewRepoFileMemory(storage)
	repo.URLs, err = storage.Consumer.Load()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки данных из файла-хранилища: %w", err)
	}

	return repo, nil
}

// InitRepoMemory - инициализация репозитория для работы с памятью.
func InitRepoMemory(cfg *config.Config) (*repository.RepoFileMemory, error) {
	urls := make([]models.URLBase, 0)
	consumer := storage.NewConsumerMemory(urls)
	producer := storage.NewProducerMemory(urls)

	storage := &repository.Storage{
		Consumer: consumer,
		Producer: producer,
	}

	var err error
	repo := repository.NewRepoFileMemory(storage)
	repo.URLs, err = storage.Consumer.Load()
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
	if cfg.FileStoragePath != "" {
		return InitRepoFile(cfg.FileStoragePath)
	}
	return InitRepoMemory(cfg)
}
