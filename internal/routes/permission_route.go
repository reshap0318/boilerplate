package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/handlers"
)

// RegisterPermissionRoutes registers protected permission routes.
func RegisterPermissionRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
	permissions := r.Group("/permissions")
	{
		permissions.POST("", handlers.PermissionCreate)
		permissions.GET("", handlers.PermissionGetAll)
		permissions.GET("/:id", handlers.PermissionGetByID)
		permissions.PUT("/:id", handlers.PermissionUpdate)
		permissions.DELETE("/:id", handlers.PermissionDelete)
	}
}
