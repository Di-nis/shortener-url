package app

import (
	"fmt"
	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/handler"
	"net/http"
)

func Run() error {
	address := fmt.Sprintf(":%s", config.Port)

	router := handler.CreateRouter()
	return http.ListenAndServe(address, router)
}
