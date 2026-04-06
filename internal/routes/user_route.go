package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/handlers"
)

// RegisterUserRoutes registers protected user routes.
func RegisterUserRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
	users := r.Group("/users")
	{
		users.POST("", handlers.UserCreate)
		users.GET("", handlers.UserGetAll)
		users.GET("/:id", handlers.UserGetByID)
		users.PUT("/:id", handlers.UserUpdate)
		users.DELETE("/:id", handlers.UserDelete)
	}
}
