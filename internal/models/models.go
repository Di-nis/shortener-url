package models

import "encoding/json"

type URL struct {
	ID          string
	URLShort    string
	URLOriginal string
}

func (url URL) MarshalJSON() ([]byte, error) {
	aliasValue := struct {
		URLShort    string `json:"url_short"`
		URLOriginal string `json:"url_original"`
		ID          string `json:"uuid"`
	}{
		URLShort:    url.URLShort,
		URLOriginal: url.URLOriginal,
		ID:          url.ID,
	}

	return json.Marshal(aliasValue)
}

func (url *URL) UnmarshalJSON(data []byte) (err error) {
	aliasValue := &struct {
		URLShort    string `json:"url_short"`
		URLOriginal string `json:"url_original"`
		ID          string `json:"uuid"`
	}{
		URLShort:    url.URLShort,
		URLOriginal: url.URLOriginal,
		ID:          url.ID,
	}

	if err = json.Unmarshal(data, aliasValue); err != nil {
		return
	}
	return
}

type Request struct {
	URLOriginal string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}
