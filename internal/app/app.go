package app

import (
	"net/http"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/handler"
	"github.com/Di-nis/shortener-url/internal/repository"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/Di-nis/shortener-url/internal/usecase"
)

func Run() error {
	config := new(config.Config)
	config.Parse()

	repo := repository.NewRepo()
	svc := service.NewService()

	urlUseCase := usecase.NewURLUseCase(repo, svc)
	controller := handler.NewСontroller(urlUseCase, config)

	router := controller.CreateRouter()
	return http.ListenAndServe(config.ServerAddress, router)
}
