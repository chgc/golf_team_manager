package middleware

import (
	nethttp "net/http"
	"strings"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

const authorizationHeader = "Authorization"

func JWTAuth(tokens auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerValue := strings.TrimSpace(c.GetHeader(authorizationHeader))
		if !strings.HasPrefix(headerValue, "Bearer ") {
			respondUnauthorized(c)
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(headerValue, "Bearer "))
		if token == "" {
			respondUnauthorized(c)
			return
		}

		claims, err := tokens.Validate(token)
		if err != nil {
			respondUnauthorized(c)
			return
		}

		SetPrincipal(c, claims.Principal())
		c.Next()
	}
}

func respondUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(nethttp.StatusUnauthorized, gin.H{
		"error": gin.H{
			"code":    "unauthorized",
			"message": "valid bearer token is required",
		},
	})
}
