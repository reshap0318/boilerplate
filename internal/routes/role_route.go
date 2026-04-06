package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/handlers"
)

// RegisterRoleRoutes registers protected role routes.
func RegisterRoleRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
	roles := r.Group("/roles")
	{
		roles.POST("", handlers.RoleCreate)
		roles.GET("", handlers.RoleGetAll)
		roles.GET("/:id", handlers.RoleGetByID)
		roles.PUT("/:id", handlers.RoleUpdate)
		roles.DELETE("/:id", handlers.RoleDelete)
		roles.GET("/:id/permissions", handlers.RoleGetPermissions)
		roles.POST("/:id/permissions", handlers.RoleAttachPermissions)
		roles.DELETE("/:id/permissions", handlers.RoleDetachPermissions)
	}
}
