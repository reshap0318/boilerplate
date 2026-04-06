package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/handlers"
)

// RegisterPermissionRoutes registers protected permission routes.
func RegisterPermissionRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
	permissions := r.Group("/permissions")
	{
		permissions.POST("", handlers.CreatePermission)
		permissions.GET("", handlers.GetAllPermissions)
		permissions.GET("/:id", handlers.GetPermissionByID)
		permissions.PUT("/:id", handlers.UpdatePermission)
		permissions.DELETE("/:id", handlers.DeletePermission)
	}
}
