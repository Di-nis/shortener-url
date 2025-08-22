package config

import (
	"flag"
	"github.com/joho/godotenv"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	LogLevel      string `env:"LOG_LEVEL"`
}

func (a *Config) Parse() {
	_ = godotenv.Load()

	// первый приоритет - из переменных окружения
	_ = env.Parse(a)

	if a.ServerAddress == "" {
		flag.StringVar(&a.ServerAddress, "a", "localhost:8080", "URL")
	}
	if a.BaseURL == "" {
		flag.StringVar(&a.BaseURL, "b", "http://localhost:8080", "base URL")
	}
	if a.LogLevel == "" {
		flag.StringVar(&a.LogLevel, "c", "info", "log level")
	}

	flag.Parse()
}
