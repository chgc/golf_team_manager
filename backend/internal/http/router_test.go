package apihttp

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
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

	if response["userId"] != "dev-user:dev-manager" {
		t.Fatalf("userId = %v, want %q", response["userId"], "dev-user:dev-manager")
	}
}

func TestGetCurrentPrincipalReturnsLineModeUnauthorizedWithoutToken(t *testing.T) {
	router, cleanup := newLineTestRouter(t, RouterDependencies{})
	defer cleanup()

	responseRecorder := performJSONRequest(t, router, nethttp.MethodGet, "/api/auth/me", nil)
	if responseRecorder.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status code = %d, want %d", responseRecorder.Code, nethttp.StatusUnauthorized)
	}
}

func TestGetCurrentPrincipalReturnsLineModePrincipalWithoutPlayerLink(t *testing.T) {
	cfg := newLineTestConfig()
	now := time.Now().UTC()
	tokenManager := auth.NewHMACTokenManager(cfg.Auth.JWTSecret)

	router, cleanup := newLineTestRouter(t, RouterDependencies{
		TokenManager: tokenManager,
	})
	defer cleanup()

	token, err := tokenManager.Sign(auth.Claims{
		Subject:         "user-1",
		Provider:        auth.ProviderLINEOAuth,
		ProviderSubject: "line-user-1",
		Role:            auth.RolePlayer,
		DisplayName:     "王小明",
		IssuedAt:        now.Unix(),
		ExpiresAt:       now.Add(time.Hour).Unix(),
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	responseRecorder := performJSONRequestWithHeaders(
		t,
		router,
		nethttp.MethodGet,
		"/api/auth/me",
		nil,
		map[string]string{
			"Authorization": "Bearer " + token,
		},
	)
	if responseRecorder.Code != nethttp.StatusOK {
		t.Fatalf("status code = %d, want %d", responseRecorder.Code, nethttp.StatusOK)
	}

	var response map[string]any
	decodeResponseBody(t, responseRecorder, &response)
	if _, ok := response["playerId"]; ok {
		t.Fatalf("playerId present = true, want false")
	}
}

func TestLineCallbackCreatesUserAndRedirectsWithJWT(t *testing.T) {
	cfg := newLineTestConfig()
	tokenManager := auth.NewHMACTokenManager(cfg.Auth.JWTSecret)
	lineProvider := &stubLineProvider{
		authorizeURL: "https://line.example/authorize",
		tokenResponse: auth.LineTokenResponse{
			IDToken: "line-id-token",
		},
		identity: auth.LineIdentity{
			Subject:     "line-user-1",
			DisplayName: "王小明",
		},
	}

	router, cleanup, database := newLineTestRouterWithDatabase(t, RouterDependencies{
		LineProvider: lineProvider,
		TokenManager: tokenManager,
	})
	defer cleanup()

	loginRequest := httptest.NewRequest(nethttp.MethodGet, "/api/auth/line/login", nil)
	loginResponse := httptest.NewRecorder()
	router.ServeHTTP(loginResponse, loginRequest)

	if loginResponse.Code != nethttp.StatusFound {
		t.Fatalf("login status code = %d, want %d", loginResponse.Code, nethttp.StatusFound)
	}

	cookies := loginResponse.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("login cookies length = 0, want at least 1")
	}

	stateCookie := cookies[0]
	flow, err := auth.DecodeOAuthState(stateCookie.Value)
	if err != nil {
		t.Fatalf("DecodeOAuthState() error = %v", err)
	}

	callbackRequest := httptest.NewRequest(
		nethttp.MethodGet,
		"/api/auth/line/callback?code=test-code&state="+flow.State,
		nil,
	)
	callbackRequest.AddCookie(stateCookie)
	callbackResponse := httptest.NewRecorder()
	router.ServeHTTP(callbackResponse, callbackRequest)

	if callbackResponse.Code != nethttp.StatusFound {
		t.Fatalf("callback status code = %d, want %d", callbackResponse.Code, nethttp.StatusFound)
	}

	location := callbackResponse.Header().Get("Location")
	if !strings.HasPrefix(location, cfg.Auth.FrontendURL+"/auth/done#token=") {
		t.Fatalf("redirect location = %q, want prefix %q", location, cfg.Auth.FrontendURL+"/auth/done#token=")
	}

	if lineProvider.exchangedCode != "test-code" {
		t.Fatalf("exchanged code = %q, want %q", lineProvider.exchangedCode, "test-code")
	}

	if lineProvider.verifiedNonce != flow.Nonce {
		t.Fatalf("verified nonce = %q, want %q", lineProvider.verifiedNonce, flow.Nonce)
	}

	var (
		userID      string
		playerID    sql.NullString
		displayName string
		role        string
	)
	err = database.QueryRowContext(
		context.Background(),
		`SELECT id, player_id, display_name, role FROM users WHERE auth_provider = 'line' AND provider_subject = 'line-user-1'`,
	).Scan(&userID, &playerID, &displayName, &role)
	if err != nil {
		t.Fatalf("select user error = %v", err)
	}

	if playerID.Valid {
		t.Fatalf("playerID.Valid = %t, want false", playerID.Valid)
	}

	if displayName != "王小明" {
		t.Fatalf("displayName = %q, want %q", displayName, "王小明")
	}

	if role != "player" {
		t.Fatalf("role = %q, want %q", role, "player")
	}

	if userID == "" {
		t.Fatal("userID is empty")
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

func TestSessionDetailUpdateStatusAndAutoCloseFlow(t *testing.T) {
	router, cleanup := newTestRouter(t)
	defer cleanup()

	expiringSessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions",
		map[string]any{
			"date":                 "2026-06-01T08:00:00Z",
			"courseName":           "林口高爾夫球場",
			"maxPlayers":           8,
			"registrationDeadline": "2025-01-01T00:00:00Z",
			"status":               "open",
		},
	)
	if expiringSessionResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create expiring session status = %d, want %d", expiringSessionResponse.Code, nethttp.StatusCreated)
	}

	var expiringSession map[string]any
	decodeResponseBody(t, expiringSessionResponse, &expiringSession)

	detailResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodGet,
		"/api/sessions/"+expiringSession["id"].(string),
		nil,
	)
	if detailResponse.Code != nethttp.StatusOK {
		t.Fatalf("get session status = %d, want %d", detailResponse.Code, nethttp.StatusOK)
	}

	var detailSession map[string]any
	decodeResponseBody(t, detailResponse, &detailSession)
	if detailSession["status"] != "closed" {
		t.Fatalf("auto-closed session status = %v, want %q", detailSession["status"], "closed")
	}

	activeSessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions",
		map[string]any{
			"date":                 "2026-06-15T08:00:00Z",
			"courseName":           "台中高爾夫俱樂部",
			"maxPlayers":           12,
			"registrationDeadline": "2026-06-10T23:59:00Z",
			"status":               "open",
			"notes":                "Morning round",
		},
	)
	if activeSessionResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create active session status = %d, want %d", activeSessionResponse.Code, nethttp.StatusCreated)
	}

	var activeSession map[string]any
	decodeResponseBody(t, activeSessionResponse, &activeSession)

	updateResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/sessions/"+activeSession["id"].(string),
		map[string]any{
			"date":                 "2026-06-15T09:00:00Z",
			"courseName":           "台中高爾夫俱樂部",
			"courseAddress":        "台中市大雅區",
			"maxPlayers":           16,
			"registrationDeadline": "2026-06-10T23:59:00Z",
			"status":               "open",
			"notes":                "Updated tee time",
		},
	)
	if updateResponse.Code != nethttp.StatusOK {
		t.Fatalf("update session status = %d, want %d", updateResponse.Code, nethttp.StatusOK)
	}

	invalidTransitionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/sessions/"+activeSession["id"].(string),
		map[string]any{
			"date":                 "2026-06-15T09:00:00Z",
			"courseName":           "台中高爾夫俱樂部",
			"courseAddress":        "台中市大雅區",
			"maxPlayers":           16,
			"registrationDeadline": "2026-06-10T23:59:00Z",
			"status":               "confirmed",
			"notes":                "Updated tee time",
		},
	)
	if invalidTransitionResponse.Code != nethttp.StatusUnprocessableEntity {
		t.Fatalf("invalid transition status = %d, want %d", invalidTransitionResponse.Code, nethttp.StatusUnprocessableEntity)
	}

	closeResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/sessions/"+activeSession["id"].(string),
		map[string]any{
			"date":                 "2026-06-15T09:00:00Z",
			"courseName":           "台中高爾夫俱樂部",
			"courseAddress":        "台中市大雅區",
			"maxPlayers":           16,
			"registrationDeadline": "2026-06-10T23:59:00Z",
			"status":               "closed",
			"notes":                "Registration closed",
		},
	)
	if closeResponse.Code != nethttp.StatusOK {
		t.Fatalf("close session status = %d, want %d", closeResponse.Code, nethttp.StatusOK)
	}

	confirmResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/sessions/"+activeSession["id"].(string),
		map[string]any{
			"date":                 "2026-06-15T09:00:00Z",
			"courseName":           "台中高爾夫俱樂部",
			"courseAddress":        "台中市大雅區",
			"maxPlayers":           16,
			"registrationDeadline": "2026-06-10T23:59:00Z",
			"status":               "confirmed",
			"notes":                "Pairings locked",
		},
	)
	if confirmResponse.Code != nethttp.StatusOK {
		t.Fatalf("confirm session status = %d, want %d", confirmResponse.Code, nethttp.StatusOK)
	}

	completeResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/sessions/"+activeSession["id"].(string),
		map[string]any{
			"date":                 "2026-06-15T09:00:00Z",
			"courseName":           "台中高爾夫俱樂部",
			"courseAddress":        "台中市大雅區",
			"maxPlayers":           16,
			"registrationDeadline": "2026-06-10T23:59:00Z",
			"status":               "completed",
			"notes":                "Round finished",
		},
	)
	if completeResponse.Code != nethttp.StatusOK {
		t.Fatalf("complete session status = %d, want %d", completeResponse.Code, nethttp.StatusOK)
	}

	listResponse := performJSONRequest(t, router, nethttp.MethodGet, "/api/sessions", nil)
	if listResponse.Code != nethttp.StatusOK {
		t.Fatalf("list sessions status = %d, want %d", listResponse.Code, nethttp.StatusOK)
	}

	var sessions []map[string]any
	decodeResponseBody(t, listResponse, &sessions)
	if len(sessions) != 2 {
		t.Fatalf("sessions length = %d, want %d", len(sessions), 2)
	}
}

