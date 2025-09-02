package constants

import (
	"errors"
)

var (
	// Ошибки:
	// URL уже существует
	ErrorURLAlreadyExist = errors.New("короткий URL уже существует")
	// URL не существует
	ErrorURLNotExist = errors.New("URL не существует")
)
