package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	defaultHTTPHost          = "127.0.0.1"
	defaultHTTPPort          = 8080
	defaultHTTPReadTimeout   = 5 * time.Second
	envHTTPHost              = "HTTP_HOST"
	envHTTPPort              = "HTTP_PORT"
	envHTTPReadHeaderTimeout = "HTTP_READ_HEADER_TIMEOUT"
)

type Config struct {
	HTTP HTTPConfig
}

type HTTPConfig struct {
	Host        string
	Port        int
	ReadTimeout time.Duration
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

	return Config{
		HTTP: HTTPConfig{
			Host:        loadStringEnv(envHTTPHost, defaultHTTPHost),
			Port:        port,
			ReadTimeout: readTimeout,
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
