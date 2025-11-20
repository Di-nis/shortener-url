// Package app. Инициализация зависимостей
// конфигурации, инфраструктурных компонентов и запуск серверв.
package app

import (
	"net/http"

	"github.com/Di-nis/shortener-url/internal/service"
)

// Run - запуск приложения.
func Run() error {
	cfg, err := initConfigAndLogger()
	if err != nil {
		return err
	}
	repo, err := initStorage(cfg)
	if err != nil {
		return err
	}

	// defer repo.Close()

	svc := service.NewService()
	router := setupRouter(cfg, repo, svc)
	return http.ListenAndServe(cfg.ServerAddress, router)
}
