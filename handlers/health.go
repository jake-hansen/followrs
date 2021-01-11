package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jake-hansen/followrs/domain"
)

type HealthHandler struct {
	HealthService domain.HealthService
}

// NewHealthHandler initializes the endpoints for Health
func NewHealthHandler(parentGroup *gin.RouterGroup, service domain.HealthService) {
	handler := &HealthHandler{
		HealthService: service,
	}

	healthGroup := parentGroup.Group("health")
	{
		healthGroup.GET("", handler.Status)
	}

}

// Status retrieves the health status of the server.
func (h *HealthHandler) Status(c *gin.Context) {
	health, _ := h.HealthService.GetHealth()
	c.JSON(http.StatusOK, health)
}
