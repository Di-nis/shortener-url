package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/models"
)

func WithAudit(auditFile, auditURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			var (
				action, userID, url string
				urlInOut            models.URLCopyOne
			)

			// определение поля user_id
			userID = req.Context().Value(constants.UserIDKey).(string)

			// определение поля action
			switch req.Method {
			case http.MethodGet:
				action = "follow"
			case http.MethodPost:
				action = "shorten"
			}

			if req.Method == http.MethodGet {
				url = res.Header().Get("Location")
				// TODO как перехватить url
				// url = ?
				fmt.Printf("%+v\n", res.Header())
			}
			if req.Method == http.MethodPost {
				bodyBytes, err := io.ReadAll(req.Body)
				if err != nil {
					logger.Sugar.Errorf("path: internal/audit/audit.go, errror - %w", err)
					// return
				}
				err = json.Unmarshal(bodyBytes, &urlInOut)
				if err != nil {
					url = string(bodyBytes)
				} else {
					url = urlInOut.Original
				}
			}

			if auditFile != "" {
				audit := models.NewAudit(action, userID, url)
				producer, err := NewProducer(auditFile)
				if err != nil {
					logger.Sugar.Errorf("path: internal/audit/audit.go, errror - %w", err)
					// return
				}
				defer producer.Close()

				err = producer.Write(audit)
				if err != nil {
					logger.Sugar.Errorf("path: internal/audit/audit.go, errror - %w", err)
					// return
				}
			}

			if auditURL != "" {
				audit := models.NewAudit(action, userID, url)
				data, err := json.Marshal(&audit)
				if err != nil {
					logger.Sugar.Errorf("path: internal/audit/audit.go, errror - %w", err)
					// return
				}

				client := &http.Client{}
				_, err = client.Post(auditURL, "application/json", bytes.NewBuffer(data))
				if err != nil {
					logger.Sugar.Errorf("path: internal/audit/audit.go, errror - %w", err)
					// return
				}
			}
			next.ServeHTTP(res, req)
		})
	}
}