func TestRegistrationUpdateAndValidationFlow(t *testing.T) {
	router, cleanup := newTestRouter(t)
	defer cleanup()

	sessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions",
		map[string]any{
			"date":                 "2026-08-20T08:00:00Z",
			"courseName":           "新竹高爾夫俱樂部",
			"maxPlayers":           1,
			"registrationDeadline": "2026-08-15T23:59:00Z",
			"status":               "open",
		},
	)
	if sessionResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create session status = %d, want %d", sessionResponse.Code, nethttp.StatusCreated)
	}

	var session map[string]any
	decodeResponseBody(t, sessionResponse, &session)

	activePlayerResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/players",
		map[string]any{
			"name":     "張小美",
			"handicap": 10,
			"status":   "active",
		},
	)
	if activePlayerResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create active player status = %d, want %d", activePlayerResponse.Code, nethttp.StatusCreated)
	}

	var activePlayer map[string]any
	decodeResponseBody(t, activePlayerResponse, &activePlayer)

	registrationResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions/"+session["id"].(string)+"/registrations",
		map[string]any{
			"playerId": activePlayer["id"],
			"status":   "confirmed",
		},
	)
	if registrationResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create registration status = %d, want %d", registrationResponse.Code, nethttp.StatusCreated)
	}

	var registration map[string]any
	decodeResponseBody(t, registrationResponse, &registration)

	cancelResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/registrations/"+registration["id"].(string),
		map[string]any{
			"status": "cancelled",
		},
	)
	if cancelResponse.Code != nethttp.StatusOK {
		t.Fatalf("cancel registration status = %d, want %d", cancelResponse.Code, nethttp.StatusOK)
	}

	restoreResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/registrations/"+registration["id"].(string),
		map[string]any{
			"status": "confirmed",
		},
	)
	if restoreResponse.Code != nethttp.StatusOK {
		t.Fatalf("restore registration status = %d, want %d", restoreResponse.Code, nethttp.StatusOK)
	}

	duplicateResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions/"+session["id"].(string)+"/registrations",
		map[string]any{
			"playerId": activePlayer["id"],
			"status":   "confirmed",
		},
	)
	if duplicateResponse.Code != nethttp.StatusConflict {
		t.Fatalf("duplicate registration status = %d, want %d", duplicateResponse.Code, nethttp.StatusConflict)
	}

	inactivePlayerResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/players",
		map[string]any{
			"name":     "陳小安",
			"handicap": 20,
			"status":   "inactive",
		},
	)
	if inactivePlayerResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create inactive player status = %d, want %d", inactivePlayerResponse.Code, nethttp.StatusCreated)
	}

	var inactivePlayer map[string]any
	decodeResponseBody(t, inactivePlayerResponse, &inactivePlayer)

	inactiveRegistrationResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions/"+session["id"].(string)+"/registrations",
		map[string]any{
			"playerId": inactivePlayer["id"],
			"status":   "confirmed",
		},
	)
	if inactiveRegistrationResponse.Code != nethttp.StatusConflict {
		t.Fatalf("inactive player registration status = %d, want %d", inactiveRegistrationResponse.Code, nethttp.StatusConflict)
	}

	cancelAgainResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/registrations/"+registration["id"].(string),
		map[string]any{
			"status": "cancelled",
		},
	)
	if cancelAgainResponse.Code != nethttp.StatusOK {
		t.Fatalf("second cancel registration status = %d, want %d", cancelAgainResponse.Code, nethttp.StatusOK)
	}

	closeSessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/sessions/"+session["id"].(string),
		map[string]any{
			"date":                 "2026-08-20T08:00:00Z",
			"courseName":           "新竹高爾夫俱樂部",
			"maxPlayers":           1,
			"registrationDeadline": "2026-08-15T23:59:00Z",
			"status":               "closed",
		},
	)
	if closeSessionResponse.Code != nethttp.StatusOK {
		t.Fatalf("close session status = %d, want %d", closeSessionResponse.Code, nethttp.StatusOK)
	}

	restoreClosedResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/registrations/"+registration["id"].(string),
		map[string]any{
			"status": "confirmed",
		},
	)
	if restoreClosedResponse.Code != nethttp.StatusConflict {
		t.Fatalf("restore closed session registration status = %d, want %d", restoreClosedResponse.Code, nethttp.StatusConflict)
	}

	fullSessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions",
		map[string]any{
			"date":                 "2026-09-01T08:00:00Z",
			"courseName":           "桃園高爾夫球場",
			"maxPlayers":           1,
			"registrationDeadline": "2026-08-25T23:59:00Z",
			"status":               "open",
		},
	)
	if fullSessionResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create full session status = %d, want %d", fullSessionResponse.Code, nethttp.StatusCreated)
	}

	var fullSession map[string]any
	decodeResponseBody(t, fullSessionResponse, &fullSession)

	firstFullResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions/"+fullSession["id"].(string)+"/registrations",
		map[string]any{
			"playerId": activePlayer["id"],
			"status":   "confirmed",
		},
	)
	if firstFullResponse.Code != nethttp.StatusCreated {
		t.Fatalf("first full-session registration status = %d, want %d", firstFullResponse.Code, nethttp.StatusCreated)
	}

	anotherPlayerResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/players",
		map[string]any{
			"name":     "林小綠",
			"handicap": 14,
			"status":   "active",
		},
	)
	if anotherPlayerResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create another player status = %d, want %d", anotherPlayerResponse.Code, nethttp.StatusCreated)
	}

	var anotherPlayer map[string]any
	decodeResponseBody(t, anotherPlayerResponse, &anotherPlayer)

	capacityFullResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions/"+fullSession["id"].(string)+"/registrations",
		map[string]any{
			"playerId": anotherPlayer["id"],
			"status":   "confirmed",
		},
	)
	if capacityFullResponse.Code != nethttp.StatusConflict {
		t.Fatalf("capacity full registration status = %d, want %d", capacityFullResponse.Code, nethttp.StatusConflict)
	}
}

