package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/handler"
	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/Di-nis/shortener-url/internal/usecase"
)

// setupRouter - настройка маршрутизатора.
func setupRouter(cfg *config.Config, repo usecase.URLRepository, svc *service.Service) http.Handler {
	urlUseCase := usecase.NewURLUseCase(repo, svc)

	controller := handler.NewСontroller(urlUseCase, cfg)
	return controller.SetupRouter()
}

// Run - запуск HTTP-сервера.
func Run(ctx context.Context, cfg *config.Config, repo usecase.URLRepository, svc *service.Service) error {
	var err error
	routerHandler := setupRouter(cfg, repo, svc)

	httpServer := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: routerHandler,
	}

	go func() {
		if cfg.EnableHTTPS {
			if err = httpServer.ListenAndServeTLS(cfg.CertFilePath, cfg.KeyFilePath); err != nil && err != http.ErrServerClosed {
				logger.Sugar.Fatalf("failed start TLS-server: %w", err)
			}
		}
		if err = httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Sugar.Fatalf("failed start server: %w", err)
		}
	}()

	<-ctx.Done()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	go func() {
		if err := httpServer.Shutdown(context.Background()); err != nil {
			logger.Sugar.Errorf("failed graceful shutdown: %w", err)
		}
		if err = repo.Close(); err != nil {
			logger.Sugar.Errorf("failed closing database: %w", err)
		}
	}()
	<-shutDownCtx.Done()
	return nil
}
