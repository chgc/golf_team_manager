package apihttp

import (
	"errors"
	"fmt"
	nethttp "net/http"
	"net/url"
	"strings"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

type LineAuthHandlers struct {
	service      *auth.LineAuthService
	frontendURL  string
	cookieSecure bool
}

func NewLineAuthHandlers(service *auth.LineAuthService, frontendURL string, cookieSecure bool) *LineAuthHandlers {
	return &LineAuthHandlers{
		service:      service,
		frontendURL:  frontendURL,
		cookieSecure: cookieSecure,
	}
}

func (h *LineAuthHandlers) Login(c *gin.Context) {
	flow, redirectURL, err := h.service.NewLoginFlow()
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, ErrorResponse{
			Error: APIError{Code: "internal_error", Message: "failed to start LINE login"},
		})
		return
	}

	encodedState, err := auth.EncodeOAuthState(flow)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, ErrorResponse{
			Error: APIError{Code: "internal_error", Message: "failed to persist LINE login state"},
		})
		return
	}

	c.SetSameSite(nethttp.SameSiteLaxMode)
	c.SetCookie(auth.OAuthStateCookieName, encodedState, 300, "/", "", h.cookieSecure, true)
	c.Redirect(nethttp.StatusFound, redirectURL)
}

func (h *LineAuthHandlers) Callback(c *gin.Context) {
	encodedState, err := c.Cookie(auth.OAuthStateCookieName)
	if err != nil {
		h.respondLineCallbackError(c, auth.ErrInvalidOAuthState)
		return
	}

	flow, err := auth.DecodeOAuthState(encodedState)
	if err != nil {
		h.respondLineCallbackError(c, auth.ErrInvalidOAuthState)
		return
	}

	token, _, err := h.service.CompleteCallback(c.Request.Context(), c.Query("code"), flow, c.Query("state"))
	if err != nil {
		h.respondLineCallbackError(c, err)
		return
	}

	c.SetSameSite(nethttp.SameSiteLaxMode)
	c.SetCookie(auth.OAuthStateCookieName, "", -1, "/", "", h.cookieSecure, true)
	c.Redirect(nethttp.StatusFound, buildAuthDoneRedirect(h.frontendURL, token))
}

func (h *LineAuthHandlers) respondLineCallbackError(c *gin.Context, err error) {
	statusCode := nethttp.StatusInternalServerError
	message := "failed to complete LINE login"
	code := "internal_error"

	switch {
	case errors.Is(err, auth.ErrInvalidOAuthState), errors.Is(err, auth.ErrMissingOAuthCode):
		statusCode = nethttp.StatusBadRequest
		message = "invalid LINE auth callback request"
		code = "invalid_auth_callback"
	case errors.Is(err, auth.ErrLineTokenExchange), errors.Is(err, auth.ErrLineTokenVerification):
		statusCode = nethttp.StatusBadGateway
		message = "LINE auth provider request failed"
		code = "line_auth_failed"
	}

	c.JSON(statusCode, ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
			Details: []string{err.Error()},
		},
	})
}

func buildAuthDoneRedirect(frontendURL string, token string) string {
	parsedURL, err := url.Parse(frontendURL)
	if err != nil {
		return strings.TrimRight(frontendURL, "/") + "/auth/done#token=" + token
	}

	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/") + "/auth/done"
	parsedURL.Fragment = fmt.Sprintf("token=%s", token)
	return parsedURL.String()
}
