// Package handler предоставляет HTTP-транспорт для приложения, реализуя обработчики
// маршрутов и преобразуя сетевые запросы в вызовы сервисного слоя.
package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/pprof"
	"reflect"

	"github.com/Di-nis/shortener-url/internal/audit"
	"github.com/Di-nis/shortener-url/internal/authn"
	"github.com/Di-nis/shortener-url/internal/compress"
	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/models"

	"github.com/go-chi/chi/v5"

	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// URLUseCase - интерфейс для бизнес-логики.
type URLUseCase interface {
	Ping(context.Context) error
	CreateURLOrdinary(context.Context, any) (models.URL, error)
	CreateURLBatch(context.Context, []models.URL) ([]models.URL, error)
	GetOriginalURL(context.Context, string) (string, error)
	GetAllURLs(context.Context, string) ([]models.URL, error)
	DeleteURLs(context.Context, []models.URL) error
}

// Controller - структура HTTP-хендлера.
type Controller struct {
	URLUseCase URLUseCase
	Config     *config.Config
}

// NewСontroller - создание структуры Controller.
func NewСontroller(urlUseCase URLUseCase, config *config.Config) *Controller {
	return &Controller{
		URLUseCase: urlUseCase,
		Config:     config,
	}
}

// CreateRouter - маршрутизация запросов.
func (c *Controller) CreateRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(authn.AuthMiddleware, logger.WithLogging, compress.GzipMiddleware)

	router.Post("/api/shorten/batch", c.createURLShortJSONBatch)
	router.Get("/api/user/urls", c.getAllURLs)
	router.Delete("/api/user/urls", c.deleteURLs)
	router.Get("/ping", c.pingDB)

	router.Group(func(r chi.Router) {
		r.Use(audit.WithAudit(c.Config.AuditFile, c.Config.AuditURL))

		r.Post("/", c.createURLShortText)
		r.Post("/api/shorten", c.createURLShortJSON)
		r.Get("/{short_url}", c.getURLOriginal)
	})

	router.Mount("/debug/pprof/", http.HandlerFunc(pprof.Index)) // для перенаправления
	router.Get("/debug/pprof/cmdline", pprof.Cmdline)
	router.Get("/debug/pprof/profile", pprof.Profile)
	router.Get("/debug/pprof/symbol", pprof.Symbol)
	router.Get("/debug/pprof/trace", pprof.Trace)

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

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, constants.ReadRequestError, http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	if len(bodyBytes) == 0 {
		http.Error(res, constants.EmptyBodyError, http.StatusBadRequest)
		return
	}

	userID := req.Context().Value(constants.UserIDKey).(string)

	var urls []models.URL
	if err := json.Unmarshal(bodyBytes, &urls); err != nil {
		http.Error(res, constants.InvalidJSONError, http.StatusBadRequest)
		return
	}

	for i := range urls {
		urls[i].UUID = userID
	}

	createdURLs, err := c.URLUseCase.CreateURLBatch(ctx, urls)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			res.WriteHeader(http.StatusConflict)
		} else {
			http.Error(res, constants.InternalError, http.StatusInternalServerError)
		}
		return
	}

	// Добавление базового URL
	for i := range createdURLs {
		createdURLs[i].Short = addBaseURLToResponse(c.Config.BaseURL, createdURLs[i].Short)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	bodyResult, err := json.Marshal(createdURLs)
	if err != nil {
		http.Error(res, constants.InvalidJSONError, http.StatusInternalServerError)
		return
	}

	_, err = res.Write([]byte(bodyResult))
	if err != nil {
		http.Error(res, constants.WriteResponseError, http.StatusInternalServerError)
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
		http.Error(res, constants.EmptyBodyError, http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	// Получение userID через middleware Auth
	userID := req.Context().Value(constants.UserIDKey).(string)

	var (
		urlInOut models.URLCopyOne
		url      models.URL
	)

	if err := json.Unmarshal(bodyBytes, &urlInOut); err != nil {
		http.Error(res, constants.InvalidJSONError, http.StatusBadRequest)
		return
	}

	urlInOut.UUID = userID

	url, err := c.URLUseCase.CreateURLOrdinary(ctx, urlInOut)

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
		http.Error(res, constants.WriteResponseError, http.StatusInternalServerError)
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
		http.Error(res, constants.EmptyBodyError, http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	// Получение userID через middleware Auth
	userID := req.Context().Value(constants.UserIDKey).(string)

	urlIn := models.URL{
		Original: string(bodyBytes),
		UUID:     userID,
	}

	urlOut, err := c.URLUseCase.CreateURLOrdinary(ctx, urlIn)
	urlOut.Short = addBaseURLToResponse(c.Config.BaseURL, urlOut.Short)

	res.Header().Set("Content-Type", "text/plain")

	getStatusCode(res, err)

	_, err = res.Write([]byte(urlOut.Short))
	if err != nil {
		http.Error(res, constants.WriteResponseError, http.StatusInternalServerError)
	}
}

// getAllURLs - получение всех когда-либо сокращенных пользователем URL.
func (c *Controller) getAllURLs(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer req.Body.Close()

	var err error
	userID := req.Context().Value(constants.UserIDKey).(string)

	urls, err := c.URLUseCase.GetAllURLs(ctx, userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	var (
		urlsOut []models.URLCopyFour
		urlOut  models.URLCopyFour
	)

	for _, url := range urls {
		urlOut = models.URLCopyFour(url)
		urlOut.Short = addBaseURLToResponse(c.Config.BaseURL, urlOut.Short)
		urlsOut = append(urlsOut, urlOut)
	}

	if len(urlsOut) == 0 {
		http.Error(res, constants.NoContentError, http.StatusNoContent)
		return
	}

	bodyResult, err := json.Marshal(urlsOut)
	if err != nil {
		http.Error(res, constants.InvalidJSONError, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	_, err = res.Write([]byte(bodyResult))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// getURLOriginal - обрабатка HTTP-запроса: тип запроcа - GET, возвращает оригинальный URL.
func (c *Controller) getURLOriginal(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	URLShort := chi.URLParam(req, "short_url")
	defer req.Body.Close()

	urlOriginal, err := c.URLUseCase.GetOriginalURL(ctx, URLShort)
	if err != nil {
		if errors.Is(err, constants.ErrorURLNotExist) {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, constants.ErrorURLAlreadyDeleted) {
			res.WriteHeader(http.StatusGone)
			return
		}

	}
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Add("Location", urlOriginal)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

// pingDB - пинг БД.
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

// deleteURLs - удаление сокращенных URL.
func (c *Controller) deleteURLs(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if req.Method != http.MethodDelete {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Получение userID через middleware Auth
	userID := req.Context().Value(constants.UserIDKey).(string)

	var shorts []string
	urls := []models.URL{}

	if err := json.NewDecoder(req.Body).Decode(&shorts); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	for _, short := range shorts {
		urls = append(urls, models.URL{
			Short: short,
			UUID:  userID,
		})
	}

	if err := c.URLUseCase.DeleteURLs(ctx, urls); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}
