package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"message": "Service is healthy",
	})
}

func (h *HealthHandler) RegisterPublicRoutes(r *gin.Engine) {
	r.GET("/health", h.HealthCheck)
}

func (h *HealthHandler) RegisterProtectedRoutes(r *gin.RouterGroup) {}
