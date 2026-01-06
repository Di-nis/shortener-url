package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Di-nis/shortener-url/internal/constants"
)

// addBaseURLToResponse - добавление базового URL к ответу.
func addBaseURLToResponse(baseURL string, urlShort string) string {
	return fmt.Sprintf("%s/%s", baseURL, urlShort)
}

// writeStatusCreate - запись статус-кода в ответ для функций создания url.
func writeStatusCreate(res http.ResponseWriter, err error) {
	if err == nil {
		res.WriteHeader(http.StatusCreated)
	} else if errors.Is(err, constants.ErrorURLAlreadyExist) {
		res.WriteHeader(http.StatusConflict)
	} else {
		res.WriteHeader(http.StatusServiceUnavailable)
	}
}

// writeStatusCodePing - запись статус-кода в ответ для pingDB.
func writeStatusCodePing(res http.ResponseWriter, err error) {
	if err != nil {
		if errors.Is(err, constants.ErrorMethodNotAllowed) {
			res.WriteHeader(http.StatusMethodNotAllowed)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		res.WriteHeader(http.StatusOK)
	}
}
