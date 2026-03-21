package config

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHTTPHost          = "localhost"
	defaultHTTPPort          = 8080
	defaultHTTPReadTimeout   = 5 * time.Second
	defaultDBPath            = "data\\golf_team_manager.sqlite"
	defaultDBAutoMigrate     = true
	defaultJWTTTL            = time.Hour
	envHTTPHost              = "HTTP_HOST"
	envHTTPPort              = "HTTP_PORT"
	envHTTPReadHeaderTimeout = "HTTP_READ_HEADER_TIMEOUT"
	envDBPath                = "DB_PATH"
	envDBAutoMigrate         = "DB_AUTO_MIGRATE"
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
	LineClientID     string
	LineClientSecret string
	LineRedirectURI  string
	FrontendURL      string
	JWTSecret        string
	JWTTTL           time.Duration
}

func LoadFromEnv() (Config, error) {
	if err := loadRootDotEnv(); err != nil {
		return Config{}, err
	}

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

	jwtTTL, err := loadDurationEnv(envJWTTTL, defaultJWTTTL)
	if err != nil {
		return Config{}, err
	}

	authConfig := AuthConfig{
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
	requiredValues := map[string]string{
		envLineClientID:     cfg.LineClientID,
		envLineClientSecret: cfg.LineClientSecret,
		envLineRedirectURI:  cfg.LineRedirectURI,
		envFrontendURL:      cfg.FrontendURL,
		envJWTSecret:        cfg.JWTSecret,
	}

	for key, value := range requiredValues {
		if value == "" {
			return fmt.Errorf("%s is required", key)
		}
	}

	return nil
}

func loadRootDotEnv() error {
	envPath, err := resolveRootDotEnvPath()
	if err != nil {
		return err
	}

	if envPath == "" {
		return nil
	}

	file, err := os.Open(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("open %s: %w", envPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			return fmt.Errorf("parse %s:%d: expected KEY=VALUE", envPath, lineNumber)
		}

		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("parse %s:%d: empty key", envPath, lineNumber)
		}

		if os.Getenv(key) != "" {
			continue
		}

		if err := os.Setenv(key, trimDotEnvValue(value)); err != nil {
			return fmt.Errorf("set %s from %s:%d: %w", key, envPath, lineNumber, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read %s: %w", envPath, err)
	}

	return nil
}

func resolveRootDotEnvPath() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	moduleRoot, found, err := findModuleRoot(workingDir)
	if err != nil {
		return "", err
	}

	if found {
		return filepath.Join(filepath.Dir(moduleRoot), ".env"), nil
	}

	return filepath.Join(workingDir, ".env"), nil
}

func findModuleRoot(startDir string) (string, bool, error) {
	currentDir := filepath.Clean(startDir)

	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir, true, nil
		} else if !os.IsNotExist(err) {
			return "", false, fmt.Errorf("stat go.mod in %s: %w", currentDir, err)
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", false, nil
		}

		currentDir = parentDir
	}
}

func trimDotEnvValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) >= 2 {
		if (trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"') ||
			(trimmed[0] == '\'' && trimmed[len(trimmed)-1] == '\'') {
			return trimmed[1 : len(trimmed)-1]
		}
	}

	return trimmed
}
