package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID int
}

type URL struct {
	UUID        string `db:"user_id"`
	Short       string `db:"short"`
	Original    string `db:"original"`
	URLID       string
	DeletedFlag bool `db:"is_deleted"`
}

func (url URL) MarshalJSON() ([]byte, error) {
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

func (url *URL) UnmarshalJSON(data []byte) error {
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

type URLCopyOne struct {
	UUID        string
	Short       string
	Original    string
	URLID       string
	DeletedFlag bool
}

func (url URLCopyOne) MarshalJSON() ([]byte, error) {
	urlAlias := struct {
		Short string `json:"result"`
	}{
		Short: url.Short,
	}

	return json.Marshal(urlAlias)
}

func (url *URLCopyOne) UnmarshalJSON(data []byte) error {
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

type URLCopyTwo struct {
	UUID        string `json:"uuid"`
	Short       string `json:"url_short"`
	Original    string `json:"url_original"`
	URLID       string `json:"-"`
	DeletedFlag bool   `json:"-"`
}

type URLCopyThree struct {
	UUID        string `json:"-"`
	Short       string `json:"url_short"`
	Original    string `json:"url_original"`
	URLID       string `json:"-"`
	DeletedFlag bool   `json:"-"`
}

type URLCopyFour struct {
	UUID        string `json:"-"`
	Short       string `json:"short_url"`
	Original    string `json:"original_url"`
	URLID       string `json:"-"`
	DeletedFlag bool   `json:"-"`
}

// Audit - структура для хранения данных аудита.
type Audit struct {
	TS     int64  `json:"ts"`
	Action string `json:"action"`
	UserID string `json:"user_id"`
	URL    string `json:"url"`
}

// NewAudit - функция для создания нового экземпляра Audit.
func NewAudit(action, userID, url string) *Audit {
	return &Audit{
		TS:     time.Now().Unix(),
		Action: action,
		UserID: userID,
		URL:    url,
	}
}
