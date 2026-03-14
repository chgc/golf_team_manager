package apihttp

import (
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/chgc/golf_team_manager/backend/internal/http/handlers"
)

func TestNewRouterHealthEndpoint(t *testing.T) {
	router := NewRouter()

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
}
