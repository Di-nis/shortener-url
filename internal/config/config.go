package config

import (
	"flag"
	"github.com/joho/godotenv"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func (a *Config) Parse() {
	// первый приоритет - из переменных окружения
	_ = godotenv.Load()
	_ = env.Parse(a)

	// второй приоритет - из аргументов командной строки
	var serverAddress, baseURL, fileStoragePath string
	flag.StringVar(&serverAddress, "a", "localhost:8080", "URL")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "base URL")
	flag.StringVar(&fileStoragePath, "f", "database.log", "File storage path")

	flag.Parse()

	if a.ServerAddress == "" {
		a.ServerAddress = serverAddress
	}
	if a.BaseURL == "" {
		a.BaseURL = baseURL
	}
	if a.FileStoragePath == "" {
		a.FileStoragePath = fileStoragePath
	}
	// if a.LogLevel == "" {
	//  a.LogLevel = logLevel
	// }
}
