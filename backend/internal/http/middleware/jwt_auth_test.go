package middleware

import (
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func TestJWTAuthInjectsPrincipal(t *testing.T) {
	gin.SetMode(gin.TestMode)

	now := time.Now().UTC()
	tokenManager := auth.NewHMACTokenManager("test-secret")
	token, err := tokenManager.Sign(auth.Claims{
		Subject:         "user-1",
		PlayerID:        "player-1",
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

	router := gin.New()
	router.Use(JWTAuth(tokenManager))
	router.GET("/protected", func(c *gin.Context) {
		principal, ok := PrincipalFromContext(c)
		if !ok {
			t.Fatal("PrincipalFromContext() ok = false, want true")
		}

		c.JSON(nethttp.StatusOK, principal)
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/protected", nil)
	request.Header.Set(authorizationHeader, "Bearer "+token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, nethttp.StatusOK)
	}
}

func TestJWTAuthRejectsMissingBearerToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(JWTAuth(auth.NewHMACTokenManager("test-secret")))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(nethttp.StatusOK)
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/protected", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status code = %d, want %d", response.Code, nethttp.StatusUnauthorized)
	}
}

func TestJWTAuthRejectsInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(JWTAuth(auth.NewHMACTokenManager("test-secret")))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(nethttp.StatusOK)
	})

	request := httptest.NewRequest(nethttp.MethodGet, "/protected", nil)
	request.Header.Set(authorizationHeader, "Bearer invalid-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusUnauthorized {
		t.Fatalf("status code = %d, want %d", response.Code, nethttp.StatusUnauthorized)
	}
}
