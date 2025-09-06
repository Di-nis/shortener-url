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
	DataBaseDSN     string `env:"DATABASE_DSN"`
}

func (a *Config) Parse() {
	// первый приоритет - из переменных окружения
	_ = godotenv.Load()
	_ = env.Parse(a)

	// второй приоритет - из аргументов командной строки
	var serverAddress, baseURL, fileStoragePath, dataBaseDSN string
	flag.StringVar(&serverAddress, "a", "localhost:8080", "URL")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "Base URL")
	flag.StringVar(&fileStoragePath, "f", "database.log", "File Storage Path")
	flag.StringVar(&dataBaseDSN, "d", "", "Database Source Name")

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

	if a.DataBaseDSN == "" {
		a.DataBaseDSN = dataBaseDSN
	}
	// if a.LogLevel == "" {
	//  a.LogLevel = logLevel
	// }
}
