package models

import (
	"encoding/json"
)

type URL struct {
	ID       string
	Short    string
	Original string
}

func (url URL) MarshalJSON() ([]byte, error) {
	urlAlias := struct {
		Short    string `json:"short_url"`
		Original string `json:"-"`
		ID       string `json:"correlation_id"`
	}{
		Short:    url.Short,
		Original: url.Original,
		ID:       url.ID,
	}

	return json.Marshal(urlAlias)
}

func (url *URL) UnmarshalJSON(data []byte) error {
	type URLAlias struct {
		Short    string `json:"-"`
		Original string `json:"original_url"`
		ID       string `json:"correlation_id"`
	}

	var urlAlias URLAlias

	if err := json.Unmarshal(data, &urlAlias); err != nil {
		return err
	}
	url.Original = urlAlias.Original
	url.ID = urlAlias.ID
	return nil
}

type URLCopyOne struct {
	ID       string
	Short    string
	Original string
}

func (url URLCopyOne) MarshalJSON() ([]byte, error) {
	urlAlias := struct {
		Short    string `json:"result"`
	}{
		Short:    url.Short,
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

type Request struct {
	URLOriginal string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}
