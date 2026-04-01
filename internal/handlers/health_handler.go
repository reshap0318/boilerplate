package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck handles health check endpoint.
// @Summary Health check
// @Description Check system health status (database, redis, etc)
// @Tags health
// @Produce json
// @Success 200 {object} dtos.HealthStatus
// @Router /api/health [get]
func (h *Handlers) HealthCheck(c *gin.Context) {
	status := h.svcs.HealthGetStatus(c.Request.Context())
	
	if status.Status == "unhealthy" {
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}
	
	c.JSON(http.StatusOK, status)
}
