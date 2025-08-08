package handler

import (
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/service"

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

func (controller *Controller) CreateRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/", controller.createShortURL)
	router.Get("/{short_url}", controller.getOriginalURL)
	return router
}

// createShortURL обрабатывает HTTP-запрос.
func (controller *Controller) createShortURL(res http.ResponseWriter, req *http.Request) {
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
	urlShort := controller.Service.Repo.CreateURL(urlOriginal)
	bodyResult := fmt.Sprintf("%s/%s", controller.Options.BaseURL, urlShort)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(bodyResult))
}

// getOriginalURL обрабатывает HTTP-запрос.
func (controller *Controller) getOriginalURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	URLShort := chi.URLParam(req, "short_url")
	defer req.Body.Close()

	urlOriginal, err := controller.Service.Repo.GetURL(URLShort)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.Header().Add("Location", urlOriginal)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}
