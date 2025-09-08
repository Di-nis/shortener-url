package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Di-nis/shortener-url/internal/constants"
	_ "github.com/jackc/pgx/v5/stdlib"
    "strings"
    // "time"
)

type Config struct {
	Localhost string
	Port      string
	User      string
	Password  string
	Name      string
    SSLMode string
}

func NewConfig(dataBaseDSN string) *Config {
    // host=localhost port=5432 user=postgres password=postgres dbname=urls sslmode=disable
	dataBaseDSNArray := strings.Split(dataBaseDSN, " ")
	localhost, post, user, password, name, sslMode := dataBaseDSNArray[0], dataBaseDSNArray[1], dataBaseDSNArray[2], dataBaseDSNArray[3], dataBaseDSNArray[4], dataBaseDSNArray[5]
	return &Config{
		Localhost: localhost,
		Port:      post,
		User:      user,
		Password:  password,
		Name:      name,
        SSLMode: sslMode,
	}
}

type RepoPostgres struct {
	config *Config
    configStr string
}

func NewRepoPostgres(config *Config, configStr string) *RepoPostgres {
	return &RepoPostgres{
		config: config,
        configStr: configStr,
	}
}

func (repo *RepoPostgres) Create(ctx context.Context, urlOriginal, urlShort string) error {
    db, err := sql.Open("pgx", repo.configStr)
	if err != nil {return err}
	defer db.Close()
    // тут должно быть result, err
	_, err = db.ExecContext(ctx, "INSERT INTO urls (original, short) VALUES ($1, $2)", urlOriginal, urlShort)
    // написать возврат ошибки, когда запись уже существует
	if err != nil {
	    return err
	}
    // return result.LastInsertId()
    return nil
}

func (repo *RepoPostgres) Get(ctx context.Context, urlShort string) (string, error) {
    db, err := sql.Open("pgx", repo.configStr)
	if err != nil {return "", err}
	defer db.Close()
	row := db.QueryRowContext(ctx, "SELECT original FROM urls WHERE short = ?", urlShort)

	var URLOriginal string
	err = row.Scan(&URLOriginal)
	if err != nil {
		// в константах создать ошибку
		return "", errors.New("new errors get")
	}
	if URLOriginal == "" {
		return "", constants.ErrorURLNotExist
	}
	return URLOriginal, nil
}
