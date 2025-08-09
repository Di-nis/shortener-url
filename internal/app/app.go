package app

import (
	"fmt"
	"net/http"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/handler"
	// "github.com/Di-nis/shortener-url/internal/repository"
	// "github.com/Di-nis/shortener-url/internal/service"
)

func Run() error {
	options := new(config.Options)
	options.Parse()
	fmt.Println(options.Port)

	// repo := repository.NewRepo()
	// service := service.NewService(repo)
	// controller := handler.NewСontroller(service, options)

	router := handler.CreateRouter()
	return http.ListenAndServe(":" + options.Port, router)
}
