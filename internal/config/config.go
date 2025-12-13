// Package config содержит структуры и функции для загрузки и проверки
// конфигурации приложения из переменных окружения, файлов или других источников.
package config

import (
	"flag"

	"github.com/joho/godotenv"

	"github.com/caarlos0/env/v6"
)

// Config - структура конфигурации приложения.
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DataBaseDSN     string `env:"DATABASE_DSN"`
	JWTSecret       string `env:"JWT_SECRET"`
	AuditFile       string `env:"AUDIT_FILE"`
	AuditURL        string `env:"AUDIT_URL"`
	UseMockAuth     bool
	EnableHTTPS     bool   `env:"ENABLE_HTTPS"`
	CertFilePath    string `env:"CERT_FILE_PATH"`
	KeyFilePath     string `env:"KEY_FILE_PATH"`
}

// NewConfig - функция для создания конфигурации.
func NewConfig() *Config {
	return &Config{}
}

// Load - метод для парсинга конфигурации из переменных окружения и аргументов командной строки.
func (a *Config) Load() {
	// первый приоритет - из переменных окружения
	_ = godotenv.Load()
	_ = env.Parse(a)

	// второй приоритет - из аргументов командной строки
	var (
		serverAddress, baseURL, fileStoragePath string
		dataBaseDSN, auditFile, auditURL        string
		enableHTTPS                             bool
	)
	flag.StringVar(&serverAddress, "a", "localhost:8080", "URL")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "Base URL")
	flag.StringVar(&fileStoragePath, "f", "db.log", "File Storage Path")
	flag.StringVar(&dataBaseDSN, "d", "", "Database Source Name")
	flag.StringVar(&auditFile, "audit-file", "", "Audit File Path")
	flag.StringVar(&auditURL, "audit-url", "", "Audit URL Path")
	flag.BoolVar(&enableHTTPS, "s", false, "Use HTTPS web-server")

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
	if a.AuditFile == "" {
		a.AuditFile = auditFile
	}
	if a.AuditURL == "" {
		a.AuditURL = auditURL
	}
	if !a.EnableHTTPS {
		a.EnableHTTPS = enableHTTPS
	}
}
