package handler

import (
	"io"
	"net/http"
	"strings"
)

var (
	ShortenerArray = make(map[string]string)
)

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, webhook1)
	mux.HandleFunc(`/{id}`, webhook2)

	return http.ListenAndServe(":8080", mux)
}

// функция webhook1 — обработчик HTTP-запроса
func webhook1(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// разрешаем только POST-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Не удалось прочитать тело запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ShortenerArray["EwHXdJfB"] = string(bodyBytes)

	// установим правильный заголовок для типа данных
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	// пока установим ответ-заглушку, без проверки ошибок
	_, _ = w.Write([]byte("http://localhost:8080/EwHXdJfB"))
}

// функция webhook2 — обработчик HTTP-запроса
func webhook2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// разрешаем только Get-запросы
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	url := strings.Trim(r.URL.Path, "/")
	urlOut, ok := ShortenerArray[url]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusTemporaryRedirect)
	// пока установим ответ-заглушку, без проверки ошибок
	_, _ = w.Write([]byte(urlOut))
}