func TestGetReservationSummaryFlow(t *testing.T) {
	router, cleanup := newTestRouter(t)
	defer cleanup()

	firstPlayerResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/players",
		map[string]any{
			"name":     "Alice",
			"handicap": 10,
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
			"name":     "Bob",
			"handicap": 14,
			"status":   "active",
		},
	)
	if secondPlayerResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create second player status = %d, want %d", secondPlayerResponse.Code, nethttp.StatusCreated)
	}

	var secondPlayer map[string]any
	decodeResponseBody(t, secondPlayerResponse, &secondPlayer)

	sessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions",
		map[string]any{
			"date":                 "2026-10-01T08:00:00Z",
			"courseName":           "Sunrise Golf Club",
			"courseAddress":        "",
			"maxPlayers":           8,
			"registrationDeadline": "2026-09-25T23:59:00Z",
			"status":               "open",
		},
	)
	if sessionResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create session status = %d, want %d", sessionResponse.Code, nethttp.StatusCreated)
	}

	var session map[string]any
	decodeResponseBody(t, sessionResponse, &session)

	for _, player := range []map[string]any{firstPlayer, secondPlayer} {
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
	}

	closeResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/sessions/"+session["id"].(string),
		map[string]any{
			"date":                 "2026-10-01T08:00:00Z",
			"courseName":           "Sunrise Golf Club",
			"courseAddress":        "",
			"maxPlayers":           8,
			"registrationDeadline": "2026-09-25T23:59:00Z",
			"status":               "closed",
		},
	)
	if closeResponse.Code != nethttp.StatusOK {
		t.Fatalf("close session status = %d, want %d", closeResponse.Code, nethttp.StatusOK)
	}

	confirmResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/sessions/"+session["id"].(string),
		map[string]any{
			"date":                 "2026-10-01T08:00:00Z",
			"courseName":           "Sunrise Golf Club",
			"courseAddress":        "",
			"maxPlayers":           8,
			"registrationDeadline": "2026-09-25T23:59:00Z",
			"status":               "confirmed",
		},
	)
	if confirmResponse.Code != nethttp.StatusOK {
		t.Fatalf("confirm session status = %d, want %d", confirmResponse.Code, nethttp.StatusOK)
	}

	summaryResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodGet,
		"/api/reports/sessions/"+session["id"].(string)+"/reservation-summary",
		nil,
	)
	if summaryResponse.Code != nethttp.StatusOK {
		t.Fatalf("get reservation summary status = %d, want %d", summaryResponse.Code, nethttp.StatusOK)
	}

	var summary map[string]any
	decodeResponseBody(t, summaryResponse, &summary)

	if summary["confirmedPlayerCount"] != float64(2) {
		t.Fatalf("confirmedPlayerCount = %v, want %v", summary["confirmedPlayerCount"], float64(2))
	}

	if summary["estimatedGroups"] != float64(1) {
		t.Fatalf("estimatedGroups = %v, want %v", summary["estimatedGroups"], float64(1))
	}

	confirmedPlayers, ok := summary["confirmedPlayers"].([]any)
	if !ok || len(confirmedPlayers) != 2 {
		t.Fatalf("confirmedPlayers length = %d, want %d", len(confirmedPlayers), 2)
	}

	firstConfirmedPlayer := confirmedPlayers[0].(map[string]any)
	secondConfirmedPlayer := confirmedPlayers[1].(map[string]any)
	if firstConfirmedPlayer["playerName"] != "Alice" || secondConfirmedPlayer["playerName"] != "Bob" {
		t.Fatalf("confirmedPlayers order = %v, want Alice then Bob", confirmedPlayers)
	}

	summaryText, ok := summary["summaryText"].(string)
	if !ok {
		t.Fatal("summaryText is not a string")
	}

	for _, expectedLine := range []string{
		"Course: Sunrise Golf Club",
		"Address: N/A",
		"Confirmed Players: 2",
		"Estimated Groups: 1",
		"- Alice",
		"- Bob",
	} {
		if !strings.Contains(summaryText, expectedLine) {
			t.Fatalf("summaryText %q does not contain %q", summaryText, expectedLine)
		}
	}
}

