package repository

import (
	"context"
	"database/sql"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

type RepoPostgres struct {
	dataSourceName string
}

func NewRepoPostgres(dataSourceName string) *RepoPostgres {
	return &RepoPostgres{
		dataSourceName: dataSourceName,
	}
}

func (repo *RepoPostgres) Migrations() error {
	db, err := sql.Open("postgres", repo.dataSourceName)
	if err != nil {return  err}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {return  err}

	m, err := migrate.NewWithDatabaseInstance(
		"file:migrations",
		"postgres",
		driver)

	if err != nil {return  err}
	m.Up()
	return nil
}

func (repo *RepoPostgres) Create(ctx context.Context, urls []models.URL) error {
	
	db, err := sql.Open("pgx", repo.dataSourceName)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := db.PrepareContext(ctx, "INSERT INTO urls (original, short) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	for _, url := range urls {
		_, err = stmt.ExecContext(ctx, url.Original, url.Short)
	}
		if err != nil {
			return err
		}
	// return result.LastInsertId()
	return tx.Commit()
}

func (repo *RepoPostgres) Get(ctx context.Context, urlShort string) (string, error) {
	db, err := sql.Open("pgx", repo.dataSourceName)
	if err != nil {
		return "", err
	}
	defer db.Close()
	row := db.QueryRowContext(ctx, "SELECT original FROM urls WHERE short = $1", urlShort)

	var URLOriginal string
	err = row.Scan(&URLOriginal)
	if err != nil {
		return "", constants.ErrorURLNotExist
	}
	return URLOriginal, nil
}
