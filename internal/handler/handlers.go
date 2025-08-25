package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type URL struct {
	URLOriginal string
	URLShort    string
}

// Кастомная сериализация
func (u URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		URLOriginal string `json:"-"`
		URLShort    string `json:"result"`
	}{
		URLShort: u.URLShort,
	})
}

// Кастомная десериализация
func (u *URL) UnmarshalJSON(data []byte) error {
	aux := struct {
		URLOriginal string `json:"url"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	u.URLOriginal = aux.URLOriginal
	return nil
}

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

	router.Post("/api/shorten", c.createURLShortJSON)
	router.Post("/", c.createURLShort)
	router.Get("/{short_url}", c.getlURLOriginal)
	return router
}

// createURLShortJSON - обрабатка HTTP-запроса: тип запроcа - POST, вовзвращает короткий URL.
func (c *Controller) createURLShortJSON(res http.ResponseWriter, req *http.Request) {
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

	var url URL
	if err := json.Unmarshal(bodyBytes, &url); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	urlShort, err := c.URLUseCase.CreateURL(url.URLOriginal)
	if err != nil {
		res.WriteHeader(http.StatusConflict)
	}
	url.URLShort = urlShort

	bodyResult, err := json.Marshal(url)
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
