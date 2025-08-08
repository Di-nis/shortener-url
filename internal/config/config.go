package config

import (
	"flag"
)

const (
	Port = "8080"
)

type Run struct {
	URL     string
	BaseURL string
}

func (a *Run) ParseOptions() {
	flag.StringVar(&a.URL, "a", "localhost:8080", "URL")
	flag.StringVar(&a.BaseURL, "b", "http://localhost:8080", "base URL")

	flag.Parse()
}