func TestGetReservationSummaryValidationAndAuthorizationFlow(t *testing.T) {
	router, cleanup := newTestRouter(t)
	defer cleanup()

	openSessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions",
		map[string]any{
			"date":                 "2026-11-10T08:00:00Z",
			"courseName":           "North Hills",
			"maxPlayers":           4,
			"registrationDeadline": "2026-11-01T23:59:00Z",
			"status":               "open",
		},
	)
	if openSessionResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create open session status = %d, want %d", openSessionResponse.Code, nethttp.StatusCreated)
	}

	var openSession map[string]any
	decodeResponseBody(t, openSessionResponse, &openSession)

	ineligibleResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodGet,
		"/api/reports/sessions/"+openSession["id"].(string)+"/reservation-summary",
		nil,
	)
	if ineligibleResponse.Code != nethttp.StatusUnprocessableEntity {
		t.Fatalf("open session summary status = %d, want %d", ineligibleResponse.Code, nethttp.StatusUnprocessableEntity)
	}

	var ineligibleError ErrorResponse
	decodeResponseBody(t, ineligibleResponse, &ineligibleError)
	if ineligibleError.Error.Code != "session_not_eligible_for_report" {
		t.Fatalf("ineligible error code = %q, want %q", ineligibleError.Error.Code, "session_not_eligible_for_report")
	}

	emptySessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions",
		map[string]any{
			"date":                 "2026-11-20T08:00:00Z",
			"courseName":           "Ocean View",
			"maxPlayers":           4,
			"registrationDeadline": "2026-11-10T23:59:00Z",
			"status":               "open",
		},
	)
	if emptySessionResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create empty session status = %d, want %d", emptySessionResponse.Code, nethttp.StatusCreated)
	}

	var emptySession map[string]any
	decodeResponseBody(t, emptySessionResponse, &emptySession)

	closeEmptySessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/sessions/"+emptySession["id"].(string),
		map[string]any{
			"date":                 "2026-11-20T08:00:00Z",
			"courseName":           "Ocean View",
			"maxPlayers":           4,
			"registrationDeadline": "2026-11-10T23:59:00Z",
			"status":               "closed",
		},
	)
	if closeEmptySessionResponse.Code != nethttp.StatusOK {
		t.Fatalf("close empty session status = %d, want %d", closeEmptySessionResponse.Code, nethttp.StatusOK)
	}

	confirmEmptySessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/sessions/"+emptySession["id"].(string),
		map[string]any{
			"date":                 "2026-11-20T08:00:00Z",
			"courseName":           "Ocean View",
			"maxPlayers":           4,
			"registrationDeadline": "2026-11-10T23:59:00Z",
			"status":               "confirmed",
		},
	)
	if confirmEmptySessionResponse.Code != nethttp.StatusOK {
		t.Fatalf("confirm empty session status = %d, want %d", confirmEmptySessionResponse.Code, nethttp.StatusOK)
	}

	emptySummaryResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodGet,
		"/api/reports/sessions/"+emptySession["id"].(string)+"/reservation-summary",
		nil,
	)
	if emptySummaryResponse.Code != nethttp.StatusUnprocessableEntity {
		t.Fatalf("empty summary status = %d, want %d", emptySummaryResponse.Code, nethttp.StatusUnprocessableEntity)
	}

	var emptySummaryError ErrorResponse
	decodeResponseBody(t, emptySummaryResponse, &emptySummaryError)
	if emptySummaryError.Error.Code != "reservation_summary_empty" {
		t.Fatalf("empty summary error code = %q, want %q", emptySummaryError.Error.Code, "reservation_summary_empty")
	}

	notFoundResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodGet,
		"/api/reports/sessions/missing-session/reservation-summary",
		nil,
	)
	if notFoundResponse.Code != nethttp.StatusNotFound {
		t.Fatalf("missing session summary status = %d, want %d", notFoundResponse.Code, nethttp.StatusNotFound)
	}

	var notFoundError ErrorResponse
	decodeResponseBody(t, notFoundResponse, &notFoundError)
	if notFoundError.Error.Code != "session_not_found" {
		t.Fatalf("missing session error code = %q, want %q", notFoundError.Error.Code, "session_not_found")
	}

	playerResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/players",
		map[string]any{
			"name":     "Charlie",
			"handicap": 9,
			"status":   "active",
		},
	)
	if playerResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create player status = %d, want %d", playerResponse.Code, nethttp.StatusCreated)
	}

	var player map[string]any
	decodeResponseBody(t, playerResponse, &player)

	eligibleSessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions",
		map[string]any{
			"date":                 "2026-12-01T08:00:00Z",
			"courseName":           "Evergreen",
			"maxPlayers":           4,
			"registrationDeadline": "2026-11-20T23:59:00Z",
			"status":               "open",
		},
	)
	if eligibleSessionResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create eligible session status = %d, want %d", eligibleSessionResponse.Code, nethttp.StatusCreated)
	}

	var eligibleSession map[string]any
	decodeResponseBody(t, eligibleSessionResponse, &eligibleSession)

	eligibleRegistrationResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPost,
		"/api/sessions/"+eligibleSession["id"].(string)+"/registrations",
		map[string]any{
			"playerId": player["id"],
			"status":   "confirmed",
		},
	)
	if eligibleRegistrationResponse.Code != nethttp.StatusCreated {
		t.Fatalf("create eligible registration status = %d, want %d", eligibleRegistrationResponse.Code, nethttp.StatusCreated)
	}

	closeEligibleSessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/sessions/"+eligibleSession["id"].(string),
		map[string]any{
			"date":                 "2026-12-01T08:00:00Z",
			"courseName":           "Evergreen",
			"maxPlayers":           4,
			"registrationDeadline": "2026-11-20T23:59:00Z",
			"status":               "closed",
		},
	)
	if closeEligibleSessionResponse.Code != nethttp.StatusOK {
		t.Fatalf("close eligible session status = %d, want %d", closeEligibleSessionResponse.Code, nethttp.StatusOK)
	}

	confirmEligibleSessionResponse := performJSONRequest(
		t,
		router,
		nethttp.MethodPatch,
		"/api/sessions/"+eligibleSession["id"].(string),
		map[string]any{
			"date":                 "2026-12-01T08:00:00Z",
			"courseName":           "Evergreen",
			"maxPlayers":           4,
			"registrationDeadline": "2026-11-20T23:59:00Z",
			"status":               "confirmed",
		},
	)
	if confirmEligibleSessionResponse.Code != nethttp.StatusOK {
		t.Fatalf("confirm eligible session status = %d, want %d", confirmEligibleSessionResponse.Code, nethttp.StatusOK)
	}

	forbiddenResponse := performJSONRequestWithHeaders(
		t,
		router,
		nethttp.MethodGet,
		"/api/reports/sessions/"+eligibleSession["id"].(string)+"/reservation-summary",
		nil,
		map[string]string{
			"X-Debug-Role": "player",
		},
	)
	if forbiddenResponse.Code != nethttp.StatusForbidden {
		t.Fatalf("player role summary status = %d, want %d", forbiddenResponse.Code, nethttp.StatusForbidden)
	}

	var forbiddenError ErrorResponse
	decodeResponseBody(t, forbiddenResponse, &forbiddenError)
	if forbiddenError.Error.Code != "forbidden" {
		t.Fatalf("forbidden error code = %q, want %q", forbiddenError.Error.Code, "forbidden")
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

func newLineTestRouter(t *testing.T, deps RouterDependencies) (*gin.Engine, func()) {
	router, cleanup, _ := newLineTestRouterWithDatabase(t, deps)
	return router, cleanup
}

func newLineTestRouterWithDatabase(t *testing.T, deps RouterDependencies) (*gin.Engine, func(), *sql.DB) {
	t.Helper()

	database := openTestDatabase(t)
	cleanup := func() {
		database.Close()
	}

	return newRouterWithDependencies(database, newLineTestConfig(), deps), cleanup, database
}

func newLineTestConfig() config.Config {
	return config.Config{
		HTTP: config.HTTPConfig{
			Host:        "localhost",
			Port:        8080,
			ReadTimeout: 5 * time.Second,
		},
		DB: config.DBConfig{
			Path:        "test.sqlite",
			AutoMigrate: true,
		},
		Auth: config.AuthConfig{
			Mode:             "line",
			LineClientID:     "line-client",
			LineClientSecret: "line-secret",
			LineRedirectURI:  "http://localhost:8080/api/auth/line/callback",
			FrontendURL:      "http://localhost:4200",
			JWTSecret:        "jwt-secret",
			JWTTTL:           time.Hour,
		},
	}
}

type stubLineProvider struct {
	authorizeURL    string
	tokenResponse   auth.LineTokenResponse
	identity        auth.LineIdentity
	exchangedCode   string
	verifiedIDToken string
	verifiedNonce   string
}

func (p *stubLineProvider) BuildAuthorizationURL(state string, nonce string) (string, error) {
	return p.authorizeURL + "?state=" + state + "&nonce=" + nonce, nil
}

func (p *stubLineProvider) ExchangeCode(_ context.Context, code string) (auth.LineTokenResponse, error) {
	p.exchangedCode = code
	return p.tokenResponse, nil
}

func (p *stubLineProvider) VerifyIDToken(_ context.Context, idToken string, nonce string) (auth.LineIdentity, error) {
	p.verifiedIDToken = idToken
	p.verifiedNonce = nonce
	identity := p.identity
	identity.Nonce = nonce
	return identity, nil
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
	return performJSONRequestWithHeaders(t, router, method, target, payload, nil)
}

func performJSONRequestWithHeaders(
	t *testing.T,
	router httpHandler,
	method string,
	target string,
	payload any,
	headers map[string]string,
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
	for key, value := range headers {
		request.Header.Set(key, value)
	}
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
