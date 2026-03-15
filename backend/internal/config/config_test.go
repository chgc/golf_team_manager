package config

import (
	"testing"
	"time"
)

func TestLoadFromEnvUsesDefaults(t *testing.T) {
	t.Setenv(envHTTPHost, "")
	t.Setenv(envHTTPPort, "")
	t.Setenv(envHTTPReadHeaderTimeout, "")
	t.Setenv(envDBPath, "")
	t.Setenv(envDBAutoMigrate, "")
	t.Setenv(envAuthMode, "")
	t.Setenv(envAuthRole, "")
	t.Setenv(envAuthDisplayName, "")
	t.Setenv(envAuthSubject, "")
	t.Setenv(envAuthUserID, "")
	t.Setenv(envAuthPlayerID, "")
	t.Setenv(envLineClientID, "")
	t.Setenv(envLineClientSecret, "")
	t.Setenv(envLineRedirectURI, "")
	t.Setenv(envFrontendURL, "")
	t.Setenv(envJWTSecret, "")
	t.Setenv(envJWTTTL, "")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}

	if cfg.HTTP.Host != defaultHTTPHost {
		t.Fatalf("Host = %q, want %q", cfg.HTTP.Host, defaultHTTPHost)
	}

	if cfg.HTTP.Port != defaultHTTPPort {
		t.Fatalf("Port = %d, want %d", cfg.HTTP.Port, defaultHTTPPort)
	}

	if cfg.HTTP.ReadTimeout != defaultHTTPReadTimeout {
		t.Fatalf("ReadTimeout = %s, want %s", cfg.HTTP.ReadTimeout, defaultHTTPReadTimeout)
	}

	if cfg.DB.Path != defaultDBPath {
		t.Fatalf("DB.Path = %q, want %q", cfg.DB.Path, defaultDBPath)
	}

	if cfg.DB.AutoMigrate != defaultDBAutoMigrate {
		t.Fatalf("DB.AutoMigrate = %t, want %t", cfg.DB.AutoMigrate, defaultDBAutoMigrate)
	}

	if cfg.Auth.Mode != defaultAuthMode {
		t.Fatalf("Auth.Mode = %q, want %q", cfg.Auth.Mode, defaultAuthMode)
	}

	if cfg.Auth.DevRole != defaultAuthRole {
		t.Fatalf("Auth.DevRole = %q, want %q", cfg.Auth.DevRole, defaultAuthRole)
	}

	if cfg.Auth.DevUserID != deriveDevelopmentUserID(defaultAuthSubject) {
		t.Fatalf("Auth.DevUserID = %q, want %q", cfg.Auth.DevUserID, deriveDevelopmentUserID(defaultAuthSubject))
	}

	if cfg.Auth.JWTTTL != defaultJWTTTL {
		t.Fatalf("Auth.JWTTTL = %s, want %s", cfg.Auth.JWTTTL, defaultJWTTTL)
	}
}

func TestLoadFromEnvUsesOverrides(t *testing.T) {
	t.Setenv(envHTTPHost, "0.0.0.0")
	t.Setenv(envHTTPPort, "9090")
	t.Setenv(envHTTPReadHeaderTimeout, "3s")
	t.Setenv(envDBPath, "data\\test.sqlite")
	t.Setenv(envDBAutoMigrate, "false")
	t.Setenv(envAuthMode, "dev_stub")
	t.Setenv(envAuthRole, "player")
	t.Setenv(envAuthDisplayName, "Demo Player")
	t.Setenv(envAuthSubject, "dev-player")
	t.Setenv(envAuthUserID, "dev-user-player")
	t.Setenv(envAuthPlayerID, "player-1")
	t.Setenv(envJWTTTL, "30m")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}

	if cfg.HTTP.Host != "0.0.0.0" {
		t.Fatalf("Host = %q, want %q", cfg.HTTP.Host, "0.0.0.0")
	}

	if cfg.HTTP.Port != 9090 {
		t.Fatalf("Port = %d, want %d", cfg.HTTP.Port, 9090)
	}

	if cfg.HTTP.ReadTimeout != 3*time.Second {
		t.Fatalf("ReadTimeout = %s, want %s", cfg.HTTP.ReadTimeout, 3*time.Second)
	}

	if cfg.DB.Path != "data\\test.sqlite" {
		t.Fatalf("DB.Path = %q, want %q", cfg.DB.Path, "data\\test.sqlite")
	}

	if cfg.DB.AutoMigrate {
		t.Fatal("DB.AutoMigrate = true, want false")
	}

	if cfg.Auth.DevRole != "player" {
		t.Fatalf("Auth.DevRole = %q, want %q", cfg.Auth.DevRole, "player")
	}

	if cfg.Auth.DevDisplayName != "Demo Player" {
		t.Fatalf("Auth.DevDisplayName = %q, want %q", cfg.Auth.DevDisplayName, "Demo Player")
	}

	if cfg.Auth.DevSubject != "dev-player" {
		t.Fatalf("Auth.DevSubject = %q, want %q", cfg.Auth.DevSubject, "dev-player")
	}

	if cfg.Auth.DevPlayerID != "player-1" {
		t.Fatalf("Auth.DevPlayerID = %q, want %q", cfg.Auth.DevPlayerID, "player-1")
	}

	if cfg.Auth.DevUserID != "dev-user-player" {
		t.Fatalf("Auth.DevUserID = %q, want %q", cfg.Auth.DevUserID, "dev-user-player")
	}

	if cfg.Auth.JWTTTL != 30*time.Minute {
		t.Fatalf("Auth.JWTTTL = %s, want %s", cfg.Auth.JWTTTL, 30*time.Minute)
	}
}

