package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/jackc/pgx/v5/pgconn"
)

// addBaseURLToResponse - добавление базового URL к ответу.
func addBaseURLToResponse(baseURL string, urlShort string) string {
	return fmt.Sprintf("%s/%s", baseURL, urlShort)

}

// getStatusCode - определение статус-кода ответа.
func getStatusCode(res http.ResponseWriter, err error) {
	var PgErr *pgconn.PgError
	if err != nil && errors.As(err, &PgErr) {
		switch PgErr.Code {
		case "23505":
			res.WriteHeader(http.StatusConflict)
		}
	} else if err != nil && errors.Is(err, constants.ErrorURLAlreadyExist) {
		res.WriteHeader(http.StatusConflict)
	} else if err != nil {
		res.WriteHeader(http.StatusServiceUnavailable)
	} else {
		res.WriteHeader(http.StatusCreated)
	}
}
