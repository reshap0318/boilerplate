package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/reshap0318/go-boilerplate/internal/di"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/middleware"
	"github.com/reshap0318/go-boilerplate/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	host := helpers.GetEnv("APP_HOST", "0.0.0.0")
	port := helpers.GetEnv("APP_PORT", "8080")
	trustedProxies := helpers.GetEnv("TRUSTED_PROXIES", "")
	allowedOrigins := helpers.GetEnv("ALLOWED_ORIGINS", "*")

	gin.SetMode(helpers.GetEnv("GIN_MODE", "release"))

	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	r := gin.Default()

	if trustedProxies != "" {
		if err := r.SetTrustedProxies(strings.Split(trustedProxies, ",")); err != nil {
			log.Printf("Warning: failed to set trusted proxies: %v", err)
		}
	}

	r.Use(middleware.RateLimit(container.RateLimiter))
	r.Use(middleware.CORS(allowedOrigins))

	apiGroup := r.Group("/api")
	{
		routes.RegisterHealthRoutes(apiGroup, container.Handlers)
		routes.RegisterAuthRoutes(apiGroup, container.Handlers)
	}

	protected := apiGroup.Group("")
	protected.Use(middleware.JWTAuth(container.Services))
	{
		routes.RegisterAuthProtectedRoutes(protected, container.Handlers)
		routes.RegisterPermissionRoutes(protected, container.Handlers)
		routes.RegisterRoleRoutes(protected, container.Handlers)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	addr := host + ":" + port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
