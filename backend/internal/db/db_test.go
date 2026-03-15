package db

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/chgc/golf_team_manager/backend/internal/config"
)

func TestOpenCreatesSQLiteDatabaseFile(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")

	database, err := Open(config.DBConfig{
		Path:        dbPath,
		AutoMigrate: false,
	})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer database.Close()

	if _, err := database.Exec("SELECT 1"); err != nil {
		t.Fatalf("database exec failed: %v", err)
	}
}

func TestRunMigrationsCreatesBaselineSchema(t *testing.T) {
	database := openTestDatabase(t)
	defer database.Close()

	if err := RunMigrations(context.Background(), database); err != nil {
		t.Fatalf("RunMigrations() error = %v", err)
	}

	assertTableExists(t, database, "app_metadata")
	assertTableExists(t, database, "players")
	assertTableExists(t, database, "sessions")
	assertTableExists(t, database, "registrations")
	assertTableExists(t, database, "users")
	assertAppliedMigrationCount(t, database, 4)
}

func TestRunMigrationsIsIdempotent(t *testing.T) {
	database := openTestDatabase(t)
	defer database.Close()

	if err := RunMigrations(context.Background(), database); err != nil {
		t.Fatalf("first RunMigrations() error = %v", err)
	}

	if err := RunMigrations(context.Background(), database); err != nil {
		t.Fatalf("second RunMigrations() error = %v", err)
	}

	assertAppliedMigrationCount(t, database, 4)
}

func openTestDatabase(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	database, err := Open(config.DBConfig{
		Path:        dbPath,
		AutoMigrate: true,
	})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	return database
}

func assertTableExists(t *testing.T, database *sql.DB, tableName string) {
	t.Helper()

	var name string
	if err := database.QueryRow(
		`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`,
		tableName,
	).Scan(&name); err != nil {
		t.Fatalf("table %q does not exist: %v", tableName, err)
	}
}

func assertAppliedMigrationCount(t *testing.T, database *sql.DB, want int) {
	t.Helper()

	var got int
	if err := database.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&got); err != nil {
		t.Fatalf("count schema_migrations failed: %v", err)
	}

	if got != want {
		t.Fatalf("schema_migrations count = %d, want %d", got, want)
	}
}
