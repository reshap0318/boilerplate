package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/handlers"
)

// RegisterJWKSRoutes registers JWKS routes
func RegisterJWKSRoutes(r *gin.Engine, handlers *handlers.Handlers) {
	// JWKS endpoint - public, no auth required
	r.GET("/.well-known/jwks.json", handlers.JWKSGetKeys)
}
