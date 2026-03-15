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
	defaultAuthMode          = "dev_stub"
	defaultAuthRole          = "manager"
	defaultAuthDisplayName   = "Demo Manager"
	defaultAuthSubject       = "dev-manager"
	defaultJWTTTL            = time.Hour
	envHTTPHost              = "HTTP_HOST"
	envHTTPPort              = "HTTP_PORT"
	envHTTPReadHeaderTimeout = "HTTP_READ_HEADER_TIMEOUT"
	envDBPath                = "DB_PATH"
	envDBAutoMigrate         = "DB_AUTO_MIGRATE"
	envAuthMode              = "AUTH_MODE"
	envAuthRole              = "AUTH_DEV_DEFAULT_ROLE"
	envAuthDisplayName       = "AUTH_DEV_DEFAULT_NAME"
	envAuthSubject           = "AUTH_DEV_DEFAULT_SUBJECT"
	envAuthUserID            = "AUTH_DEV_DEFAULT_USER_ID"
	envAuthPlayerID          = "AUTH_DEV_DEFAULT_PLAYER_ID"
	envLineClientID          = "LINE_CLIENT_ID"
	envLineClientSecret      = "LINE_CLIENT_SECRET"
	envLineRedirectURI       = "LINE_REDIRECT_URI"
	envFrontendURL           = "FRONTEND_URL"
	envJWTSecret             = "JWT_SECRET"
	envJWTTTL                = "JWT_TTL"
)

type Config struct {
	HTTP HTTPConfig
	DB   DBConfig
	Auth AuthConfig
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

type AuthConfig struct {
	Mode             string
	DevDisplayName   string
	DevPlayerID      string
	DevRole          string
	DevSubject       string
	DevUserID        string
	LineClientID     string
	LineClientSecret string
	LineRedirectURI  string
	FrontendURL      string
	JWTSecret        string
	JWTTTL           time.Duration
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

	authMode := loadStringEnv(envAuthMode, defaultAuthMode)
	if authMode != "dev_stub" && authMode != "line" {
		return Config{}, fmt.Errorf("%s must be dev_stub or line", envAuthMode)
	}

	authRole := loadStringEnv(envAuthRole, defaultAuthRole)
	if authRole != "manager" && authRole != "player" {
		return Config{}, fmt.Errorf("%s must be manager or player", envAuthRole)
	}

	devSubject := loadStringEnv(envAuthSubject, defaultAuthSubject)
	jwtTTL, err := loadDurationEnv(envJWTTTL, defaultJWTTTL)
	if err != nil {
		return Config{}, err
	}

	authConfig := AuthConfig{
		Mode:             authMode,
		DevDisplayName:   loadStringEnv(envAuthDisplayName, defaultAuthDisplayName),
		DevPlayerID:      loadStringEnv(envAuthPlayerID, ""),
		DevRole:          authRole,
		DevSubject:       devSubject,
		DevUserID:        loadStringEnv(envAuthUserID, deriveDevelopmentUserID(devSubject)),
		LineClientID:     loadStringEnv(envLineClientID, ""),
		LineClientSecret: loadStringEnv(envLineClientSecret, ""),
		LineRedirectURI:  loadStringEnv(envLineRedirectURI, ""),
		FrontendURL:      loadStringEnv(envFrontendURL, ""),
		JWTSecret:        loadStringEnv(envJWTSecret, ""),
		JWTTTL:           jwtTTL,
	}

	if err := validateAuthConfig(authConfig); err != nil {
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
		Auth: authConfig,
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

func validateAuthConfig(cfg AuthConfig) error {
	if cfg.Mode != "line" {
		return nil
	}

	requiredValues := map[string]string{
		envLineClientID:     cfg.LineClientID,
		envLineClientSecret: cfg.LineClientSecret,
		envLineRedirectURI:  cfg.LineRedirectURI,
		envFrontendURL:      cfg.FrontendURL,
		envJWTSecret:        cfg.JWTSecret,
	}

	for key, value := range requiredValues {
		if value == "" {
			return fmt.Errorf("%s is required when %s=line", key, envAuthMode)
		}
	}

	return nil
}

func deriveDevelopmentUserID(subject string) string {
	return "dev-user:" + subject
}
