package repository

import (
	"context"
	"database/sql"

	"github.com/Di-nis/shortener-url/internal/constants"
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

func (repo *RepoPostgres) Create(ctx context.Context, urlOriginal, urlShort string) error {
	db, err := sql.Open("pgx", repo.dataSourceName)
	if err != nil {
		return err
	}
	defer db.Close()

	// тут должно быть result, err
	_, err = db.ExecContext(ctx, "INSERT INTO urls (original, short) VALUES ($1, $2)", urlOriginal, urlShort)
	// написать более подробно обработку ошибок
	if err != nil {
		return err
	}
	// return result.LastInsertId()
	return nil
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
