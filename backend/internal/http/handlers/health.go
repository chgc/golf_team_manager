package handlers

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func Health(c *gin.Context) {
	c.JSON(nethttp.StatusOK, HealthResponse{
		Status: "ok",
	})
}
