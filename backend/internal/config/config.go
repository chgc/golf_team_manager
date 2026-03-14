package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	defaultHTTPHost          = "127.0.0.1"
	defaultHTTPPort          = 8080
	defaultHTTPReadTimeout   = 5 * time.Second
	defaultDBPath            = "data\\golf_team_manager.sqlite"
	defaultDBAutoMigrate     = true
	envHTTPHost              = "HTTP_HOST"
	envHTTPPort              = "HTTP_PORT"
	envHTTPReadHeaderTimeout = "HTTP_READ_HEADER_TIMEOUT"
	envDBPath                = "DB_PATH"
	envDBAutoMigrate         = "DB_AUTO_MIGRATE"
)

type Config struct {
	HTTP HTTPConfig
	DB   DBConfig
}

type HTTPConfig struct {
	Host        string
	Port        int
	ReadTimeout time.Duration
}

type DBConfig struct {
	Path        string
	AutoMigrate bool
}

func LoadFromEnv() (Config, error) {
	port, err := loadIntEnv(envHTTPPort, defaultHTTPPort)
	if err != nil {
		return Config{}, err
	}

	readTimeout, err := loadDurationEnv(envHTTPReadHeaderTimeout, defaultHTTPReadTimeout)
	if err != nil {
		return Config{}, err
	}

	autoMigrate, err := loadBoolEnv(envDBAutoMigrate, defaultDBAutoMigrate)
	if err != nil {
		return Config{}, err
	}

	return Config{
		HTTP: HTTPConfig{
			Host:        loadStringEnv(envHTTPHost, defaultHTTPHost),
			Port:        port,
			ReadTimeout: readTimeout,
		},
		DB: DBConfig{
			Path:        filepath.Clean(loadStringEnv(envDBPath, defaultDBPath)),
			AutoMigrate: autoMigrate,
		},
	}, nil
}

func (c HTTPConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func loadStringEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func loadIntEnv(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}

	if parsedValue < 1 || parsedValue > 65535 {
		return 0, fmt.Errorf("%s must be between 1 and 65535", key)
	}

	return parsedValue, nil
}

func loadDurationEnv(key string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsedValue, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}

	if parsedValue <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", key)
	}

	return parsedValue, nil
}

func loadBoolEnv(key string, fallback bool) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("parse %s: %w", key, err)
	}

	return parsedValue, nil
}
