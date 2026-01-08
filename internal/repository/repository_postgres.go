// Package repository содержать инфраструктурные реализации доступа к данным.
// Реализации: PostgreSQL, хранение в памяти приложения и файловое хранилище.
// Сервисный слой работает только с интерфейсами, не зная деталей реализации.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"fmt"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

// RepoPostgres - репозиторий для работы с БД Postgres.
type RepoPostgres struct {
	db *sql.DB
}

// NewRepoPostgres - конструктор репозитория.
func NewRepoPostgres(dataSourceName string) (*RepoPostgres, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &RepoPostgres{
		db: db,
	}, nil
}

// Ping - проверка соединения с БД.
func (repo *RepoPostgres) Ping(ctx context.Context) error {
	return repo.db.PingContext(ctx)
}

// Close - закрытие соединения с БД.
func (repo *RepoPostgres) Close() error {
	return repo.db.Close()
}

// Migrations - миграции БД.
func (repo *RepoPostgres) Migrations() error {
	var err error
	driver, err := postgres.WithInstance(repo.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("path: internal/repository/postgres_repository.go, func Migrations(), failed create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:migrations",
		"postgres",
		driver)
	if err != nil {
		return fmt.Errorf("path: internal/repository/postgres_repository.go, func Migrations(), failed new Migrate Instance: %w", err)
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	} else if err != nil {
		return fmt.Errorf("path: internal/repository/postgres_repository.go, func Migrations(), failed make migrations: %w", err)
	}
	return nil
}

// InsertOrdinary - добавление ординарного URL в БД.
func (repo *RepoPostgres) InsertOrdinary(ctx context.Context, url models.URLBase) error {
	query := "INSERT INTO urls (original, short, user_id) VALUES ($1, $2, $3)"
	_, err := repo.db.ExecContext(ctx, query, url.Original, url.Short, url.UUID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("path: internal/repository/postgres_repository.go, func InsertOrdinary(), failed to insert url: %w", constants.ErrorURLAlreadyExist)
		}
		return fmt.Errorf("path: internal/repository/postgres_repository.go, func InsertOrdinary(), failed to insert url: %w", err)
	}
	return nil
}

// InsertBatch - добавление нескольких URL в БД.
func (repo *RepoPostgres) InsertBatch(ctx context.Context, urls []models.URLBase) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := repo.db.PrepareContext(ctx, "INSERT INTO urls (original, short, user_id) VALUES ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("path: internal/repository/postgres_repository.go, func InsertBatch(), failed to prepare statement: %w", err)
	}

	for _, url := range urls {
		_, err = stmt.ExecContext(ctx, url.Original, url.Short, url.UUID)
	}
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("path: internal/repository/postgres_repository.go, func InsertBatch(), failed to insert urls: %w", constants.ErrorURLAlreadyExist)
		}
		return fmt.Errorf("path: internal/repository/postgres_repository.go, func InsertBatch(), failed to insert urls: %w", err)
	}
	return tx.Commit()
}

// SelectShort - получение короткого URL по оригинальному.
func (repo *RepoPostgres) SelectShort(ctx context.Context, urlOriginal string) (string, error) {
	query := "SELECT short FROM urls WHERE original = $1"
	row := repo.db.QueryRowContext(ctx, query, urlOriginal)

	var urlShort string
	err := row.Scan(&urlShort)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("path: internal/repository/postgres_repository.go, func SelectShort(): %w", constants.ErrorURLNotExist)
	}
	return urlShort, nil
}

// SelectOriginal - получение оригинального URL по короткому.
func (repo *RepoPostgres) SelectOriginal(ctx context.Context, urlShort string) (string, error) {
	query := "SELECT original, is_deleted FROM urls WHERE short = $1"
	row := repo.db.QueryRowContext(ctx, query, urlShort)

	var url models.URLBase
	err := row.Scan(&url.Original, &url.DeletedFlag)

	if url.DeletedFlag {
		return "", fmt.Errorf("path: internal/repository/postgres_repository.go, func SelectOriginal(): %w", constants.ErrorURLAlreadyDeleted)
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("path: internal/repository/postgres_repository.go, func SelectOriginal(): %w", constants.ErrorURLNotExist)
	}

	return url.Original, nil
}

// SelectAll - получение всех когда-либо сокращенных пользователем URL.
func (repo *RepoPostgres) SelectAll(ctx context.Context, userID string) ([]models.URLBase, error) {
	stmt, err := repo.db.PrepareContext(ctx, "SELECT original, short FROM urls WHERE user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("path: internal/repository/postgres_repository.go, func SelectAll(), failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("path: internal/repository/postgres_repository.go, func SelectAll(), failed to get urls: %w", err)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("path: internal/repository/postgres_repository.go, func SelectAll(), row iteration failed: %w", err)
	}

	urls := make([]models.URLBase, 0, 20)

	for rows.Next() {
		var url models.URLBase
		err = rows.Scan(&url.Original, &url.Short)
		if err != nil {
			return nil, fmt.Errorf("path: internal/repository/postgres_repository.go, func SelectAll(), failed to scan url: %w", err)
		}

		urls = append(urls, url)
	}

	return urls, nil
}

// Delete - удаление URL из БД.
func (repo *RepoPostgres) Delete(ctx context.Context, urls []models.URLBase) error {
	if len(urls) == 0 {
		return constants.ErrorNoData
	}

	var values []string
	var args []any

	for i, url := range urls {
		base := i * 2
		params := fmt.Sprintf("($%d, $%d)", base+1, base+2)
		values = append(values, params)
		args = append(args, url.Short, url.UUID)
	}

	query := `
	UPDATE urls AS u SET is_deleted = true FROM (VALUES ` + strings.Join(values, ",") + `) AS v(short, user_id) WHERE u.short = v.short AND u.user_id = v.user_id;`

	result, err := repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("path: internal/repository/postgres_repository.go, func Delete(), failed to delete url: %w", err)
	}

	count, _ := result.RowsAffected()
	if count == 0 {
		return fmt.Errorf("path: internal/repository/postgres_repository.go, func Delete(), url not found: %w", constants.ErrorNotFound)
	}
	return nil
}

// GetCountURLs - получение количества записей.
func (repo *RepoPostgres) GetCountURLs(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM urls`
	row := repo.db.QueryRowContext(ctx, query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("path: internal/repository/postgres_repository.go, func GetCountURLs(): %w", err)
	}
	return count, nil
}

// GetCountUsers - получение количества уникальных пользователей.
func (repo *RepoPostgres) GetCountUsers(ctx context.Context) (int, error) {
	query := "SELECT COUNT(DISTINCT user_id) FROM urls"
	row := repo.db.QueryRowContext(ctx, query)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("path: internal/repository/postgres_repository.go, func GetCountUsers(): %w", err)
	}
	return count, nil
}
