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
		var userID string
		JWTSecret := os.Getenv("JWT_SECRET")
		tokenString := ""
		cookie, err := req.Cookie("auth_token")
		if err != nil {
			tokenString = ""
		} else {
			tokenString = cookie.Value
		}

		if tokenString == "" {
			userID = GenerateUserID()
			newToken, err := BuildJWTString(userID, JWTSecret)
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
			}
			http.SetCookie(res, newCookie)
			res.Header().Set("Authorization", newToken)
		} else {
			userID = GetUserID(tokenString, JWTSecret)
			if userID == "-1" {
				http.Error(res, "Невалидный токен", http.StatusUnauthorized)
				return
			}
			http.SetCookie(res, cookie)
			res.Header().Set("Authorization", tokenString)
		}

		ctx := context.WithValue(req.Context(), constants.UserIDKey, userID)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
