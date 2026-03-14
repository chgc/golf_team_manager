package apihttp

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/chgc/golf_team_manager/backend/internal/config"
	appdb "github.com/chgc/golf_team_manager/backend/internal/db"
	"github.com/chgc/golf_team_manager/backend/internal/http/handlers"
	"github.com/chgc/golf_team_manager/backend/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

func TestNewRouterHealthEndpoint(t *testing.T) {
	router, cleanup := newTestRouter(t)
	defer cleanup()

	request := httptest.NewRequest(nethttp.MethodGet, "/health", nil)
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != nethttp.StatusOK {
		t.Fatalf("status code = %d, want %d", responseRecorder.Code, nethttp.StatusOK)
	}

	var response handlers.HealthResponse
	if err := json.Unmarshal(responseRecorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if response.Status != "ok" {
		t.Fatalf("response status = %q, want %q", response.Status, "ok")
	}

	if responseRecorder.Header().Get(middleware.RequestIDHeader) == "" {
		t.Fatal("request ID header is empty")
	}
}

func TestGetCurrentPrincipalReturnsDevelopmentStubIdentity(t *testing.T) {
	router, cleanup := newTestRouter(t)
	defer cleanup()

	responseRecorder := performJSONRequest(t, router, nethttp.MethodGet, "/api/auth/me", nil)
	if responseRecorder.Code != nethttp.StatusOK {
		t.Fatalf("status code = %d, want %d", responseRecorder.Code, nethttp.StatusOK)
	}

	var response map[string]any
	decodeResponseBody(t, responseRecorder, &response)

	if response["role"] != "manager" {
		t.Fatalf("role = %v, want %q", response["role"], "manager")
	}
}

func TestCreatePlayerReturnsValidationError(t *testing.T) {
	router, cleanup := newTestRouter(t)
	defer cleanup()

	responseRecorder := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/players",
		map[string]any{
			"name":     "",
			"handicap": 100,
			"status":   "invalid",
		},
	)

	if responseRecorder.Code != nethttp.StatusUnprocessableEntity {
		t.Fatalf("status code = %d, want %d", responseRecorder.Code, nethttp.StatusUnprocessableEntity)
	}
}

func TestPlayerSessionAndRegistrationFlow(t *testing.T) {
	router, cleanup := newTestRouter(t)
	defer cleanup()

	playerResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/players",
		map[string]any{
			"name":     "王大明",
			"handicap": 12.5,
			"status":   "active",
		},
	)
	if playerResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create player status = %d, want %d", playerResponse.Code, nethttp.StatusCreated)
	}

	var player map[string]any
	decodeResponseBody(t, playerResponse, &player)

	sessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions",
		map[string]any{
			"date":                 "2026-04-05T08:00:00Z",
			"courseName":           "台北高爾夫球場",
			"maxPlayers":           4,
			"registrationDeadline": "2026-03-29T23:59:00Z",
			"status":               "open",
		},
	)
	if sessionResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create session status = %d, want %d", sessionResponse.Code, nethttp.StatusCreated)
	}

	var session map[string]any
	decodeResponseBody(t, sessionResponse, &session)

	registrationResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions/"+session["id"].(string)+"/registrations",
		map[string]any{
			"playerId": player["id"],
			"status":   "confirmed",
		},
	)
	if registrationResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create registration status = %d, want %d", registrationResponse.Code, nethttp.StatusCreated)
	}

	listPlayersResponse := performJSONRequest(t, router, nethttp.MethodGet, "/api/players", nil)
	if listPlayersResponse.Code != nethttp.StatusOK {
		t.Fatalf("list players status = %d, want %d", listPlayersResponse.Code, nethttp.StatusOK)
	}

	var players []map[string]any
	decodeResponseBody(t, listPlayersResponse, &players)
	if len(players) != 1 {
		t.Fatalf("players length = %d, want %d", len(players), 1)
	}

	listRegistrationsResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodGet,
		"/api/sessions/"+session["id"].(string)+"/registrations",
		nil,
	)
	if listRegistrationsResponse.Code != nethttp.StatusOK {
		t.Fatalf("list registrations status = %d, want %d", listRegistrationsResponse.Code, nethttp.StatusOK)
	}

	var registrations []map[string]any
	decodeResponseBody(t, listRegistrationsResponse, &registrations)
	if len(registrations) != 1 {
		t.Fatalf("registrations length = %d, want %d", len(registrations), 1)
	}
}

