package authn

import (
	"errors"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/golang-jwt/jwt/v5"

	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	SID    string
	UserID string
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString(secretKey, userID, sessionID string) (string, error) {
	if userID == "" {
		return "", errors.New("user id is empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.TokenExp)),
		},
		SID:    sessionID,
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
