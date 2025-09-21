package authn

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

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
		fmt.Println("Token is not valid")
		return "-1"
	}

	fmt.Println("Token is valid")
	return claims.UserID
}
