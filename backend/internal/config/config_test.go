package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func isolateWorkingDir(t *testing.T) string {
	t.Helper()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir(%q) error = %v", tempDir, err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	return tempDir
}

func TestLoadFromEnvUsesDefaults(t *testing.T) {
	isolateWorkingDir(t)
	t.Setenv(envHTTPHost, "")
	t.Setenv(envHTTPPort, "")
	t.Setenv(envHTTPReadHeaderTimeout, "")
	t.Setenv(envDBPath, "")
	t.Setenv(envDBAutoMigrate, "")
	t.Setenv(envLineClientID, "")
	t.Setenv(envLineClientSecret, "")
	t.Setenv(envLineRedirectURI, "")
	t.Setenv(envFrontendURL, "")
	t.Setenv(envJWTSecret, "")
	t.Setenv(envJWTTTL, "")

	_, err := LoadFromEnv()
	if err == nil {
		t.Fatal("LoadFromEnv() error = nil, want error")
	}
	if !strings.HasSuffix(err.Error(), " is required") {
		t.Fatalf("error = %q, want suffix %q", err.Error(), " is required")
	}
}

func TestLoadFromEnvUsesOverrides(t *testing.T) {
	isolateWorkingDir(t)
	t.Setenv(envHTTPHost, "0.0.0.0")
	t.Setenv(envHTTPPort, "9090")
	t.Setenv(envHTTPReadHeaderTimeout, "3s")
	t.Setenv(envDBPath, "data\\test.sqlite")
	t.Setenv(envDBAutoMigrate, "false")
	t.Setenv(envLineClientID, "line-client")
	t.Setenv(envLineClientSecret, "line-secret")
	t.Setenv(envLineRedirectURI, "http://localhost:8080/api/auth/line/callback")
	t.Setenv(envFrontendURL, "http://localhost:4200")
	t.Setenv(envJWTSecret, "jwt-secret")
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

	if cfg.Auth.JWTTTL != 30*time.Minute {
		t.Fatalf("Auth.JWTTTL = %s, want %s", cfg.Auth.JWTTTL, 30*time.Minute)
	}
}

func TestLoadFromEnvRejectsInvalidPort(t *testing.T) {
	isolateWorkingDir(t)
	t.Setenv(envHTTPPort, "70000")

	_, err := LoadFromEnv()
	if err == nil {
		t.Fatal("LoadFromEnv() error = nil, want error")
	}
}

func TestLoadFromEnvRejectsInvalidAutoMigrateValue(t *testing.T) {
	isolateWorkingDir(t)
	t.Setenv(envDBAutoMigrate, "not-a-bool")

	_, err := LoadFromEnv()
	if err == nil {
		t.Fatal("LoadFromEnv() error = nil, want error")
	}
}

func TestLoadFromEnvRejectsMissingLineConfig(t *testing.T) {
	isolateWorkingDir(t)
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
	isolateWorkingDir(t)
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

func TestLoadFromEnvLoadsRootDotEnv(t *testing.T) {
	tempDir := t.TempDir()
	backendDir := filepath.Join(tempDir, "backend")
	if err := os.MkdirAll(backendDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", backendDir, err)
	}
	if err := os.WriteFile(filepath.Join(backendDir, "go.mod"), []byte("module example.com/test\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(go.mod) error = %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(tempDir, ".env"),
		[]byte("LINE_CLIENT_ID=line-client\nLINE_CLIENT_SECRET=line-secret\nLINE_REDIRECT_URI=http://localhost:8080/api/auth/line/callback\nFRONTEND_URL=http://localhost:4200\nJWT_SECRET=jwt-secret\nJWT_TTL=45m\n"),
		0o644,
	); err != nil {
		t.Fatalf("WriteFile(.env) error = %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	if err := os.Chdir(backendDir); err != nil {
		t.Fatalf("Chdir(%q) error = %v", backendDir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

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

	if cfg.Auth.LineClientID != "line-client" {
		t.Fatalf("Auth.LineClientID = %q, want %q", cfg.Auth.LineClientID, "line-client")
	}

	if cfg.Auth.JWTSecret != "jwt-secret" {
		t.Fatalf("Auth.JWTSecret = %q, want %q", cfg.Auth.JWTSecret, "jwt-secret")
	}
}

func TestLoadFromEnvPrefersProcessEnvOverRootDotEnv(t *testing.T) {
	tempDir := t.TempDir()
	backendDir := filepath.Join(tempDir, "backend")
	if err := os.MkdirAll(backendDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", backendDir, err)
	}
	if err := os.WriteFile(filepath.Join(backendDir, "go.mod"), []byte("module example.com/test\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(go.mod) error = %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(tempDir, ".env"),
		[]byte("LINE_CLIENT_ID=line-client\nLINE_CLIENT_SECRET=line-secret\nLINE_REDIRECT_URI=http://localhost:8080/api/auth/line/callback\nFRONTEND_URL=http://localhost:4200\nJWT_SECRET=jwt-secret\n"),
		0o644,
	); err != nil {
		t.Fatalf("WriteFile(.env) error = %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	if err := os.Chdir(backendDir); err != nil {
		t.Fatalf("Chdir(%q) error = %v", backendDir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	t.Setenv(envLineClientID, "shell-client")
	t.Setenv(envLineClientSecret, "line-secret")
	t.Setenv(envLineRedirectURI, "http://localhost:8080/api/auth/line/callback")
	t.Setenv(envFrontendURL, "http://localhost:4200")
	t.Setenv(envJWTSecret, "jwt-secret")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}

	if cfg.Auth.LineClientID != "shell-client" {
		t.Fatalf("Auth.LineClientID = %q, want %q", cfg.Auth.LineClientID, "shell-client")
	}
}
