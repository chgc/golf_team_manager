package apihttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"strings"
	"time"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/chgc/golf_team_manager/backend/internal/domain"
	"github.com/chgc/golf_team_manager/backend/internal/http/middleware"
	"github.com/chgc/golf_team_manager/backend/internal/repository"
	"github.com/chgc/golf_team_manager/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type adminUserResponse struct {
	CreatedAt   time.Time     `json:"createdAt"`
	DisplayName string        `json:"displayName"`
	PlayerID    string        `json:"playerId,omitempty"`
	Provider    auth.Provider `json:"provider"`
	Role        auth.Role     `json:"role"`
	Subject     string        `json:"subject"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	UserID      string        `json:"userId"`
}

type adminUserPatchRequest struct {
	PlayerID optionalNullableString `json:"playerId"`
	Role     optionalRole           `json:"role"`
}

type optionalNullableString struct {
	Set   bool
	Value *string
}

type optionalRole struct {
	Set   bool
	Value *auth.Role
}

func (h *Handlers) ListAdminUsers(c *gin.Context) {
	if !requireManagerAccess(c) {
		return
	}

	filter, err := buildAdminUserListFilter(c.Query("linkState"), c.Query("role"))
	if err != nil {
		respondError(c, domain.ValidationErrors{err})
		return
	}

	users, err := h.userAdminService.List(c.Request.Context(), filter)
	if err != nil {
		respondError(c, err)
		return
	}

	response := make([]adminUserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, mapAdminUserResponse(user))
	}

	c.JSON(nethttp.StatusOK, response)
}

func (h *Handlers) UpdateAdminUser(c *gin.Context) {
	if !requireManagerAccess(c) {
		return
	}

	var request adminUserPatchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, domain.ValidationErrors{err})
		return
	}

	updateInput, err := request.toServiceInput()
	if err != nil {
		respondError(c, domain.ValidationErrors{err})
		return
	}

	user, err := h.userAdminService.Update(c.Request.Context(), c.Param("userId"), updateInput)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(nethttp.StatusOK, mapAdminUserResponse(user))
}

func buildAdminUserListFilter(linkState string, role string) (repository.UserListFilter, error) {
	filter := repository.UserListFilter{}

	switch strings.TrimSpace(linkState) {
	case "", "all":
	case string(repository.UserLinkStateLinked):
		filter.LinkState = repository.UserLinkStateLinked
	case string(repository.UserLinkStateUnlinked):
		filter.LinkState = repository.UserLinkStateUnlinked
	default:
		return repository.UserListFilter{}, fmt.Errorf("invalid linkState filter %q", linkState)
	}

	switch strings.TrimSpace(role) {
	case "", "all":
	case string(auth.RoleManager):
		filter.Role = auth.RoleManager
	case string(auth.RolePlayer):
		filter.Role = auth.RolePlayer
	default:
		return repository.UserListFilter{}, fmt.Errorf("invalid role filter %q", role)
	}

	return filter, nil
}

func mapAdminUserResponse(user auth.User) adminUserResponse {
	return adminUserResponse{
		CreatedAt:   user.CreatedAt,
		DisplayName: user.DisplayName,
		PlayerID:    user.PlayerID,
		Provider:    user.Provider,
		Role:        user.Role,
		Subject:     user.Subject,
		UpdatedAt:   user.UpdatedAt,
		UserID:      user.ID,
	}
}

func requireManagerAccess(c *gin.Context) bool {
	principal, ok := middleware.PrincipalFromContext(c)
	if !ok {
		c.JSON(nethttp.StatusInternalServerError, ErrorResponse{
			Error: APIError{
				Code:    "internal_error",
				Message: "auth principal is unavailable",
			},
		})
		return false
	}

	if principal.Role != auth.RoleManager {
		c.JSON(nethttp.StatusForbidden, ErrorResponse{
			Error: APIError{
				Code:    "forbidden",
				Message: "manager access is required",
			},
		})
		return false
	}

	return true
}

func (f *optionalNullableString) UnmarshalJSON(data []byte) error {
	f.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		f.Value = nil
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	trimmedValue := strings.TrimSpace(value)
	f.Value = &trimmedValue
	return nil
}

func (f *optionalRole) UnmarshalJSON(data []byte) error {
	f.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		f.Value = nil
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	normalizedRole := auth.Role(strings.TrimSpace(value))
	f.Value = &normalizedRole
	return nil
}

func (r adminUserPatchRequest) toServiceInput() (service.UserAdminUpdateInput, error) {
	if !r.Role.Set && !r.PlayerID.Set {
		return service.UserAdminUpdateInput{}, fmt.Errorf("at least one of role or playerId must be provided")
	}

	input := service.UserAdminUpdateInput{}
	if r.Role.Set {
		if r.Role.Value == nil {
			return service.UserAdminUpdateInput{}, fmt.Errorf("role must not be null")
		}

		switch *r.Role.Value {
		case auth.RoleManager, auth.RolePlayer:
			input.Role = r.Role.Value
		default:
			return service.UserAdminUpdateInput{}, fmt.Errorf("role must be manager or player")
		}
	}

	if r.PlayerID.Set {
		if r.PlayerID.Value == nil {
			input.ClearPlayerLink = true
		} else {
			if *r.PlayerID.Value == "" {
				return service.UserAdminUpdateInput{}, fmt.Errorf("playerId must not be empty")
			}
			input.PlayerID = r.PlayerID.Value
		}
	}

	return input, nil
}
