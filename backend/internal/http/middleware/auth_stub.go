package middleware

import (
	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/chgc/golf_team_manager/backend/internal/config"
	"github.com/gin-gonic/gin"
)

const (
	authContextKey         = "auth.principal"
	debugDisplayNameHeader = "X-Debug-Display-Name"
	debugPlayerIDHeader    = "X-Debug-Player-ID"
	debugRoleHeader        = "X-Debug-Role"
	debugSubjectHeader     = "X-Debug-Subject"
	debugUserIDHeader      = "X-Debug-User-ID"
)

func DevelopmentAuth(cfg config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := valueOrFallback(c.GetHeader(debugSubjectHeader), cfg.DevSubject)
		principal := auth.Principal{
			DisplayName: valueOrFallback(c.GetHeader(debugDisplayNameHeader), cfg.DevDisplayName),
			PlayerID:    valueOrFallback(c.GetHeader(debugPlayerIDHeader), cfg.DevPlayerID),
			Provider:    auth.ProviderDevelopmentStub,
			Role:        auth.Role(valueOrFallback(c.GetHeader(debugRoleHeader), cfg.DevRole)),
			Subject:     subject,
			UserID:      valueOrFallback(c.GetHeader(debugUserIDHeader), cfg.DevUserID),
		}

		SetPrincipal(c, principal)
		c.Next()
	}
}

func SetPrincipal(c *gin.Context, principal auth.Principal) {
	c.Set(authContextKey, principal)
}

func PrincipalFromContext(c *gin.Context) (auth.Principal, bool) {
	value, ok := c.Get(authContextKey)
	if !ok {
		return auth.Principal{}, false
	}

	principal, ok := value.(auth.Principal)
	return principal, ok
}

func valueOrFallback(value string, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
