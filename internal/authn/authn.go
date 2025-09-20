package authn

import (
	"errors"
	"fmt"
	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString(userID int, secretKey string) (string, error) {
	if userID == 0 {
		return "", errors.New("user id is zero")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.TOKEN_EXP)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString, secretKey string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return -1
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return -1
	}

	fmt.Println("Token os valid")
	return claims.UserID
}


// // getToken - получение токена из заголовка запроса.
// func getToken(headerAuthorization string) string {
// 	if headerAuthorization != "" {
// 		bearerToken := strings.Split(headerAuthorization, " ")
// 		if len(bearerToken) == 2 {
// 			return bearerToken[1]
// 		}
// 	}
// 	return ""
// }