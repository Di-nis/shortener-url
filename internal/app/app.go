package app

import (
	"github.com/Di-nis/shortener-url/internal/handler"
	"net/http"
)

func Run() error {
	router := handler.CreateRouter()
	return http.ListenAndServe(":8080", router)
}
