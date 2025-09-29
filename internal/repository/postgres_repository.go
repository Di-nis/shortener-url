package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"sync"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

type RepoPostgres struct {
	db *sql.DB
}

func NewRepoPostgres(dataSourceName string) (*RepoPostgres, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &RepoPostgres{
		db: db,
	}, nil
}

func (repo *RepoPostgres) Ping(ctx context.Context) error {
	return repo.db.PingContext(ctx)
}

func (repo *RepoPostgres) Close() error {
	return repo.db.Close()
}

func (repo *RepoPostgres) Migrations() error {
	driver, err1 := postgres.WithInstance(repo.db, &postgres.Config{})
	if err1 != nil {
		return err1
	}

	m, err2 := migrate.NewWithDatabaseInstance(
		"file:migrations",
		"postgres",
		driver)
	if err2 != nil {
		return err2
	}

	err3 := m.Up()
	// тут надо подумать над обработкой ошибок
	if errors.Is(err3, migrate.ErrNoChange) {
		return nil
	} else if err3 != nil {
		return err3
	}
	return nil
}

func (repo *RepoPostgres) CreateOrdinary(ctx context.Context, url models.URL) error {
	stmt, err := repo.db.PrepareContext(ctx, "INSERT INTO urls (original, short, user_id) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, url.Original, url.Short, url.UUID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *RepoPostgres) CreateBatch(ctx context.Context, urls []models.URL) error {
	tx, err := repo.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := repo.db.PrepareContext(ctx, "INSERT INTO urls (original, short, user_id) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	for _, url := range urls {
		_, err = stmt.ExecContext(ctx, url.Original, url.Short, url.UUID)
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (repo *RepoPostgres) GetShortURL(ctx context.Context, urlOriginal string) (string, error) {
	stmt, err := repo.db.PrepareContext(ctx, "SELECT short FROM urls WHERE original = $1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, urlOriginal)

	var urlShort string
	err = row.Scan(&urlShort)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", constants.ErrorURLNotExist
	}
	return urlShort, nil
}

func (repo *RepoPostgres) GetOriginalURL(ctx context.Context, urlShort string) (string, error) {
	stmt, err := repo.db.PrepareContext(ctx, "SELECT original, is_deleted FROM urls WHERE short = $1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, urlShort)

	var url models.URL
	err = row.Scan(&url.Original, &url.DeletedFlag)

	if url.DeletedFlag {
		return "", constants.ErrorURLAlreadyDeleted
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", constants.ErrorURLNotExist
	}

	return url.Original, nil
}

// GetAllURLs - получение всех когда-либо сокращенных пользователем URL.
func (repo *RepoPostgres) GetAllURLs(ctx context.Context, userID string) ([]models.URL, error) {
	stmt, err := repo.db.PrepareContext(ctx, "SELECT original, short FROM urls WHERE user_id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	// требуется дополнительно уточнить, как писать
	if rows.Err() != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	// TODO на подумать над послежним параметром
	urls := make([]models.URL, 0, 20)

	for rows.Next() {
		var url models.URL
		err = rows.Scan(&url.Original, &url.Short)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (repo *RepoPostgres) DeleteURL(ctx context.Context, urls []models.URL) error {
	// fmt.Println("я тут", urls)
	var m sync.RWMutex
	tx, err := repo.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := repo.db.PrepareContext(ctx, "UPDATE urls SET is_deleted = true WHERE short = $1 AND user_id = $2")
	if err != nil {
		return err
	}

	for _, url := range urls {
		// fmt.Println("и тут", url)
		m.Lock()
		_, err = stmt.ExecContext(ctx, url.Short, url.UUID)
		// fmt.Println(err)
		m.Unlock()
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}