func TestLoadFromEnvRejectsInvalidPort(t *testing.T) {
	t.Setenv(envHTTPPort, "70000")

	_, err := LoadFromEnv()
	if err == nil {
		t.Fatal("LoadFromEnv() error = nil, want error")
	}
}

func TestLoadFromEnvRejectsInvalidAutoMigrateValue(t *testing.T) {
	t.Setenv(envDBAutoMigrate, "not-a-bool")

	_, err := LoadFromEnv()
	if err == nil {
		t.Fatal("LoadFromEnv() error = nil, want error")
	}
}

func TestLoadFromEnvRejectsInvalidAuthRole(t *testing.T) {
	t.Setenv(envAuthRole, "guest")

	_, err := LoadFromEnv()
	if err == nil {
		t.Fatal("LoadFromEnv() error = nil, want error")
	}
}

func TestLoadFromEnvRejectsInvalidAuthMode(t *testing.T) {
	t.Setenv(envAuthMode, "unsupported")

	_, err := LoadFromEnv()
	if err == nil {
		t.Fatal("LoadFromEnv() error = nil, want error")
	}
}

func TestLoadFromEnvRejectsMissingLineConfig(t *testing.T) {
	t.Setenv(envAuthMode, "line")
	t.Setenv(envLineClientID, "line-client")
	t.Setenv(envLineClientSecret, "line-secret")
	t.Setenv(envLineRedirectURI, "http://localhost:8080/api/auth/line/callback")
	t.Setenv(envFrontendURL, "http://localhost:4200")
	t.Setenv(envJWTSecret, "")

	_, err := LoadFromEnv()
	if err == nil {
		t.Fatal("LoadFromEnv() error = nil, want error")
	}
}

func TestLoadFromEnvLoadsLineConfig(t *testing.T) {
	t.Setenv(envAuthMode, "line")
	t.Setenv(envLineClientID, "line-client")
	t.Setenv(envLineClientSecret, "line-secret")
	t.Setenv(envLineRedirectURI, "http://localhost:8080/api/auth/line/callback")
	t.Setenv(envFrontendURL, "http://localhost:4200")
	t.Setenv(envJWTSecret, "jwt-secret")
	t.Setenv(envJWTTTL, "45m")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}

	if cfg.Auth.Mode != "line" {
		t.Fatalf("Auth.Mode = %q, want %q", cfg.Auth.Mode, "line")
	}

	if cfg.Auth.LineClientID != "line-client" {
		t.Fatalf("Auth.LineClientID = %q, want %q", cfg.Auth.LineClientID, "line-client")
	}

	if cfg.Auth.LineClientSecret != "line-secret" {
		t.Fatalf("Auth.LineClientSecret = %q, want %q", cfg.Auth.LineClientSecret, "line-secret")
	}

	if cfg.Auth.LineRedirectURI != "http://localhost:8080/api/auth/line/callback" {
		t.Fatalf("Auth.LineRedirectURI = %q, want %q", cfg.Auth.LineRedirectURI, "http://localhost:8080/api/auth/line/callback")
	}

	if cfg.Auth.FrontendURL != "http://localhost:4200" {
		t.Fatalf("Auth.FrontendURL = %q, want %q", cfg.Auth.FrontendURL, "http://localhost:4200")
	}

	if cfg.Auth.JWTSecret != "jwt-secret" {
		t.Fatalf("Auth.JWTSecret = %q, want %q", cfg.Auth.JWTSecret, "jwt-secret")
	}

	if cfg.Auth.JWTTTL != 45*time.Minute {
		t.Fatalf("Auth.JWTTTL = %s, want %s", cfg.Auth.JWTTTL, 45*time.Minute)
	}
}
