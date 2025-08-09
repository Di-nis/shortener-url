package config

import (
	"flag"
)

type Options struct {
	Port    string
	BaseURL string
}

func (a *Options) Parse() {
	flag.StringVar(&a.Port, "a", ":8080", "URL")
	flag.StringVar(&a.BaseURL, "b", "http://localhost:8080", "base URL")

	flag.Parse()
}
