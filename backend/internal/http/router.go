package apihttp

import (
	"database/sql"

	"github.com/chgc/golf_team_manager/backend/internal/http/handlers"
	"github.com/chgc/golf_team_manager/backend/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(database *sql.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.RequestID())

	apiHandlers := NewHandlers(database)

	router.GET("/health", handlers.Health)

	apiGroup := router.Group("/api")
	apiGroup.GET("/players", apiHandlers.ListPlayers)
	apiGroup.POST("/players", apiHandlers.CreatePlayer)
	apiGroup.GET("/sessions", apiHandlers.ListSessions)
	apiGroup.POST("/sessions", apiHandlers.CreateSession)
	apiGroup.GET("/sessions/:sessionId/registrations", apiHandlers.ListRegistrationsBySession)
	apiGroup.POST("/sessions/:sessionId/registrations", apiHandlers.CreateRegistration)
	apiGroup.GET("/reports/sessions/:sessionId/reservation-summary", apiHandlers.NotImplemented)

	return router
}
