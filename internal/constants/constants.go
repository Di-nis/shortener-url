package constants

import (
	"errors"
	"time"
)

const (
	HashLength = 8
	TOKEN_EXP = time.Hour * 3
)

var (
	// URL уже существует
	ErrorURLAlreadyExist = errors.New("URL already exists")
	// URL не существует
	ErrorURLNotExist = errors.New("URL doesn't exist")
	// Метод не разрешен
	ErrorMethodNotAllowed = errors.New("method not allowed")
)
