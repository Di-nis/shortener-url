package app

import (
	"net/http"

	"github.com/Di-nis/shortener-url/internal/compress"
	"github.com/Di-nis/shortener-url/internal/logger"
)

func Run() error {
	cfg, err := initConfigAndLogger()
	if err != nil {
		return err
	}
	repo, svc, err := initStorageAndServices(cfg)
	if err != nil {
		return err
	}
	router := setupRouter(cfg, repo, svc)
	return http.ListenAndServe(cfg.ServerAddress, logger.WithLogging(compress.GzipMiddleware(router)))
}
