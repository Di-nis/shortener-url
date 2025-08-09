package handler

import (
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/Di-nis/shortener-url/internal/repository"

	"github.com/go-chi/chi/v5"
)

type Controller struct {
	Service *service.Service
	Options *config.Options
}

func NewСontroller(service *service.Service, options *config.Options) *Controller {
	return &Controller{
		Service: service,
		Options: options,
	}
}

func CreateRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/", createURLShort)
	router.Get("/{short_url}", getlURLOriginal)
	return router
}

// createURLShort обрабатывает HTTP-запрос.
func createURLShort(res http.ResponseWriter, req *http.Request) {
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

	repo := repository.NewRepo()

	urlOriginal := string(bodyBytes)
	urlShort := repo.CreateURL(urlOriginal)

	bodyResult := fmt.Sprintf("http://localhost:8080/%s", urlShort)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(bodyResult))
}

// getlURLOriginal обрабатывает HTTP-запрос.
func getlURLOriginal(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	URLShort := chi.URLParam(req, "short_url")
	defer req.Body.Close()

	repo := repository.NewRepo()

	urlOriginal, err := repo.GetURL(URLShort)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.Header().Add("Location", urlOriginal)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}
