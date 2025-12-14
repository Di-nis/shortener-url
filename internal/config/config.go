// Package config содержит структуры и функции для загрузки и проверки
// конфигурации приложения из переменных окружения, файлов или других источников.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

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
	Config          string `env:"CONFIG"`
}

// NewConfig - функция для создания конфигурации.
func NewConfig() *Config {
	return &Config{}
}

// Load - метод для парсинга конфигурации из переменных окружения и аргументов командной строки.
func (c *Config) Load() {
	// первый приоритет - из переменных окружения
	c.loanFromEnv()

	// второй приоритет - из аргументов командной строки
	c.loanFromFlags()

	// третий приоритет - из файла
	c.loanFromFile()
}

// loanFromEnv - загрузка конфигурации из переменных окружения.
func (c *Config) loanFromEnv() error {
	var err error
	if err = godotenv.Load(); err != nil {
		return fmt.Errorf("path: internal/config/config.go, func loanFromEnv(), failed to load env: %w", err)
	}
	if err = env.Parse(c); err != nil {
		return fmt.Errorf("path: internal/config/config.go, func loanFromEnv(), failed to parse env: %w", err)
	}
	return nil
}

// loanFromFlags - загрузка конфигурации из аргументов командной строки.
func (c *Config) loanFromFlags() {
	var (
		serverAddress, baseURL, fileStoragePath  string
		dataBaseDSN, auditFile, auditURL, config string
		enableHTTPS                              bool
	)
	flag.StringVar(&serverAddress, "a", "", "URL")
	flag.StringVar(&baseURL, "b", "", "base URL")
	flag.StringVar(&fileStoragePath, "f", "", "path to storage file")
	flag.StringVar(&dataBaseDSN, "d", "", "database dource name")
	flag.StringVar(&auditFile, "audit-file", "", "path to audit file")
	flag.StringVar(&auditURL, "audit-url", "", "path to audit URL")
	flag.StringVar(&config, "config", "", "path to the configuration file")
	flag.BoolVar(&enableHTTPS, "s", false, "use HTTPS web-server")

	flag.Parse()

	if c.ServerAddress == "" {
		c.ServerAddress = serverAddress
	}
	if c.BaseURL == "" {
		c.BaseURL = baseURL
	}
	if c.FileStoragePath == "" {
		c.FileStoragePath = fileStoragePath
	}
	if c.DataBaseDSN == "" {
		c.DataBaseDSN = dataBaseDSN
	}
	if c.AuditFile == "" {
		c.AuditFile = auditFile
	}
	if c.AuditURL == "" {
		c.AuditURL = auditURL
	}
	if !c.EnableHTTPS {
		c.EnableHTTPS = enableHTTPS
	}
	if c.Config == "" {
		c.Config = config
	}
}

// loanFromFile - загрузка конфигурации из файла.
func (c *Config) loanFromFile() error {
	if c.Config == "" {
		return nil
	}

	type ConfigAlias struct {
		ServerAddress   string `json:"server_address"`
		BaseURL         string `json:"base_url"`
		LogLevel        string `json:"log_level"`
		FileStoragePath string `json:"file_storage_path"`
		DataBaseDSN     string `json:"database_dsn"`
		AuditFile       string `json:"audit_file"`
		AuditURL        string `json:"audit_url"`
		EnableHTTPS     bool   `json:"enable_https"`
		CertFilePath    string `json:"cert_file_path"`
		KeyFilePath     string `json:"key_file_path"`
	}

	var configAlias ConfigAlias

	jsonFile, err := os.Open(c.Config)
	if err != nil {
		return fmt.Errorf("path: internal/config/config.go, func loanFromJSON(), failed to open json file: %w", err)
	}

	jsonFileData, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("path: internal/config/config.go, func loanFromJSON(), failed to read json file: %w", err)
	}
	defer jsonFile.Close()

	if err := json.Unmarshal(jsonFileData, &configAlias); err != nil {
		return fmt.Errorf("path: internal/config/config.go, func loanFromJSON(), failed to unmarshal data: %w", err)
	}

	if c.ServerAddress == "" {
		c.ServerAddress = configAlias.ServerAddress
	}

	if c.BaseURL == "" {
		c.BaseURL = configAlias.BaseURL
	}

	if c.LogLevel == "" {
		c.LogLevel = configAlias.LogLevel
	}

	if c.FileStoragePath == "" {
		c.FileStoragePath = configAlias.FileStoragePath
	}

	if c.DataBaseDSN == "" {
		c.DataBaseDSN = configAlias.DataBaseDSN
	}

	if c.AuditFile == "" {
		c.AuditFile = configAlias.AuditFile
	}

	if c.AuditURL == "" {
		c.AuditURL = configAlias.AuditURL
	}

	if !c.EnableHTTPS {
		c.EnableHTTPS = configAlias.EnableHTTPS
	}

	if c.CertFilePath == "" {
		c.CertFilePath = configAlias.CertFilePath
	}

	if c.KeyFilePath == "" {
		c.KeyFilePath = configAlias.KeyFilePath
	}

	return nil
}
