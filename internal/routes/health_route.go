package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/handlers"
)

// RegisterHealthRoutes registers health check routes.
func RegisterHealthRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
	r.GET("/health", handlers.HealthCheck)
}
