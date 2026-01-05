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
	"github.com/Di-nis/shortener-url/internal/cidr"
	"github.com/Di-nis/shortener-url/internal/compress"
	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/models"

	"github.com/go-chi/chi/v5"

	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Pinger - интерфейс для проверки соединения с БД.
type Pinger interface {
	Ping(context.Context) error
}

// URLCreator - интерфейс, включащий методы по созданию URL.
type URLCreator interface {
	CreateURLOrdinary(context.Context, any) (models.URLBase, error)
	CreateURLBatch(context.Context, []models.URLBase) ([]models.URLBase, error)
}

// URLReader - интерфейс, включащий методы по получению URL.
type URLReader interface {
	GetOriginalURL(context.Context, string) (string, error)
	GetAllURLs(context.Context, string) ([]models.URLBase, error)
}

// URLDeleter - интерфейс, включащий методы по удалению URL.
type URLDeleter interface {
	DeleteURLs(context.Context, []models.URLBase) error
}

// URLStats - интерфейс, включающий методы по получению статистики.
type URLStats interface {
	GetStats(context.Context) (int, int, error)
}

// URLUseCase - объединенный интерфейс.
type URLUseCase interface {
	Pinger
	URLCreator
	URLReader
	URLDeleter
	URLStats
}

// Controller - структура HTTP-хендлера.
type Controller struct {
	Pinger     Pinger
	URLCreator URLCreator
	URLReader  URLReader
	URLDeleter URLDeleter
	URLStats   URLStats

	Config *config.Config
	Client *audit.Client
}

// NewСontroller - создание структуры Controller.
func NewСontroller(urlUseCase URLUseCase, config *config.Config) *Controller {
	return &Controller{
		Pinger:     urlUseCase,
		URLCreator: urlUseCase,
		URLReader:  urlUseCase,
		URLDeleter: urlUseCase,
		URLStats:   urlUseCase,
		Config:     config,
		Client:     audit.NewClient(&http.Client{}, config.AuditURL),
	}
}

// SetupRouter - маршрутизация запросов.
func (c *Controller) SetupRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(logger.WithLogging, compress.GzipMiddleware)
	c.UseAuthMiddleware(router)

	c.RegisterRoutes(router)
	return router
}

// UseAuthMiddleware - использование middleware для аутентификации.
func (c *Controller) UseAuthMiddleware(router *chi.Mux) {
	if c.Config.UseMockAuth {
		router.Use(authn.MockAuthMiddleware)
	} else {
		router.Use(authn.Middleware)
	}
}

// RegisterRoutes - регистрация маршрутов.
func (c *Controller) RegisterRoutes(router *chi.Mux) {
	router.Post("/api/shorten/batch", c.CreateURLShortJSONBatch)
	router.Get("/api/user/urls", c.getAllURLs)
	router.Delete("/api/user/urls", c.deleteURLs)
	router.Get("/ping", c.pingDB)

	router.Group(func(r chi.Router) {
		r.Use(cidr.WithCheckCIDR(c.Config.TrustedSubnet))

		r.Get("/api/internal/stats", c.stats)
	})

	router.Group(func(r chi.Router) {
		r.Use(audit.WithAudit(c.Client, c.Config.AuditFile))

		r.Post("/", c.createURLShortText)
		r.Post("/api/shorten", c.createURLShortJSON)
		r.Get("/{short_url}", c.getURLOriginal)
	})

	// pprof
	router.Mount("/debug/pprof/", http.HandlerFunc(pprof.Index))
	router.Get("/debug/pprof/cmdline", pprof.Cmdline)
	router.Get("/debug/pprof/profile", pprof.Profile)
	router.Get("/debug/pprof/symbol", pprof.Symbol)
	router.Get("/debug/pprof/trace", pprof.Trace)
}

// CreateURLShortJSONBatch - обрабатка HTTP-запроса: тип запроcа - POST, вовзвращает короткий URL.
func (c *Controller) CreateURLShortJSONBatch(res http.ResponseWriter, req *http.Request) {
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

	var urls []models.URLBase
	if err :=

		json.Unmarshal(bodyBytes, &urls); err != nil {
		http.Error(res, constants.InvalidJSONError, http.StatusBadRequest)
		return
	}

	for i := range urls {
		urls[i].UUID = userID
	}

	createdURLs, err := c.URLCreator.CreateURLBatch(ctx, urls)

	// Добавление базового URL
	for i := range createdURLs {
		createdURLs[i].Short = addBaseURLToResponse(c.Config.BaseURL, createdURLs[i].Short)
	}

	bodyResult, marshalErr := json.Marshal(createdURLs)
	if marshalErr != nil {
		http.Error(res, constants.InvalidJSONError, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	writeStatusCreate(res, err)

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
		urlInOut models.URLJSON
		url      models.URLBase
	)

	if err := json.Unmarshal(bodyBytes, &urlInOut); err != nil {
		http.Error(res, constants.InvalidJSONError, http.StatusBadRequest)
		return
	}

	urlInOut.UUID = userID

	url, err := c.URLCreator.CreateURLOrdinary(ctx, urlInOut)

	url.Short = addBaseURLToResponse(c.Config.BaseURL, url.Short)
	urlInOut = models.URLJSON(url)

	bodyResult, marshalErr := json.Marshal(urlInOut)
	if marshalErr != nil {
		http.Error(res, constants.InvalidJSONError, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	writeStatusCreate(res, err)

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

	urlIn := models.URLBase{
		Original: string(bodyBytes),
		UUID:     userID,
	}

	urlOut, err := c.URLCreator.CreateURLOrdinary(ctx, urlIn)

	urlOut.Short = addBaseURLToResponse(c.Config.BaseURL, urlOut.Short)

	res.Header().Set("Content-Type", "text/plain")
	writeStatusCreate(res, err)

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

	urls, err := c.URLReader.GetAllURLs(ctx, userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	var (
		urlsOut []models.URLGetAll
		urlOut  models.URLGetAll
	)

	for _, url := range urls {
		urlOut = models.URLGetAll(url)
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
	res.WriteHeader(http.StatusOK)

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

	urlOriginal, err := c.URLReader.GetOriginalURL(ctx, URLShort)
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
	urls := []models.URLBase{}

	if err := json.NewDecoder(req.Body).Decode(&shorts); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	for _, short := range shorts {
		urls = append(urls, models.URLBase{
			Short: short,
			UUID:  userID,
		})
	}

	if err := c.URLDeleter.DeleteURLs(ctx, urls); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

// pingDB - пинг БД.
func (c *Controller) pingDB(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	err := c.Pinger.Ping(ctx)
	writeStatusCodePing(res, err)
}

// stats - получение статистики по сокращенным URL.
func (c *Controller) stats(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	countURLs, countUsers, err := c.URLStats.GetStats(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	stats := models.NewStats(countURLs, countUsers)

	bodyResult, err := json.Marshal(stats)
	if err != nil {
		http.Error(res, constants.InvalidJSONError, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	_, err = res.Write([]byte(bodyResult))
	if err != nil {
		http.Error(res, constants.WriteResponseError, http.StatusInternalServerError)
	}
}
