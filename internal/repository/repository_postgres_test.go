package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"database/sql"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
)

func TestRepoPostgres_Ping(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "тест 1",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Skipf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectPing().WillReturnError(nil)

			repo := RepoPostgres{db: db}

			got := repo.Ping(context.Background())
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("TestRepoPostgres_Ping() = %v, wantErr: %v", got, tt.wantErr)
			}

		})
	}
}

func TestRepoPostgres_InsertOrdinary(t *testing.T) {
	tests := []struct {
		name    string
		url     models.URLBase
		dbErr   error
		wantErr error
	}{
		{
			name:    "тест 1",
			url:     testURLFull1,
			dbErr:   &pgconn.PgError{Code: "23505"},
			wantErr: constants.ErrorURLAlreadyExist,
		},
		{
			name:    "тест 2",
			url:     testURLFull1,
			dbErr:   errDB,
			wantErr: errDB,
		},
		{
			name:    "тест 3",
			url:     testURLFull3,
			dbErr:   nil,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Skipf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectExec(`INSERT INTO urls \(original, short, user_id\) VALUES \(\$1, \$2, \$3\)`).
				WithArgs(tt.url.Original, tt.url.Short, tt.url.UUID).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.dbErr)

			repo := RepoPostgres{db: db}

			gotErr := repo.InsertOrdinary(context.Background(), tt.url)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("TestRepoPostgres_InsertOrdinary() = %v, wantErr: %v", gotErr, tt.wantErr)
			}

		})
	}
}

func TestRepoPostgres_InsertBatch(t *testing.T) {
	tests := []struct {
		name         string
		urls         []models.URLBase
		dbErrPrepare error
		dbErr        error
		wantErr      error
	}{
		{
			name:         "тест 1",
			urls:         []models.URLBase{testURLFull1},
			dbErrPrepare: nil,
			dbErr:        &pgconn.PgError{Code: "23505"},
			wantErr:      constants.ErrorURLAlreadyExist,
		},
		{
			name:         "тест 2",
			urls:         []models.URLBase{testURLFull1},
			dbErrPrepare: nil,
			dbErr:        errDB,
			wantErr:      errDB,
		},
		{
			name:         "тест 3",
			urls:         []models.URLBase{testURLFull3},
			dbErrPrepare: nil,
			dbErr:        nil,
			wantErr:      nil,
		},
		{
			name:         "тест 4",
			urls:         []models.URLBase{testURLFull3},
			dbErrPrepare: errDBPrepare,
			dbErr:        nil,
			wantErr:      errDBPrepare,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Skipf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectBegin()

			prep := mock.ExpectPrepare(`INSERT INTO urls \(original, short, user_id\) VALUES \(\$1, \$2, \$3\)`)

			if tt.dbErrPrepare != nil {
				prep.WillReturnError(tt.dbErrPrepare)
			} else {
				for _, url := range tt.urls {
					prep.ExpectExec().
						WithArgs(url.Original, url.Short, url.UUID).
						WillReturnResult(sqlmock.NewResult(0, 1)).
						WillReturnError(tt.dbErr)
				}
			}

			mock.ExpectCommit()
			repo := RepoPostgres{db: db}

			got := repo.InsertBatch(context.Background(), tt.urls)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("TestRepoPostgres_InsertBatch() = %v, want: %v", got, tt.wantErr)
			}

		})
	}
}

func TestRepoPostgres_SelectShort(t *testing.T) {
	tests := []struct {
		name        string
		originalURL string
		dbRow       string
		dbErr       error
		want        string
		wantErr     error
	}{
		{
			name:        "тест 1",
			originalURL: url1,
			dbRow:       urlAlias1,
			dbErr:       nil,
			want:        urlAlias1,
			wantErr:     nil,
		},
		{
			name:        "тест 2",
			originalURL: url3,
			dbRow:       "",
			dbErr:       sql.ErrNoRows,
			want:        "",
			wantErr:     constants.ErrorURLNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Skipf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectQuery(`SELECT short FROM urls WHERE original = \$1`).
				WithArgs(tt.originalURL).
				WillReturnRows(sqlmock.NewRows([]string{"short"}).AddRow(tt.dbRow)).
				WillReturnError(tt.dbErr)

			repo := RepoPostgres{db: db}

			got, gotErr := repo.SelectShort(context.Background(), tt.originalURL)
			if got != tt.want {
				t.Errorf("TestRepoPostgres_SelectShort() = %v, want: %v", got, tt.want)
			}
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("TestRepoPostgres_SelectShort() = %v, wantErr: %v", gotErr, tt.wantErr)
			}

		})
	}
}

func TestRepoPostgres_SelectOriginal(t *testing.T) {
	tests := []struct {
		name     string
		shortURL string
		dbRow1   string
		dbRow2   bool
		dbErr    error
		want     string
		wantErr  error
	}{
		{
			name:     "тест 1",
			shortURL: urlAlias4,
			dbRow1:   url4,
			dbRow2:   true,
			dbErr:    nil,
			want:     "",
			wantErr:  constants.ErrorURLAlreadyDeleted,
		},
		{
			name:     "тест 2",
			shortURL: urlAlias1,
			dbRow1:   url1,
			dbRow2:   false,
			dbErr:    nil,
			want:     url1,
			wantErr:  nil,
		},
		{
			name:     "тест 3",
			shortURL: urlAlias3,
			dbErr:    sql.ErrNoRows,
			want:     "",
			wantErr:  constants.ErrorURLNotExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Skipf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectQuery(`SELECT original, is_deleted FROM urls WHERE short = \$1`).
				WithArgs(tt.shortURL).
				WillReturnRows(sqlmock.NewRows([]string{"original", "is_deleted"}).AddRow(tt.dbRow1, tt.dbRow2)).
				WillReturnError(tt.dbErr)

			repo := RepoPostgres{db: db}

			got, gotErr := repo.SelectOriginal(context.Background(), tt.shortURL)
			if got != tt.want {
				t.Errorf("TestRepoPostgres_SelectOriginal() = %v, want: %v", got, tt.want)
			}
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("TestRepoPostgres_SelectOriginal() = %v, wantErr: %v", gotErr, tt.wantErr)
			}

		})
	}
}

