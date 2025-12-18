// Package authn предоставляет middleware для авторизации и аутентификации пользователей.
package authn

import (
	"github.com/Di-nis/shortener-url/internal/constants"

	"net/http"

	"context"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// AuthMiddleware - аутентификация пользователя.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var userID, tokenString, sessionID string

		JWTSecret := os.Getenv("JWT_SECRET")
		cookie, err := req.Cookie("auth_token")
		if err != nil {
			tokenString = ""
		} else {
			tokenString = cookie.Value
		}

		if tokenString == "" {
			userID = GenerateUserID()
			sessionID = GenerateSessionID()
			newToken, err := BuildJWTString(JWTSecret, userID, sessionID)
			if err != nil {
				http.Error(res, "Ошибка создания токена", http.StatusInternalServerError)
				return
			}
			newCookie := &http.Cookie{
				Name:     "auth_token",
				Value:    newToken,
				Expires:  time.Now().Add(24 * time.Hour),
				Path:     "/",
				Domain:   "localhost",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				Secure:   true,
			}
			http.SetCookie(res, newCookie)
			res.Header().Set("Authorization", newToken)
		} else {
			claims, isTokenValid := GetClaims(tokenString, JWTSecret)
			if !isTokenValid {
				http.Error(res, "token not valid", http.StatusUnauthorized)
				return
			}
			userID = claims.UserID
			sessionID = claims.SID
			if sessionID == "" {
				http.Error(res, "session ID not valid", http.StatusUnauthorized)
				return
			}

			http.SetCookie(res, cookie)
			res.Header().Set("Authorization", tokenString)
		}

		ctx := context.WithValue(req.Context(), constants.UserIDKey, userID)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

// MockAuthMiddleware - аутентификация пользователя.
func MockAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), constants.UserIDKey, "01KA3YRQCWTNAJEGR5Z30PH6VT")
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
