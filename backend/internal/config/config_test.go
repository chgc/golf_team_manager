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
}

func TestLoadFromEnvUsesOverrides(t *testing.T) {
	t.Setenv(envHTTPHost, "0.0.0.0")
	t.Setenv(envHTTPPort, "9090")
	t.Setenv(envHTTPReadHeaderTimeout, "3s")
	t.Setenv(envDBPath, "data\\test.sqlite")
	t.Setenv(envDBAutoMigrate, "false")

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
