package audit

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Di-nis/shortener-url/internal/constants"

	"github.com/Di-nis/shortener-url/internal/models"
)

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
func getURL(w http.ResponseWriter, r *http.Request) string {
	var (
		urlInOut models.URLCopyOne
		err      error
		url      string
	)

	if r.Method == http.MethodGet {
		url = w.Header().Get("Location")
	}
	// 	// TODO как перехватить url
	// 	// url = ?
	// 	fmt.Printf("%+v\n", res.Header())
	// }

	if r.Method == http.MethodPost {
		bodyBytes, _ := io.ReadAll(r.Body)

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		err = json.Unmarshal(bodyBytes, &urlInOut)
		if err != nil {
			url = string(bodyBytes)
		} else {
			url = urlInOut.Original
		}
	}
	return url
}

// saveLogsToFile - сохранение логов в файл.
func saveLogsToFile(auditFile, action, userID, url string) {
	if auditFile != "" {
		audit := NewAudit(action, userID, url)
		producer, _ := NewProducer(auditFile)
		defer producer.Close()

		producer.Write(audit)
	}
}

// sendLogsToURL - отправка логов на URL.
func sendLogsToURL(auditURL, action, userID, url string) {
	if auditURL != "" {
		audit := NewAudit(action, userID, url)
		data, _ := json.Marshal(&audit)

		client := &http.Client{}
		client.Post(auditURL, "application/json", bytes.NewBuffer(data))
	}
}

// WithAudit - middleware-аудит.
func WithAudit(auditFile, auditURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var action, userID, url string

			userID = r.Context().Value(constants.UserIDKey).(string)
			action = getAction(r.Method)

			next.ServeHTTP(w, r)

			url = getURL(w, r)

			saveLogsToFile(auditFile, action, userID, url)
			sendLogsToURL(auditURL, action, userID, url)

		})
	}
}
