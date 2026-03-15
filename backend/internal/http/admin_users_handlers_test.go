package apihttp

import (
	"context"
	"database/sql"
	nethttp "net/http"
	"testing"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/chgc/golf_team_manager/backend/internal/config"
	"github.com/gin-gonic/gin"
)

func TestListAdminUsersRequiresManagerRole(t *testing.T) {
	router, cleanup := newTestRouter(t)
	defer cleanup()

	response := performJSONRequestWithHeaders(
		t,
		router,
		nethttp.MethodGet,
		"/api/admin/users",
		nil,
		map[string]string{
			"X-Debug-Role": "player",
		},
	)
	if response.Code != nethttp.StatusForbidden {
		t.Fatalf("status code = %d, want %d", response.Code, nethttp.StatusForbidden)
	}

	var errorResponse ErrorResponse
	decodeResponseBody(t, response, &errorResponse)
	if errorResponse.Error.Code != "forbidden" {
		t.Fatalf("error code = %q, want %q", errorResponse.Error.Code, "forbidden")
	}
}

func TestListAdminUsersReturnsFilteredUsers(t *testing.T) {
	router, cleanup, database := newTestRouterWithDatabase(t)
	defer cleanup()

	seedAdminPlayer(t, database, "player-1")
	seedAdminUser(t, database, "user-manager", sql.NullString{String: "player-1", Valid: true}, "Manager", auth.RoleManager, "manager-subject")
	seedAdminUser(t, database, "user-player", sql.NullString{}, "Player", auth.RolePlayer, "player-subject")

	response := performJSONRequest(t, router, nethttp.MethodGet, "/api/admin/users?linkState=unlinked&role=player", nil)
	if response.Code != nethttp.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, nethttp.StatusOK)
	}

	var users []adminUserResponse
	decodeResponseBody(t, response, &users)
	if len(users) != 1 {
		t.Fatalf("len(users) = %d, want 1", len(users))
	}

	if users[0].UserID != "user-player" {
		t.Fatalf("users[0].UserID = %q, want %q", users[0].UserID, "user-player")
	}
}

func TestUpdateAdminUserPromotesAndLinksUser(t *testing.T) {
	router, cleanup, database := newTestRouterWithDatabase(t)
	defer cleanup()

	seedAdminPlayer(t, database, "player-1")
	seedAdminUser(t, database, "user-player", sql.NullString{}, "Player", auth.RolePlayer, "player-subject")

	response := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/admin/users/user-player",
		map[string]any{
			"role":     "manager",
			"playerId": "player-1",
		},
	)
	if response.Code != nethttp.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, nethttp.StatusOK)
	}

	var user adminUserResponse
	decodeResponseBody(t, response, &user)
	if user.Role != auth.RoleManager {
		t.Fatalf("user.Role = %q, want %q", user.Role, auth.RoleManager)
	}

	if user.PlayerID != "player-1" {
		t.Fatalf("user.PlayerID = %q, want %q", user.PlayerID, "player-1")
	}
}

func TestUpdateAdminUserAllowsNullPlayerIDToUnlink(t *testing.T) {
	router, cleanup, database := newTestRouterWithDatabase(t)
	defer cleanup()

	seedAdminPlayer(t, database, "player-1")
	seedAdminUser(
		t,
		database,
		"user-manager",
		sql.NullString{},
		"Manager",
		auth.RoleManager,
		"manager-subject",
	)
	seedAdminUser(
		t,
		database,
		"user-player",
		sql.NullString{String: "player-1", Valid: true},
		"Player",
		auth.RolePlayer,
		"player-subject",
	)

	response := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/admin/users/user-player",
		map[string]any{
			"playerId": nil,
		},
	)
	if response.Code != nethttp.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, nethttp.StatusOK)
	}

	var user adminUserResponse
	decodeResponseBody(t, response, &user)
	if user.PlayerID != "" {
		t.Fatalf("user.PlayerID = %q, want empty", user.PlayerID)
	}
}

func TestUpdateAdminUserRejectsLastManagerDemotion(t *testing.T) {
	router, cleanup, database := newTestRouterWithDatabase(t)
	defer cleanup()

	seedAdminUser(t, database, "user-manager", sql.NullString{}, "Manager", auth.RoleManager, "manager-subject")

	response := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/admin/users/user-manager",
		map[string]any{
			"role": "player",
		},
	)
	if response.Code != nethttp.StatusConflict {
		t.Fatalf("status code = %d, want %d", response.Code, nethttp.StatusConflict)
	}

	var errorResponse ErrorResponse
	decodeResponseBody(t, response, &errorResponse)
	if errorResponse.Error.Code != "last_manager_demotion_forbidden" {
		t.Fatalf("error code = %q, want %q", errorResponse.Error.Code, "last_manager_demotion_forbidden")
	}
}

func newTestRouterWithDatabase(t *testing.T) (*gin.Engine, func(), *sql.DB) {
	t.Helper()

	database := openTestDatabase(t)
	cleanup := func() {
		database.Close()
	}

	testConfig, err := config.LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}

	return NewRouter(database, testConfig), cleanup, database
}

func seedAdminPlayer(t *testing.T, database *sql.DB, playerID string) {
	t.Helper()

	_, err := database.ExecContext(
		context.Background(),
		`INSERT INTO players (id, name, handicap, phone, email, status, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		playerID,
		"Test Player",
		12.5,
		"",
		"",
		"active",
		"",
	)
	if err != nil {
		t.Fatalf("seed player error = %v", err)
	}
}

func seedAdminUser(
	t *testing.T,
	database *sql.DB,
	userID string,
	playerID sql.NullString,
	displayName string,
	role auth.Role,
	subject string,
) {
	t.Helper()

	var nullablePlayerID any
	if playerID.Valid {
		nullablePlayerID = playerID.String
	}

	_, err := database.ExecContext(
		context.Background(),
		`INSERT INTO users (id, player_id, display_name, role, auth_provider, provider_subject, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		userID,
		nullablePlayerID,
		displayName,
		role,
		auth.ProviderLINEOAuth,
		subject,
	)
	if err != nil {
		t.Fatalf("seed user error = %v", err)
	}
}
