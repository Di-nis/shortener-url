package models

type Request struct {
	UrlOriginal string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}


// HTTP/1.1 201 OK
// Content-Type: application/json
// Content-Length: 30
// {
//  "result": "http://localhost:8080/EwHXdJfB"
// } 