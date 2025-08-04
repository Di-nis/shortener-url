package handler

import (
	// "bytes"
	"io"
	"net/http"
	"reflect"
	"strings"
	// "regexp"
)

var (
	ShortenerArray = make(map[string]string)
)

func Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, createShortURL)
	mux.HandleFunc(`/{id}`, getOriginalURL)

	return http.ListenAndServe(":8080", mux)
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

	// TODO Подумать, как сделать
	// bodyString := string(bodyBytes)
	// pattern := `/^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$/`
	// match, _ := regexp.MatchString(pattern, bodyString)

	// if !match {
	// 	http.Error(res, "Переданные данные не соответствуют структуре url-адреса", http.StatusBadRequest)
	// 	return
	// }

	defer req.Body.Close()

	// проверка на уже существование в БД
	// дополниеть реализацией связи с БД
	ShortenerArray["EwHXdJfB"] = string(bodyBytes)

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

	shortenerArray := make(map[string]string)
	shortenerArray["EwHXdJfB"] = "https://practicum.yandex.ru/ "

	url := strings.Trim(req.URL.Path, "/")
	defer req.Body.Close()

	headerLocation, ok := shortenerArray[url]
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusTemporaryRedirect)
	res.Header().Add("Location", headerLocation)
}
