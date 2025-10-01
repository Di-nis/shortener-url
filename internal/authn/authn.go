package authn

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"

	"net/http"

	"context"
	"time"
	"os"

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

func GenerateUserID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy).String()
	return id
}

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString(userID string, secretKey string) (string, error) {
	if userID == "" {
		return "", errors.New("user id is empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString, secretKey string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return "-1"
	}

	if !token.Valid {
		return "-1"
	}

	return claims.UserID
}
