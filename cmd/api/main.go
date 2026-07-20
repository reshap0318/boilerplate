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

	r.Static("/storage", "./storage")

	r.NoRoute(func(c *gin.Context) {
		helpers.NotFound(c, "Endpoint not found")
	})

	apiGroup := r.Group("/api")
	protected := apiGroup.Group("")
	protected.Use(middleware.JWTAuth(container.Services))

	routes.RegisterAll(r, apiGroup, protected, container.Handlers, container.Access)

	addr := host + ":" + port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
