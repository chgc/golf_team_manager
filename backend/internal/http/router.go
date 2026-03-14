package apihttp

import (
	"database/sql"

	"github.com/chgc/golf_team_manager/backend/internal/config"
	"github.com/chgc/golf_team_manager/backend/internal/http/handlers"
	"github.com/chgc/golf_team_manager/backend/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(database *sql.DB, cfg config.Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.RequestID(), middleware.DevelopmentAuth(cfg.Auth))

	apiHandlers := NewHandlers(database)

	router.GET("/health", handlers.Health)

	apiGroup := router.Group("/api")
	apiGroup.GET("/auth/me", apiHandlers.GetCurrentPrincipal)
	apiGroup.GET("/players", apiHandlers.ListPlayers)
	apiGroup.GET("/players/:playerId", apiHandlers.GetPlayerByID)
	apiGroup.POST("/players", apiHandlers.CreatePlayer)
	apiGroup.PATCH("/players/:playerId", apiHandlers.UpdatePlayer)
	apiGroup.GET("/sessions", apiHandlers.ListSessions)
	apiGroup.POST("/sessions", apiHandlers.CreateSession)
	apiGroup.GET("/sessions/:sessionId/registrations", apiHandlers.ListRegistrationsBySession)
	apiGroup.POST("/sessions/:sessionId/registrations", apiHandlers.CreateRegistration)
	apiGroup.GET("/reports/sessions/:sessionId/reservation-summary", apiHandlers.NotImplemented)

	return router
}
