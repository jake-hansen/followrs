package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jake-hansen/followrs/domain"
)

// HealthHandler presents information retrieved from a HealthService.
type HealthHandler struct {
	HealthService domain.HealthService // HealthService to use for performing operations on domain.
}

// NewHealthHandler initializes the endpoints for Health.
func NewHealthHandler(parentGroup *gin.RouterGroup, service domain.HealthService) {
	handler := &HealthHandler{
		HealthService: service,
	}

	healthGroup := parentGroup.Group("health")
	{
		healthGroup.GET("", handler.Status) // GET /health
	}
}

// Status retrieves the health status of the server.
func (h *HealthHandler) Status(c *gin.Context) {
	health, err := h.HealthService.GetHealth()
	if err == nil {
		c.JSON(http.StatusOK, health)
	} else {
		apiError := &APIError{
			Status:  http.StatusInternalServerError,
			Err:     err,
			Message: "An error occurred while retrieving the server's status",
		}
		c.Error(apiError).SetType(gin.ErrorTypePublic)
	}
}
