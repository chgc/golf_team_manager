package repository

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/chgc/golf_team_manager/backend/internal/config"
	appdb "github.com/chgc/golf_team_manager/backend/internal/db"
)

func TestSQLiteUserRepositoryUpsertLineUserCreatesUnlinkedUser(t *testing.T) {
	database := openUserRepositoryTestDatabase(t)
	repository := NewSQLiteUserRepository(database)

	user, err := repository.UpsertLineUser(context.Background(), auth.UpsertLineUserInput{
		DisplayName: "王小明",
		Subject:     "line-user-1",
	})
	if err != nil {
		t.Fatalf("UpsertLineUser() error = %v", err)
	}

	if user.Role != auth.RolePlayer {
		t.Fatalf("user.Role = %q, want %q", user.Role, auth.RolePlayer)
	}

	if user.PlayerID != "" {
		t.Fatalf("user.PlayerID = %q, want empty", user.PlayerID)
	}
}

func TestSQLiteUserRepositoryUpsertLineUserPreservesExistingLinkAndRole(t *testing.T) {
	database := openUserRepositoryTestDatabase(t)
	repository := NewSQLiteUserRepository(database)

	var err error
	_, err = database.ExecContext(
		context.Background(),
		`INSERT INTO players (id, name, handicap, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		"player-1",
		"測試球員",
		12.5,
		"active",
	)
	if err != nil {
		t.Fatalf("seed player error = %v", err)
	}

	_, err = database.ExecContext(
		context.Background(),
		`INSERT INTO users (id, player_id, display_name, role, auth_provider, provider_subject, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		"user-1",
		"player-1",
		"舊名字",
		"manager",
		"line",
		"line-user-1",
	)
	if err != nil {
		t.Fatalf("seed user error = %v", err)
	}

	user, err := repository.UpsertLineUser(context.Background(), auth.UpsertLineUserInput{
		DisplayName: "新名字",
		Subject:     "line-user-1",
	})
	if err != nil {
		t.Fatalf("UpsertLineUser() error = %v", err)
	}

	if user.ID != "user-1" {
		t.Fatalf("user.ID = %q, want %q", user.ID, "user-1")
	}

	if user.PlayerID != "player-1" {
		t.Fatalf("user.PlayerID = %q, want %q", user.PlayerID, "player-1")
	}

	if user.Role != auth.RoleManager {
		t.Fatalf("user.Role = %q, want %q", user.Role, auth.RoleManager)
	}

	if user.DisplayName != "新名字" {
		t.Fatalf("user.DisplayName = %q, want %q", user.DisplayName, "新名字")
	}
}

func openUserRepositoryTestDatabase(t *testing.T) *sql.DB {
	t.Helper()

	database, err := appdb.Open(config.DBConfig{
		Path:        filepath.Join(t.TempDir(), "user-repository.sqlite"),
		AutoMigrate: true,
	})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	if err := appdb.RunMigrations(context.Background(), database); err != nil {
		t.Fatalf("RunMigrations() error = %v", err)
	}

	t.Cleanup(func() {
		database.Close()
	})

	return database
}
