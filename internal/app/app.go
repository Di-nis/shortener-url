package app

import (
	"net/http"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/handler"
)

func Run() error {
	a := new(config.Run)
	a.ParseOptions()

	router := handler.CreateRouter()
	return http.ListenAndServe(a.URL, router)
}
