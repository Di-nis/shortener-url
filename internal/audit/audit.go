// Package audit предоставляет middleware для аудита запросов.
// Реализации аудита: в файл, удаленный сервер.
package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Di-nis/shortener-url/internal/constants"

	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/models"
)

// attemptsCount - максимальное количество попыток при запросах на удаленный сервер.
const attemptsCount = 5

// Client - структура для отправки логов на URL.
type Client struct {
	httpClient *http.Client
	url        string
}

// NewClient - функция для создания нового экземпляра Client.
func NewClient(httpClient *http.Client, url string) *Client {
	return &Client{
		httpClient: httpClient,
		url:        url,
	}
}

// Audit - структура для хранения данных аудита.
type Audit struct {
	TS     int64  `json:"ts"`
	Action string `json:"action"`
	UserID string `json:"user_id"`
	URL    string `json:"url"`
}

// NewAudit - функция для создания нового экземпляра Audit.
func NewAudit(action, userID, url string) *Audit {
	return &Audit{
		TS:     time.Now().Unix(),
		Action: action,
		UserID: userID,
		URL:    url,
	}
}

// getAction - получение action.
func getAction(method string) string {
	switch method {
	case http.MethodGet:
		return "follow"
	case http.MethodPost:
		return "shorten"
	default:
		return ""
	}
}

// getURL - получение url.
func getURL(w http.ResponseWriter, r *http.Request) (string, error) {
	var (
		urlInOut models.URLJSON
		url      string
	)

	if r.Method == http.MethodGet {
		url = w.Header().Get("Location")
	}

	if r.Method == http.MethodPost {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return "", fmt.Errorf("path: internal/audit/audit.go, func getURL(), read body error: %w", err)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		err = json.Unmarshal(bodyBytes, &urlInOut)
		if err != nil {
			url = string(bodyBytes)
		} else {
			url = urlInOut.Original
		}
	}
	return url, nil
}

// saveLogsToFile - сохранение логов в файл.
func saveLogsToFile(auditFile, action, userID, url string) {
	audit := NewAudit(action, userID, url)
	producer, _ := NewProducer(auditFile)
	defer producer.Close()

	err := producer.Write(audit)
	if err != nil {
		logger.Sugar.Info("path: internal/audit/audit.go, func saveLogsToFile(), save to file error", err.Error())
	}
}

// sendLogsToURL - отправка логов на URL.
func sendLogsToURL(ctx context.Context, client *Client, action, userID, url string) error {
	var (
		lastResErr, reqErr error
		req                *http.Request
	)

	for range attemptsCount {
		audit := NewAudit(action, userID, url)
		data, err := json.Marshal(&audit)
		if err != nil {
			logger.Sugar.Info(
				"path: internal/audit/audit.go, func sendLogsToURL(), marshal error",
				err.Error(),
			)
		}

		req, reqErr = http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			client.url,
			bytes.NewBuffer(data),
		)
		if reqErr != nil {
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.httpClient.Do(req)
		if err != nil {
			lastResErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			lastResErr = nil
			break
		}
	}

	if lastResErr == nil && reqErr == nil {
		return nil
	}

	logger.Sugar.Warnw(
		"path: internal/audit/audit.go, func sendLogsToURL(), audit retry",
		"attempt", attemptsCount,
		"err", lastResErr,
	)

	select {
	case <-time.After(time.Duration(attemptsCount) * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}

	return lastResErr
}

// WithAudit - middleware-аудит.
func WithAudit(client *Client, auditFile string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				action, userID, url string
				err                 error
			)

			userID = r.Context().Value(constants.UserIDKey).(string)
			action = getAction(r.Method)

			url, err = getURL(w, r)
			if err != nil {
				logger.Sugar.Info(
					"path: internal/audit/audit.go, func WithAudit(), get url error",
					err.Error(),
				)
			}

			next.ServeHTTP(w, r)

			if auditFile != "" {
				saveLogsToFile(auditFile, action, userID, url)
			}

			if client.url == "" {
				return
			}

			if err := sendLogsToURL(r.Context(), client, action, userID, url); err != nil {
				logger.Sugar.Info(
					"path: internal/audit/audit.go, func WithAudit(), send to audit URL error",
					err.Error(),
				)
			}
		})
	}
}
