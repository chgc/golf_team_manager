package apihttp

import (
	nethttp "net/http"

	"github.com/chgc/golf_team_manager/backend/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetCurrentPrincipal(c *gin.Context) {
	principal, ok := middleware.PrincipalFromContext(c)
	if !ok {
		c.JSON(nethttp.StatusInternalServerError, ErrorResponse{
			Error: APIError{
				Code:    "internal_error",
				Message: "auth principal is unavailable",
			},
		})
		return
	}

	c.JSON(nethttp.StatusOK, principal)
}
