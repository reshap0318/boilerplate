package di

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	clientEmail "github.com/reshap0318/go-boilerplate/internal/clients/email"
	"github.com/reshap0318/go-boilerplate/internal/database"
	"github.com/reshap0318/go-boilerplate/internal/handlers"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/repositories"
	"github.com/reshap0318/go-boilerplate/internal/services"
)

// Container holds all dependencies.
type Container struct {
	DB           *gorm.DB
	Redis        *database.RedisCache
	EmailClient  *clientEmail.EmailClient
	Logger       *helpers.Logger
	RateLimiter  *helpers.RateLimiter
	Repositories *repositories.Repositories
	Services     *services.Services
	Handlers     *handlers.Handlers
}

// Close closes all connections.
func (c *Container) Close() error {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return fmt.Errorf("error getting database connection: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("error closing database connection: %w", err)
		}
		log.Println("Database connection closed")
	}

	if c.Redis != nil {
		log.Println("Redis connection closed")
	}

	if c.Logger != nil {
		c.Logger.Close()
		log.Println("Logger closed")
	}

	return nil
}

// NewContainer creates and initializes all dependencies.
func NewContainer() (*Container, error) {
	container := &Container{}

	// Initialize Logger (early, before other components)
	logger, err := helpers.NewLogger("logs")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	container.Logger = logger

	// Determine database connection type (default: mysql)
	dbConnection := helpers.GetEnv("DB_CONNECTION", "mysql")

	// Initialize database based on DB_CONNECTION
	if dbConnection == "postgres" || dbConnection == "postgresql" {
		postgres, err := database.NewPostgreSQL(database.PostgreSQLConfig{
			Host:     helpers.GetEnv("DB_HOST", "127.0.0.1"),
			Port:     helpers.GetEnv("DB_PORT", "5432"),
			User:     helpers.GetEnv("DB_USERNAME", "postgres"),
			Password: helpers.GetEnv("DB_PASSWORD", ""),
			DBName:   helpers.GetEnv("DB_DATABASE", "boilerplate"),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize PostgreSQL: %w", err)
		}
		container.DB = postgres
	} else {
		// Default to MySQL
		mysql, err := database.NewMySQL(database.MySQLConfig{
			Host:     helpers.GetEnv("DB_HOST", "127.0.0.1"),
			Port:     helpers.GetEnv("DB_PORT", "3306"),
			User:     helpers.GetEnv("DB_USERNAME", "root"),
			Password: helpers.GetEnv("DB_PASSWORD", ""),
			DBName:   helpers.GetEnv("DB_DATABASE", "boilerplate"),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize MySQL: %w", err)
		}
		container.DB = mysql
	}

	// Initialize Redis (optional)
	if helpers.GetEnv("REDIS_ENABLED", "false") == "true" {
		redisClient, err := database.NewRedis(database.RedisConfig{
			Host:     helpers.GetEnv("REDIS_HOST", "localhost"),
			Port:     helpers.GetEnv("REDIS_PORT", "6379"),
			Password: helpers.GetEnv("REDIS_PASSWORD", ""),
			DB:       helpers.GetEnvInt("REDIS_DB", 0),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Redis: %w", err)
		}
		container.Redis = database.NewRedisCache(redisClient)
	}

	// Initialize Email Client (optional)
	container.EmailClient = clientEmail.NewEmailClient()

	// Initialize Rate Limiter
	rateLimitRequests := helpers.GetEnvInt("RATE_LIMIT_REQUESTS", 100)
	rateLimitWindow := helpers.GetEnvInt("RATE_LIMIT_WINDOW", 60)
	container.RateLimiter = helpers.NewRateLimiter(rateLimitRequests, rateLimitWindow)
	log.Printf("Rate limiting enabled: %d requests per %d seconds", rateLimitRequests, rateLimitWindow)

	// DB is required for repositories and services
	if container.DB == nil {
		panic("database connection is required but not initialized. Set DB_CONNECTION=mysql or DB_CONNECTION=postgres in your .env file")
	}

	// Always initialize repositories
	repos, err := repositories.NewRepositories(container.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}
	container.Repositories = repos

	// Always initialize services (Redis can be nil)
	container.Services = services.NewServices(container.Repositories, container.Redis, container.EmailClient, container.Logger)

	// Always initialize handlers
	container.Handlers = handlers.NewHandlers(container.Services)

	return container, nil
}
