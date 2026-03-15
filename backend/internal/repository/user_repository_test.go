package repository

import (
	"context"
	"database/sql"
	"errors"
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

func TestSQLiteUserRepositoryListFiltersByLinkStateAndRole(t *testing.T) {
	database := openUserRepositoryTestDatabase(t)
	repository := NewSQLiteUserRepository(database)

	seedPlayer(t, database, "player-1")
	seedUser(
		t,
		database,
		"user-manager-linked",
		sql.NullString{String: "player-1", Valid: true},
		"經理",
		auth.RoleManager,
		"line",
		"line-manager",
	)
	seedUser(
		t,
		database,
		"user-player-unlinked",
		sql.NullString{},
		"球員",
		auth.RolePlayer,
		"line",
		"line-player",
	)

	users, err := repository.List(context.Background(), UserListFilter{
		LinkState: UserLinkStateUnlinked,
		Role:      auth.RolePlayer,
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("len(users) = %d, want 1", len(users))
	}

	if users[0].ID != "user-player-unlinked" {
		t.Fatalf("users[0].ID = %q, want %q", users[0].ID, "user-player-unlinked")
	}
}

func TestSQLiteUserRepositoryCountByRole(t *testing.T) {
	database := openUserRepositoryTestDatabase(t)
	repository := NewSQLiteUserRepository(database)

	seedUser(t, database, "user-manager-1", sql.NullString{}, "經理一", auth.RoleManager, "line", "line-manager-1")
	seedUser(t, database, "user-manager-2", sql.NullString{}, "經理二", auth.RoleManager, "line", "line-manager-2")
	seedUser(t, database, "user-player-1", sql.NullString{}, "球員一", auth.RolePlayer, "line", "line-player-1")

	count, err := repository.CountByRole(context.Background(), auth.RoleManager)
	if err != nil {
		t.Fatalf("CountByRole() error = %v", err)
	}

	if count != 2 {
		t.Fatalf("count = %d, want 2", count)
	}
}

func TestSQLiteUserRepositoryUpdateRoleAndPlayerReturnsConflictWhenPlayerAlreadyLinked(t *testing.T) {
	database := openUserRepositoryTestDatabase(t)
	repository := NewSQLiteUserRepository(database)

	seedPlayer(t, database, "player-1")
	seedUser(
		t,
		database,
		"user-1",
		sql.NullString{String: "player-1", Valid: true},
		"既有綁定",
		auth.RolePlayer,
		"line",
		"line-user-1",
	)
	seedUser(t, database, "user-2", sql.NullString{}, "未綁定", auth.RolePlayer, "line", "line-user-2")

	playerID := "player-1"
	_, err := repository.UpdateRoleAndPlayer(context.Background(), "user-2", auth.RolePlayer, &playerID)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("UpdateRoleAndPlayer() error = %v, want %v", err, ErrConflict)
	}
}

func seedPlayer(t *testing.T, database *sql.DB, playerID string) {
	t.Helper()

	_, err := database.ExecContext(
		context.Background(),
		`INSERT INTO players (id, name, handicap, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		playerID,
		"測試球員",
		12.5,
		"active",
	)
	if err != nil {
		t.Fatalf("seed player error = %v", err)
	}
}

func seedUser(
	t *testing.T,
	database *sql.DB,
	userID string,
	playerID sql.NullString,
	displayName string,
	role auth.Role,
	provider auth.Provider,
	subject string,
) {
	t.Helper()

	_, err := database.ExecContext(
		context.Background(),
		`INSERT INTO users (id, player_id, display_name, role, auth_provider, provider_subject, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		userID,
		nullStringValue(playerID),
		displayName,
		role,
		provider,
		subject,
	)
	if err != nil {
		t.Fatalf("seed user error = %v", err)
	}
}

func nullStringValue(value sql.NullString) any {
	if value.Valid {
		return value.String
	}

	return nil
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
