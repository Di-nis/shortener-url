package models

import (
	"encoding/json"
)

// User - модель пользователя.
type User struct {
	ID int
}

// URL - модель URL.
type URL struct {
	UUID        string `db:"user_id"`
	Short       string `db:"short"`
	Original    string `db:"original"`
	URLID       string
	DeletedFlag bool `db:"is_deleted"`
}

// MarshalJSON - метод для сериализации модели URL.
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

// UnmarshalJSON - метод для десериализации модели URL.
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

// URLCopyOne - модель URL.
type URLCopyOne struct {
	UUID        string
	Short       string
	Original    string
	URLID       string
	DeletedFlag bool
}

// MarshalJSON - метод для сериализации модели URL.
func (url URLCopyOne) MarshalJSON() ([]byte, error) {
	urlAlias := struct {
		Short string `json:"result"`
	}{
		Short: url.Short,
	}

	return json.Marshal(urlAlias)
}

// UnmarshalJSON - метод для десериализации модели URL.
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

// URLCopyTwo - модель URL.
type URLCopyTwo struct {
	UUID        string `json:"uuid"`
	Short       string `json:"url_short"`
	Original    string `json:"url_original"`
	URLID       string `json:"-"`
	DeletedFlag bool   `json:"-"`
}

// URLCopyThree - модель URL.
type URLCopyThree struct {
	UUID        string `json:"-"`
	Short       string `json:"url_short"`
	Original    string `json:"url_original"`
	URLID       string `json:"-"`
	DeletedFlag bool   `json:"-"`
}

// URLCopyFour - модель URL.
type URLCopyFour struct {
	UUID        string `json:"-"`
	Short       string `json:"short_url"`
	Original    string `json:"original_url"`
	URLID       string `json:"-"`
	DeletedFlag bool   `json:"-"`
}
