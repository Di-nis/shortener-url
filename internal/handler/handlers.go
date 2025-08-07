package handler

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"reflect"
	"database/sql"
)


func InsertURL(dataBase *sql.DB, original, short string) error {
    insertUser := `INSERT INTO urls (original, short) VALUES (?, ?)`
    _, err := dataBase.Exec(insertUser, original, short)
	if err != nil {
		return err
	}
	return nil
}


func GetURL(dataBase *sql.DB, short string) (string, error) {
    URL := `SELECT original from urls WHERE short = (?)`
	var original string

    _ = dataBase.QueryRow(URL, short).Scan(&original)

	return original, nil
}

func CreateRouter() http.Handler {
	router := chi.NewRouter()

	router.Post("/", createShortURL)
	router.Get("/{short_url}", getOriginalURL)
	return router
}

// createShortURL обрабатывает HTTP-запрос.
func createShortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, _ := io.ReadAll(req.Body)
	if reflect.DeepEqual(bodyBytes, []byte{}) {
		http.Error(res, "Не удалось прочитать тело запроса", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	dataBase, _ := sql.Open("sqlite", "./shortener_database.db")
	InsertURL(dataBase, string(bodyBytes), "EwHXdJfB")

	defer dataBase.Close()

	// проверка на уже существование в БД
	// дополниеть реализацией связи с БД
	// ShortenerArray["EwHXdJfB"] = string(bodyBytes)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("http://localhost:8080/EwHXdJfB"))
}

// getOriginalURL обрабатывает HTTP-запрос.
func getOriginalURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	dataBase, _ := sql.Open("sqlite", "./shortener_database.db")
	defer dataBase.Close()


	shortURL := chi.URLParam(req, "short_url")
	defer req.Body.Close()

	headerLocation, _ := GetURL(dataBase, shortURL)

	// if !ok {
	// 	res.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	res.Header().Add("Location", headerLocation)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
}
