package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/handlers"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/middleware"
)

// RegisterUserRoutes registers protected user routes.
func RegisterUserRoutes(r *gin.RouterGroup, h *handlers.Handlers, acc *helpers.Access) {
	users := r.Group("/users")
	{
		users.POST("", middleware.RequirePermission(acc, "user.create"), h.UserCreate)
		users.GET("", middleware.RequirePermission(acc, "user.index"), h.UserGetAll)
		users.GET("/:id", middleware.RequirePermission(acc, "user.index"), h.UserGetByID)
		users.PUT("/:id", middleware.RequirePermission(acc, "user.edit"), h.UserUpdate)
		users.DELETE("/:id", middleware.RequirePermission(acc, "user.delete"), h.UserDelete)
	}

	me := r.Group("/me")
	{
		me.GET("", h.ProfileGet)
		me.PUT("", h.ProfileUpdate)
	}
}
