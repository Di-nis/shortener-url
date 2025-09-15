package constants

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	HashLength = 8
)

var (
	// Ошибки:
	// URL уже существует
	ErrorURLAlreadyExist = errors.New("короткий URL уже существует")
	// URL не существует
	ErrorURLNotExist = errors.New("URL не существует")
	// Ошибка PostgreSQL
	PgErr *pgconn.PgError
)
