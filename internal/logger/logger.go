package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Log *zap.Logger.
var Log *zap.Logger = zap.NewNop()

// Sugar *zap.SugaredLogger.
var Sugar *zap.SugaredLogger

// responseData - структура для хранения данных о ответе.
type responseData struct {
	status int
	size   int
}

// loggingResponseWriter - реализация http.ResponseWriter, который логирует запросы.
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write - пишет данные в loggingResponseWriter.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader - пишет заголовок в loggingResponseWriter.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Initialize - инициализирует логгер.
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}

// WithLogging - middleware-логгер.
func WithLogging(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		uri := r.RequestURI
		method := r.Method

		next.ServeHTTP(&lw, r)
		duration := time.Since(start)

		Log.Sugar().Infoln(
			"uri", uri,
			"method", method,
			"status", responseData.status,
			"size", responseData.size,
			"duration", duration,
		)

	}
	return http.HandlerFunc(logFn)
}
