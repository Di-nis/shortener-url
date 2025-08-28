package app

import (
	"net/http"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/handler"
	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/repository"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/Di-nis/shortener-url/internal/usecase"
	"github.com/Di-nis/shortener-url/internal/compress"

	"go.uber.org/zap"
)

func Run() error {
	config := &config.Config{}
	config.Parse()

	if err := logger.Initialize(config.LogLevel); err != nil {
		return err
	}
	logger.Log.Info("Running server", zap.String("address", config.ServerAddress))

	repo := repository.NewRepo()
	svc := service.NewService()

	urlUseCase := usecase.NewURLUseCase(repo, svc)
	controller := handler.NewСontroller(urlUseCase, config)

	router := controller.CreateRouter()
	return http.ListenAndServe(config.ServerAddress, logger.WithLogging(compress.GzipMiddleware(router)))
}
