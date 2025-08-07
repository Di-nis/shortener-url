package handler

import (
	"io"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
)

var (
	OriginalAndShotArray = map[string]string{}
)

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

	OriginalAndShotArray["EwHXdJfB"] = string(bodyBytes)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("http://localhost:8080/EwHXdJfB"))
}

// getOriginalURL обрабатывает HTTP-запрос.
func getOriginalURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	shortURL := chi.URLParam(req, "short_url")
	defer req.Body.Close()

	headerLocation, ok := OriginalAndShotArray[shortURL]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.Header().Add("Location", headerLocation)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}
