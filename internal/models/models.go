// Package models - реализация моделей.
package models

import (
	"encoding/json"
)

// User - модель пользователя.
type User struct {
	ID int
}

// URLBase - основная модель для сущности url.
type URLBase struct {
	UUID        string `db:"user_id"`
	Short       string `db:"short"`
	Original    string `db:"original"`
	URLID       string
	DeletedFlag bool `db:"is_deleted"`
}

// MarshalJSON - метод для сериализации модели URL.
func (url URLBase) MarshalJSON() ([]byte, error) {
	urlAlias := struct {
		Short    string `json:"short_url"`
		Original string `json:"-"`
		URLID    string `json:"correlation_id"`
	}{
		Short:    url.Short,
		Original: url.Original,
		URLID:    url.URLID,
	}

	return json.Marshal(urlAlias)
}

// UnmarshalJSON - метод для десериализации модели URL.
func (url *URLBase) UnmarshalJSON(data []byte) error {
	type URLAlias struct {
		Short    string `json:"-"`
		Original string `json:"original_url"`
		URLID    string `json:"correlation_id"`
	}

	var urlAlias URLAlias

	if err := json.Unmarshal(data, &urlAlias); err != nil {
		return err
	}
	url.Original = urlAlias.Original
	url.URLID = urlAlias.URLID
	return nil
}

// URLJSON - сопутствующая модель для сущности url.
type URLJSON struct {
	UUID        string
	Short       string
	Original    string
	URLID       string
	DeletedFlag bool
}

// MarshalJSON - метод для сериализации модели URL.
func (url URLJSON) MarshalJSON() ([]byte, error) {
	urlAlias := struct {
		Short string `json:"result"`
	}{
		Short: url.Short,
	}

	return json.Marshal(urlAlias)
}

// UnmarshalJSON - метод для десериализации модели URL.
func (url *URLJSON) UnmarshalJSON(data []byte) error {
	type URLAlias struct {
		Original string `json:"url"`
	}

	var urlAlias URLAlias

	if err := json.Unmarshal(data, &urlAlias); err != nil {
		return err
	}
	url.Original = urlAlias.Original
	return nil
}

// URLStorage - сопутствующая модель для сущности url.
type URLStorage struct {
	UUID        string `json:"uuid"`
	Short       string `json:"url_short"`
	Original    string `json:"url_original"`
	URLID       string `json:"-"`
	DeletedFlag bool   `json:"-"`
}

// URLGetAll - модель URL.
type URLGetAll struct {
	UUID        string `json:"-"`
	Short       string `json:"short_url"`
	Original    string `json:"original_url"`
	URLID       string `json:"-"`
	DeletedFlag bool   `json:"-"`
}

// Pooler - интерфейс для пула.
type Pooler interface {
	Reset()
}

// Pool - пул.
type Pool[T Pooler] struct {
	Buf []T
}

// NewPool - конструктор пула.
func NewPool[T Pooler]() *Pool[Pooler] {
	return &Pool[Pooler]{
		Buf: make([]Pooler, 0),
	}
}

// Get - получение элемента из пула.
func (p *Pool[T]) Get() T {
	if len(p.Buf) == 0 {
		return *new(T)
	}
	t := p.Buf[len(p.Buf)-1]
	p.Buf = p.Buf[:len(p.Buf)-1]
	return t
}

// Put - добавление элемента в пул.
func (p *Pool[T]) Put(t T) {
	p.Buf = append(p.Buf, t)
}

// Stats - модель статистики.
type Stats struct {
	CountURL   int `json:"urls"`
	CountUsers int `json:"users"`
}

// NewStats - конструктор модели статистики.
func NewStats(countURLs, countUsers int) *Stats {
	return &Stats{
		CountURL:   countURLs,
		CountUsers: countUsers,
	}
}
