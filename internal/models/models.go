package models

type Request struct {
	URLOriginal string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}
