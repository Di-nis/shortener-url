package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/Di-nis/shortener-url/internal/usecase"

	"github.com/go-chi/chi/v5"
)

// Controller - структура HTTP-хендлера.
type Controller struct {
	URLUseCase *usecase.URLUseCase
	Config     *config.Config
}

// NewСontroller - создание структуры Controller.
func NewСontroller(urlUseCase *usecase.URLUseCase, config *config.Config) *Controller {
	return &Controller{
		URLUseCase: urlUseCase,
		Config:     config,
	}
}

// CreateRouter - маршрутизация запросов.
func (c *Controller) CreateRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/api/shorten", c.createURLShortFromJSON)
	router.Post("/", c.createURLShortFromText)
	router.Get("/{short_url}", c.getlURLOriginal)
	return router
}

// createURLShortFromJSON - обрабатка HTTP-запроса: тип запроcа - POST, вовзвращает короткий URL.
func (c *Controller) createURLShortFromJSON(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, _ := io.ReadAll(req.Body)
	if reflect.DeepEqual(bodyBytes, []byte{}) {
		http.Error(res, "Не удалось прочитать тело запроса", http.StatusBadRequest)
		res.Header().Set("Content-Type", "")
		return
	}

	defer req.Body.Close()

	var (
		request  models.Request
		response models.Response
	)
	if err := json.Unmarshal(bodyBytes, &request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	urlShort, err := c.URLUseCase.CreateURL(request.URLOriginal)
	if err != nil {
		res.WriteHeader(http.StatusConflict)
		return
	}
	response.Result = fmt.Sprintf("%s/%s", c.Config.BaseURL, urlShort)

	bodyResult, err := json.Marshal(response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	_, err = res.Write([]byte(bodyResult))
	if err != nil {
		log.Fatal(err)
	}
}

// createURLShortFromText - обрабатка HTTP-запроса: тип запроcа - POST, вовзвращает короткий URL.
func (c *Controller) createURLShortFromText(res http.ResponseWriter, req *http.Request) {
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
		return
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
