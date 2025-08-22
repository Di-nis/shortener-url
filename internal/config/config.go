package config

import (
	"flag"
	"github.com/joho/godotenv"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func (a *Config) Parse() {
	_ = godotenv.Load()

	// первый приоритет - из переменных окружения
	_ = env.Parse(a)

	if a.ServerAddress != "" && a.BaseURL != "" {
		return
	}
	// второй приоритет - из аргументов командной строки
	flag.StringVar(&a.ServerAddress, "a", "localhost:8080", "URL")
	flag.StringVar(&a.BaseURL, "b", "http://localhost:8080", "base URL")

	flag.Parse()
}
