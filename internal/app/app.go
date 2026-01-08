// Package app. Инициализация зависимостей
// конфигурации, инфраструктурных компонентов и запуск серверв.
package app

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Di-nis/shortener-url/internal/logger"
	grpcServer "github.com/Di-nis/shortener-url/internal/server/grpc"
	httpServer "github.com/Di-nis/shortener-url/internal/server/http"
	"github.com/Di-nis/shortener-url/internal/service"

	"github.com/joho/godotenv"
)

// Run - запуск приложения.
func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	if err := godotenv.Load(); err != nil {
		logger.Sugar.Infof("path: internal/config/config.go, func loanFromEnv(), failed to load env: %w", err)
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

	// gRPC-сервер
	if cfg.EnableGRPC {
		return grpcServer.Run(ctx, cfg, repo, svc)
	}
	// HTTP-сервер
	return httpServer.Run(ctx, cfg, repo, svc)
}
