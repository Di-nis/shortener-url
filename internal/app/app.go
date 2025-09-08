package app

import (
	"net/http"
	// "context"

	"github.com/Di-nis/shortener-url/internal/compress"
	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/service"
)

func Run() error {
	// ctx := context.Background()
	cfg, err := initConfigAndLogger()
	if err != nil {
		return err
	}
	// написать реализацию выбора подключения к базе данных, postgres или хранение в файле
	repo, err := initStorage(cfg)
	if err != nil {
		return err
	}
	svc := service.NewService()
	router := setupRouter(cfg, repo, svc)
	return http.ListenAndServe(cfg.ServerAddress, logger.WithLogging(compress.GzipMiddleware(router)))
}
