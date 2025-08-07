package handler

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"reflect"
)

var (
	ShortenerArray = make(map[string]string)
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

	// проверка на уже существование в БД
	// дополниеть реализацией связи с БД
	ShortenerArray["EwHXdJfB"] = string(bodyBytes)

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

	shortenerArray := make(map[string]string)
	shortenerArray["EwHXdJfB"] = "https://practicum.yandex.ru/"

	shortURL := chi.URLParam(req, "short_url")
	defer req.Body.Close()

	headerLocation, ok := shortenerArray[shortURL]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.Header().Add("Location", headerLocation)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}