func TestPlayerDetailUpdateAndFilteringFlow(t *testing.T) {
	router, cleanup := newTestRouter(t)
	defer cleanup()

	firstPlayerResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/players",
		map[string]any{
			"name":     "王大明",
			"handicap": 12.5,
			"status":   "active",
		},
	)
	if firstPlayerResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create first player status = %d, want %d", firstPlayerResponse.Code, nethttp.StatusCreated)
	}

	var firstPlayer map[string]any
	decodeResponseBody(t, firstPlayerResponse, &firstPlayer)

	secondPlayerResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/players",
		map[string]any{
			"name":     "李小華",
			"handicap": 18,
			"status":   "active",
		},
	)
	if secondPlayerResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create second player status = %d, want %d", secondPlayerResponse.Code, nethttp.StatusCreated)
	}

	playerDetailResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodGet,
		"/api/players/"+firstPlayer["id"].(string),
		nil,
	)
	if playerDetailResponse.Code != nethttp.StatusOK {
		t.Fatalf("get player status = %d, want %d", playerDetailResponse.Code, nethttp.StatusOK)
	}

	updateResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/players/"+firstPlayer["id"].(string),
		map[string]any{
			"name":     "王大明",
			"handicap": 11.5,
			"email":    "wang@example.com",
			"status":   "inactive",
		},
	)
	if updateResponse.Code != nethttp.StatusOK {
		t.Fatalf("update player status = %d, want %d", updateResponse.Code, nethttp.StatusOK)
	}

	filteredResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodGet,
		"/api/players?status=inactive&query=%E7%8E%8B",
		nil,
	)
	if filteredResponse.Code != nethttp.StatusOK {
		t.Fatalf("filtered list status = %d, want %d", filteredResponse.Code, nethttp.StatusOK)
	}

	var filteredPlayers []map[string]any
	decodeResponseBody(t, filteredResponse, &filteredPlayers)
	if len(filteredPlayers) != 1 {
		t.Fatalf("filtered players length = %d, want %d", len(filteredPlayers), 1)
	}

	if filteredPlayers[0]["status"] != "inactive" {
		t.Fatalf("filtered player status = %v, want %q", filteredPlayers[0]["status"], "inactive")
	}
}

func newTestRouter(t *testing.T) (*gin.Engine, func()) {
	t.Helper()

	database := openTestDatabase(t)
	cleanup := func() {
		database.Close()
	}

	testConfig, err := config.LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv() error = %v", err)
	}

	return NewRouter(database, testConfig), cleanup
}

func openTestDatabase(t *testing.T) *sql.DB {
	t.Helper()

	database, err := appdb.Open(config.DBConfig{
		Path:        filepath.Join(t.TempDir(), "http-router.sqlite"),
		AutoMigrate: true,
	})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	if err := appdb.RunMigrations(context.Background(), database); err != nil {
		t.Fatalf("RunMigrations() error = %v", err)
	}

	return database
}

func performJSONRequest(
	t *testing.T,
	router httpHandler,
	method string,
	target string,
	payload any,
) *httptest.ResponseRecorder {
	t.Helper()

	var body *bytes.Buffer
	if payload == nil {
		body = bytes.NewBuffer(nil)
	} else {
		marshaledPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}

		body = bytes.NewBuffer(marshaledPayload)
	}

	request := httptest.NewRequest(method, target, body)
	request.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	return responseRecorder
}

type httpHandler interface {
	ServeHTTP(writer nethttp.ResponseWriter, request *nethttp.Request)
}

func decodeResponseBody(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()

	if err := json.Unmarshal(response.Body.Bytes(), target); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
}
