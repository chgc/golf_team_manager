package apihttp

import (
	"github.com/chgc/golf_team_manager/backend/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", handlers.Health)

	return router
}
