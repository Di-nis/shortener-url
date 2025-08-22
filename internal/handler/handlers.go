package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/usecase"

	"github.com/go-chi/chi/v5"
)

// Controller - структура HTTP-хендлера.
type Controller struct {
	URLUseCase *usecase.URLUseCase
	Config    *config.Config
}

// NewСontroller - создание структуры Controller.
func NewСontroller(urlUseCase *usecase.URLUseCase, config *config.Config) *Controller {
	return &Controller{
		URLUseCase: urlUseCase,
		Config:    config,
	}
}

// CreateRouter - маршрутизация запросов.
func (c *Controller) CreateRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/", c.createURLShort)
	router.Get("/{short_url}", c.getlURLOriginal)
	return router
}

// createURLShort - обрабатка HTTP-запроса: тип запроcа - POST, вовзвращает короткий URL.
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

	urlOriginal := string(bodyBytes)
	urlShort, err := c.URLUseCase.CreateURL(urlOriginal)
	if err != nil {
		res.WriteHeader(http.StatusConflict)
	}

	bodyResult := fmt.Sprintf("%s/%s", c.Config.BaseURL, urlShort)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)

	_, err = res.Write([]byte(bodyResult))
	if err != nil {
		log.Fatal(err)
	}
}

// getlURLOriginal - обрабатка HTTP-запроса: тип запроcа - GET, вовзвращает оригинальный URL.
func (c *Controller) getlURLOriginal(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	URLShort := chi.URLParam(req, "short_url")
	defer req.Body.Close()

	urlOriginal, err := c.URLUseCase.GetURL(URLShort)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	res.Header().Add("Location", urlOriginal)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}
