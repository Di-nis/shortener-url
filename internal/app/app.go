package app

import (
	"net/http"

	"github.com/Di-nis/shortener-url/internal/compress"
	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/handler"
	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/repository"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/Di-nis/shortener-url/internal/usecase"

	"go.uber.org/zap"
)

func Run() error {
	config := &config.Config{}
	config.Parse()

	var err error
	if err = logger.Initialize(config.LogLevel); err != nil {
		return err
	}
	logger.Log.Info("Running server", zap.String("address", config.ServerAddress))

	repo := repository.NewRepo(config.FileStoragePath)

	// Подготовка данных
	consumer, err := repository.NewConsumer(config.FileStoragePath)
	if err != nil {
		return err
	}
	repo.URLOriginalAndShort, err = consumer.LoadFromFile()
	if err != nil {
		return err
	}

	svc := service.NewService()

	urlUseCase := usecase.NewURLUseCase(repo, svc)
	controller := handler.NewСontroller(urlUseCase, config)

	router := controller.CreateRouter()
	return http.ListenAndServe(config.ServerAddress, logger.WithLogging(compress.GzipMiddleware(router)))
}
