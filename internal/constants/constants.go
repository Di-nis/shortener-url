// Package constants реализывает константы и переменные для кастомных ошибок.
package constants

import (
	"errors"
	"time"
)

type contextKey string

// Для использования при создании хэша URL.
const (
	HashLength            = 8
	TokenExp              = time.Hour * 3
	UserIDKey  contextKey = "userID"
)

// Роли пользователей.
const (
	User = "user"
)

// Кастомные ошибки.
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
	// даныне не найдены
	ErrorNotFound = errors.New("URL not found")
)

// Тексты ошибок.
var (
	ReadRequestError   = "unable to read request body"
	EmptyBodyError     = "request body is empty"
	InvalidJSONError   = "invalid JSON format"
	InternalError      = "internal error"
	WriteResponseError = "error writing response"
	NoContentError     = "no content"
)