func TestRepoPostgres_SelectAll(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		dbRows       []models.URLBase
		dbErrPrepare error
		dbErr        error
		want         []models.URLBase
		wantErr      error
	}{
		{
			name:         "тест 1",
			userID:       UUID,
			dbRows:       testURLsShort,
			dbErrPrepare: nil,
			dbErr:        nil,
			want:         testURLsShort,
			wantErr:      nil,
		},
		{
			name:         "тест 2",
			userID:       UUID,
			dbRows:       []models.URLBase{},
			dbErrPrepare: nil,
			dbErr:        sql.ErrNoRows,
			want:         nil,
			wantErr:      sql.ErrNoRows,
		},
		{
			name:         "тест 3",
			userID:       UUID,
			dbRows:       []models.URLBase{},
			dbErrPrepare: errDBPrepare,
			dbErr:        nil,
			want:         nil,
			wantErr:      errDBPrepare,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Skipf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			row := sqlmock.NewRows([]string{"original", "short"})
			for _, r := range tt.dbRows {
				row.AddRow(r.Original, r.Short)
			}

			prep := mock.ExpectPrepare(`SELECT original, short FROM urls WHERE user_id = \$1`)

			if tt.dbErrPrepare != nil {
				prep.WillReturnError(tt.dbErrPrepare)
			} else {
				prep.ExpectQuery().
					WithArgs(tt.userID).
					WillReturnRows(row).
					WillReturnError(tt.dbErr)
			}

			repo := RepoPostgres{db: db}

			got, gotErr := repo.SelectAll(context.Background(), tt.userID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestRepoPostgres_SelectAll() = %v, want: %v", got, tt.want)
			}
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("TestRepoPostgres_SelectAll() = %v, wantErr: %v", gotErr, tt.wantErr)
			}

		})
	}
}

func TestRepoPostgres_Delete(t *testing.T) {
	tests := []struct {
		name          string
		urls          []models.URLBase
		dbErr         error
		dbRowAffected int64
		wantErr       error
	}{
		{
			name:          "тест 1",
			urls:          []models.URLBase{testURLFull1},
			dbErr:         nil,
			dbRowAffected: 1,
			wantErr:       nil,
		},
		{
			name:          "тест 2",
			urls:          []models.URLBase{},
			dbErr:         nil,
			dbRowAffected: 0,
			wantErr:       constants.ErrorNoData,
		},
		{
			name:          "тест 3",
			urls:          []models.URLBase{testURLFull1},
			dbErr:         errDB,
			dbRowAffected: 0,
			wantErr:       errDB,
		},
		{
			name:          "тест 4",
			urls:          []models.URLBase{testURLFull1},
			dbErr:         nil,
			dbRowAffected: 0,
			wantErr:       constants.ErrorNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Skipf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			var values []string

			for i := range tt.urls {
				base := i * 2
				params := fmt.Sprintf("($%d, $%d)", base+1, base+2)
				values = append(values, params)
			}

			query := `
			UPDATE urls AS u SET is_deleted = true FROM (VALUES ` + strings.Join(values, ",") + `) AS v(short, user_id) WHERE u.short = v.short AND u.user_id = v.user_id;`

			expectedExec := mock.ExpectExec(regexp.QuoteMeta(query))

			if tt.dbErr != nil {
				expectedExec.WillReturnError(tt.wantErr)
			} else {
				expectedExec.WillReturnResult(sqlmock.NewResult(1, tt.dbRowAffected))
			}

			repo := RepoPostgres{db: db}

			gotErr := repo.Delete(context.Background(), tt.urls)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("TestRepoPostgres_Delete() = %v, wantErr: %v", gotErr, tt.wantErr)
			}

		})
	}
}

func TestRepoPostgres_GetCountURLs(t *testing.T) {
	tests := []struct {
		name    string
		dbRow   int
		dbErr   error
		want    int
		wantErr error
	}{
		{
			name:    "тест 1",
			dbRow:   100,
			dbErr:   nil,
			want:    100,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Skipf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM urls`)).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.dbRow)).
				WillReturnError(tt.dbErr)

			repo := RepoPostgres{db: db}

			got, gotErr := repo.GetCountURLs(context.Background())
			if got != tt.want {
				t.Errorf("TestRepoPostgres_GetCountURLs() = %v, want: %v", got, tt.want)
			}
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("TestRepoPostgres_GetCountURLs() = %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestRepoPostgres_GetCountUsers(t *testing.T) {
	tests := []struct {
		name    string
		dbRow   int
		dbErr   error
		want    int
		wantErr error
	}{
		{
			name:    "тест 1",
			dbRow:   10,
			dbErr:   nil,
			want:    10,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Skipf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(DISTINCT user_id) FROM urls`)).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.dbRow)).
				WillReturnError(tt.dbErr)

			repo := RepoPostgres{db: db}

			got, gotErr := repo.GetCountUsers(context.Background())
			if got != tt.want {
				t.Errorf("TestRepoPostgres_GetCountUsers() = %v, want: %v", got, tt.want)
			}
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("TestRepoPostgres_GetCountUsers() = %v, wantErr: %v", gotErr, tt.wantErr)
			}
		})
	}
}
