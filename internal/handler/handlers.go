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

func (c *Controller)CreateRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/", c.createURLShort)
	router.Get("/{short_url}", c.getlURLOriginal)
	return router
}

// createURLShort обрабатывает HTTP-запрос.
func (c *Controller) createURLShort(res http.ResponseWriter, req *http.Request) {
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

	bodyResult := fmt.Sprintf("%s/%s", c.Options.BaseURL, urlShort)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(bodyResult))
}

// getlURLOriginal обрабатывает HTTP-запрос.
func (c *Controller) getlURLOriginal(res http.ResponseWriter, req *http.Request) {
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
