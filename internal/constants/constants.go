package constants

import (
	"errors"
	"time"
)

type contextKey string

const (
	HashLength            = 8
	TokenExp              = time.Hour * 3
	UserIDKey  contextKey = "userID"
)

var (
	// URL уже существует
	ErrorURLAlreadyExist = errors.New("URL already exists")
	// URL не существует
	ErrorURLNotExist = errors.New("URL doesn't exist")
	// Метод не разрешен
	ErrorMethodNotAllowed = errors.New("method not allowed")
	// URL уже удален
	ErrorURLAlreadyDeleted = errors.New("URL already deleted")
	// нет валидных данных
	ErrorNoData = errors.New("URL already deleted")
)
