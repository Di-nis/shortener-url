package handler

import (
	"fmt"
	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/service"
	"io"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
)

// Создание роутера.
func CreateRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/", createShortURL)
	router.Get("/{short_url}", getOriginalURL)
	return router
}

// createShortURL обрабатывает HTTP-запрос.
func createShortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, _ := io.ReadAll(req.Body)
	if reflect.DeepEqual(bodyBytes, []byte{}) {
		http.Error(res, "Не удалось прочитать тело запроса", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	urlOriginal := string(bodyBytes)
	urlShort := service.CreateURLShort(urlOriginal)
	bodyResult := fmt.Sprintf("http://localhost:%s/%s", config.Port, urlShort)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(bodyResult))
}

// getOriginalURL обрабатывает HTTP-запрос.
func getOriginalURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	URLShort := chi.URLParam(req, "short_url")
	defer req.Body.Close()

	urlOriginal, err := service.GetURLOriginal(URLShort)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.Header().Add("Location", urlOriginal)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}
