package apihttp

import (
	"database/sql"
	nethttp "net/http"
	"strings"

	"github.com/chgc/golf_team_manager/backend/internal/auth"
	"github.com/chgc/golf_team_manager/backend/internal/config"
	"github.com/chgc/golf_team_manager/backend/internal/http/handlers"
	"github.com/chgc/golf_team_manager/backend/internal/http/middleware"
	"github.com/chgc/golf_team_manager/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	LineProvider auth.LineProvider
	TokenManager auth.TokenManager
}

func NewRouter(database *sql.DB, cfg config.Config) *gin.Engine {
	return newRouterWithDependencies(database, cfg, RouterDependencies{})
}

func newRouterWithDependencies(database *sql.DB, cfg config.Config, deps RouterDependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.RequestID())

	apiHandlers := NewHandlers(database)
	apiPublicGroup := router.Group("/api")
	apiProtectedGroup := router.Group("/api")

	switch cfg.Auth.Mode {
	case "line":
		if deps.TokenManager == nil {
			deps.TokenManager = auth.NewHMACTokenManager(cfg.Auth.JWTSecret)
		}
		if deps.LineProvider == nil {
			deps.LineProvider = auth.NewHTTPLineProvider(nethttp.DefaultClient, auth.LINEConfig{
				ClientID:     cfg.Auth.LineClientID,
				ClientSecret: cfg.Auth.LineClientSecret,
				RedirectURI:  cfg.Auth.LineRedirectURI,
			})
		}

		lineService := auth.NewLineAuthService(
			deps.LineProvider,
			repository.NewSQLiteUserRepository(database),
			deps.TokenManager,
			cfg.Auth.JWTTTL,
		)
		lineHandlers := NewLineAuthHandlers(lineService, cfg.Auth.FrontendURL, shouldSecureAuthCookie(cfg.Auth))
		apiPublicGroup.GET("/auth/line/login", lineHandlers.Login)
		apiPublicGroup.GET("/auth/line/callback", lineHandlers.Callback)
		apiProtectedGroup.Use(middleware.JWTAuth(deps.TokenManager))
	default:
		apiProtectedGroup.Use(middleware.DevelopmentAuth(cfg.Auth))
	}

	router.GET("/health", handlers.Health)

	apiProtectedGroup.GET("/auth/me", apiHandlers.GetCurrentPrincipal)
	apiProtectedGroup.GET("/players", apiHandlers.ListPlayers)
	apiProtectedGroup.GET("/players/:playerId", apiHandlers.GetPlayerByID)
	apiProtectedGroup.POST("/players", apiHandlers.CreatePlayer)
	apiProtectedGroup.PATCH("/players/:playerId", apiHandlers.UpdatePlayer)
	apiProtectedGroup.GET("/sessions", apiHandlers.ListSessions)
	apiProtectedGroup.GET("/sessions/:sessionId", apiHandlers.GetSessionByID)
	apiProtectedGroup.POST("/sessions", apiHandlers.CreateSession)
	apiProtectedGroup.PATCH("/sessions/:sessionId", apiHandlers.UpdateSession)
	apiProtectedGroup.GET("/sessions/:sessionId/registrations", apiHandlers.ListRegistrationsBySession)
	apiProtectedGroup.POST("/sessions/:sessionId/registrations", apiHandlers.CreateRegistration)
	apiProtectedGroup.PATCH("/registrations/:registrationId", apiHandlers.UpdateRegistration)
	apiProtectedGroup.GET("/reports/sessions/:sessionId/reservation-summary", apiHandlers.GetReservationSummary)
	apiProtectedGroup.GET("/admin/users", apiHandlers.ListAdminUsers)
	apiProtectedGroup.PATCH("/admin/users/:userId", apiHandlers.UpdateAdminUser)

	return router
}

func shouldSecureAuthCookie(cfg config.AuthConfig) bool {
	return strings.HasPrefix(strings.ToLower(cfg.LineRedirectURI), "https://") || strings.HasPrefix(strings.ToLower(cfg.FrontendURL), "https://")
}
