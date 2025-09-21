package models

import (
	"encoding/json"
)

type User struct {
	ID int
}

type URL struct {
	UserID   string
	Short    string
	Original string
	URLID    string
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
	UserID   string
	Short    string
	Original string
	URLID    string
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
	UserID   string `json:"uuid"`
	Short    string `json:"url_short"`
	Original string `json:"url_original"`
	URLID    string `json:"-"`
}

type URLCopyThree struct {
	UserID   string `json:"-"`
	Short    string `json:"url_short"`
	Original string `json:"url_original"`
	URLID    string `json:"-"` 
}

type URLCopyFour struct {
	UserID   string `json:"-"`
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
	URLID    string `json:"-"`
}
