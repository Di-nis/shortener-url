// Package app. Инициализация зависимостей
// конфигурации, инфраструктурных компонентов и запуск серверв.
package app

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/service"

	"github.com/joho/godotenv"
)

// Run - запуск приложения.
func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("path: internal/config/config.go, func loanFromEnv(), failed to load env: %w", err)
	}

	cfg, err := initConfigAndLogger()
	if err != nil {
		return err
	}

	repo, err := initStorage(cfg)
	if err != nil {
		return err
	}

	svc := service.NewService()
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
			logger.Sugar.Errorf("Ошибка graceful shutdown: %w", err)
		}
		if err = repo.Close(); err != nil {
			logger.Sugar.Errorf("Ошибка закрытия базы данных: %w", err)
		}
	}()
	<-shutDownCtx.Done()
	return nil
}
