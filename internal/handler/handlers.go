package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/Di-nis/shortener-url/internal/usecase"
	"github.com/Di-nis/shortener-url/internal/authn"

	"github.com/go-chi/chi/v5"

	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	router.Post("/api/shorten/batch", c.createURLShortJSONBatch)
	router.Post("/api/shorten", c.createURLShortJSON)
	router.Post("/", c.createURLShortText)
	router.Get("/api/user/urls", c.getAllURLs)
	router.Get("/{short_url}", c.getlURLOriginal)
	router.Get("/ping", c.pingDB)
	return router
}

// createURLShortJSON - обрабатка HTTP-запроса: тип запроcа - POST, вовзвращает короткий URL.
func (c *Controller) createURLShortJSONBatch(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

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

	var userID string
	// тут реализзовать вытаскивание userID из запроса

	defer req.Body.Close()

	var urls []models.URL
	if err := json.Unmarshal(bodyBytes, &urls); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	urls, err := c.URLUseCase.CreateURLBatch(ctx, urls, userID)
	if err != nil {
		res.WriteHeader(http.StatusConflict)
		return
	}

	for idx := range urls {
		urls[idx].Short = addBaseURLToResponse(c.Config.BaseURL, urls[idx].Short)
	}

	bodyResult, err := json.Marshal(urls)
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

// createURLShortJSON - обрабатка HTTP-запроса: тип запроcа - POST, вовзвращает короткий URL.
func (c *Controller) createURLShortJSON(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

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

	var userID string
	userID = "12rdscjh"
	// тут реализзовать вытаскивание userID из запроса

	defer req.Body.Close()

	var (
		urlInOut models.URLCopyOne
		url      models.URL
	)

	if err := json.Unmarshal(bodyBytes, &urlInOut); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	url, err := c.URLUseCase.CreateURLOrdinary(ctx, urlInOut, userID)

	url.Short = addBaseURLToResponse(c.Config.BaseURL, url.Short)
	urlInOut = models.URLCopyOne(url)

	bodyResult, err2 := json.Marshal(urlInOut)
	if err2 != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	getStatusCode(res, err)

	_, err = res.Write([]byte(bodyResult))
	if err != nil {
		log.Fatal(err)
	}
}

// createURLShortText - обрабатка HTTP-запроса: тип запроcа - POST, вовзвращает короткий URL.
func (c *Controller) createURLShortText(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, _ := io.ReadAll(req.Body)
	if reflect.DeepEqual(bodyBytes, []byte{}) {
		http.Error(res, "Не удалось прочитать тело запроса", http.StatusBadRequest)
		return
	}
	var userID string
	// тут реализзовать вытаскивание userID из запроса

	defer req.Body.Close()

	urlIn := models.URL{
		Original: string(bodyBytes),
	}

	urlOut, err := c.URLUseCase.CreateURLOrdinary(ctx, urlIn, userID)
	urlOut.Short = addBaseURLToResponse(c.Config.BaseURL, urlOut.Short)

	res.Header().Set("Content-Type", "text/plain")

	getStatusCode(res, err)

	_, err = res.Write([]byte(urlOut.Short))
	if err != nil {
		log.Fatal(err)
	}
}

// getAllURLs - получение всех когда-либо сокращенных пользователем URL.
func (c *Controller) getAllURLs(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	cookie, err := req.Cookie("auth_token")
	if err != nil {
		if err == http.ErrNoCookie {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	fmt.Println(cookie.Value)

	userID := authn.GetUserID(cookie.Value, c.Config.JWTSecret)
	fmt.Println(userID)

	// TODO требуется определить userID из запроса клиента
	defer req.Body.Close()

	urls, err := c.URLUseCase.GetAllURLs(ctx, userID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}

	var urlsOut []models.URLCopyThree
	var urlOut models.URLCopyThree

	for _, url := range urls {
		urlOut = models.URLCopyThree(url)
		urlsOut = append(urlsOut, urlOut)
	}

	if len(urlsOut) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	bodyResult, err2 := json.Marshal(urlsOut)
	if err2 != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	_, err = res.Write([]byte(bodyResult))
	if err != nil {
		log.Fatal(err)
	}
}

// getlURLOriginal - обрабатка HTTP-запроса: тип запроcа - GET, возвращает оригинальный URL.
func (c *Controller) getlURLOriginal(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	URLShort := chi.URLParam(req, "short_url")
	defer req.Body.Close()

	urlOriginal, err := c.URLUseCase.GetOriginalURL(ctx, URLShort)
	if err != nil && err == constants.ErrorURLNotExist {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Add("Location", urlOriginal)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (c *Controller) pingDB(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	if err := c.URLUseCase.Ping(ctx); err != nil {
		switch {
		case errors.Is(err, constants.ErrorMethodNotAllowed):
			res.WriteHeader(http.StatusMethodNotAllowed)
		default:
			res.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		res.WriteHeader(http.StatusOK)
	}
}
