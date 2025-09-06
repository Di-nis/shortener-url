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

func initStorageAndServices(cfg *config.Config) (*repository.Repo, *service.Service, error) {
	consumer, err := repository.NewConsumer(cfg.FileStoragePath)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка инициализации consumer: %w", err)
	}

	producer, err := repository.NewProducer(cfg.FileStoragePath)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка инициализации producer: %w", err)
	}

	storage := &repository.Storage{
		Consumer: consumer,
		Producer: producer,
	}

	repo := repository.NewRepo(cfg.FileStoragePath, storage)
	repo.URLOriginalAndShort, err = storage.Consumer.LoadFromFile()
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка загрузки данных из файла-хранилища: %w", err)
	}

	svc := service.NewService()
	return repo, svc, nil
}

func setupRouter(cfg *config.Config, repo *repository.Repo, svc *service.Service) http.Handler {
	urlUseCase := usecase.NewURLUseCase(repo, svc)
	controller := handler.NewСontroller(urlUseCase, cfg)
	return controller.CreateRouter()
}
