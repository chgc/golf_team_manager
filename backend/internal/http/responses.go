package apihttp

import (
	"errors"
	nethttp "net/http"

	"github.com/chgc/golf_team_manager/backend/internal/domain"
	"github.com/chgc/golf_team_manager/backend/internal/repository"
	"github.com/chgc/golf_team_manager/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func respondError(c *gin.Context, err error) {
	var validationErrors domain.ValidationErrors
	if errors.As(err, &validationErrors) {
		details := make([]string, 0, len(validationErrors))
		for _, validationError := range validationErrors {
			details = append(details, validationError.Error())
		}

		c.JSON(nethttp.StatusUnprocessableEntity, ErrorResponse{
			Error: APIError{
				Code:    "validation_failed",
				Message: "request validation failed",
				Details: details,
			},
		})
		return
	}

	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(nethttp.StatusNotFound, ErrorResponse{
			Error: APIError{Code: "not_found", Message: "requested resource was not found"},
		})
	case errors.Is(err, repository.ErrConflict):
		c.JSON(nethttp.StatusConflict, ErrorResponse{
			Error: APIError{Code: "conflict", Message: "resource conflict detected"},
		})
	case errors.Is(err, service.ErrPlayerInactive):
		c.JSON(nethttp.StatusConflict, ErrorResponse{
			Error: APIError{Code: "player_inactive", Message: "inactive players cannot register"},
		})
	case errors.Is(err, service.ErrSessionNotOpen):
		c.JSON(nethttp.StatusConflict, ErrorResponse{
			Error: APIError{Code: "session_not_open", Message: "session is not open for registration"},
		})
	case errors.Is(err, service.ErrSessionCapacityFull):
		c.JSON(nethttp.StatusConflict, ErrorResponse{
			Error: APIError{Code: "session_capacity_full", Message: "session capacity has been reached"},
		})
	default:
		c.JSON(nethttp.StatusInternalServerError, ErrorResponse{
			Error: APIError{Code: "internal_error", Message: "an unexpected error occurred"},
		})
	}
}
