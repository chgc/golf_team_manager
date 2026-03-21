package middleware

import (
	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

const authContextKey = "auth.principal"

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
