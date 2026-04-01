package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/reshap0318/go-boilerplate/cmd/migration/seeders"
	"github.com/reshap0318/go-boilerplate/internal/database"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/models"
	"gorm.io/gorm"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get database configuration from environment
	dbConnection := helpers.GetEnv("DB_CONNECTION", "mysql")
	dbHost := helpers.GetEnv("DB_HOST", "127.0.0.1")
	dbPort := helpers.GetEnv("DB_PORT", "3306")
	dbUser := helpers.GetEnv("DB_USERNAME", "root")
	dbPassword := helpers.GetEnv("DB_PASSWORD", "")
	dbName := helpers.GetEnv("DB_DATABASE", "boilerplate")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		log.Fatal("Database configuration is incomplete. Please set DB_HOST, DB_PORT, DB_USERNAME, DB_PASSWORD, and DB_DATABASE")
	}

	// Parse command
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up":
		runMigration(dbConnection, dbHost, dbPort, dbUser, dbPassword, dbName, "up")
	case "down":
		runMigration(dbConnection, dbHost, dbPort, dbUser, dbPassword, dbName, "down")
	case "seed":
		runSeed(dbConnection, dbHost, dbPort, dbUser, dbPassword, dbName)
	case "refresh":
		runMigration(dbConnection, dbHost, dbPort, dbUser, dbPassword, dbName, "down")
		runMigration(dbConnection, dbHost, dbPort, dbUser, dbPassword, dbName, "up")
		runSeed(dbConnection, dbHost, dbPort, dbUser, dbPassword, dbName)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run cmd/migration/main.go <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up       Apply all migrations (create/update tables)")
	fmt.Println("  down     Drop all migrated tables")
	fmt.Println("  seed     Insert default data")
	fmt.Println("  refresh  Drop all tables, recreate, and seed data")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migration/main.go up")
	fmt.Println("  go run cmd/migration/main.go down")
	fmt.Println("  go run cmd/migration/main.go seed")
	fmt.Println("  go run cmd/migration/main.go refresh")
}

func runMigration(dbConnection, dbHost, dbPort, dbUser, dbPassword, dbName string, command string) {
	var db *gorm.DB
	var err error

	if dbConnection == "postgres" || dbConnection == "postgresql" {
		db, err = database.NewPostgreSQL(database.PostgreSQLConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
		})
	} else {
		db, err = database.NewMySQL(database.MySQLConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
		})
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	switch command {
	case "up":
		fmt.Println("Running migrations...")

		// AutoMigrate all models
		err := db.AutoMigrate(
			&models.User{},
			&models.PasswordReset{},
		)
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		fmt.Println("✓ Migration completed successfully!")

	case "down":
		fmt.Println("Rolling back migrations...")

		// Drop tables in correct order (foreign key constraints)
		err := db.Migrator().DropTable(
			&models.PasswordReset{},
			&models.User{},
		)
		if err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}

		fmt.Println("✓ Rollback completed successfully!")
	}
}

func runSeed(dbConnection, dbHost, dbPort, dbUser, dbPassword, dbName string) {
	var db *gorm.DB
	var err error

	if dbConnection == "postgres" || dbConnection == "postgresql" {
		db, err = database.NewPostgreSQL(database.PostgreSQLConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
		})
	} else {
		db, err = database.NewMySQL(database.MySQLConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
		})
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("\n🌱 Seeding default data...\n")

	// Seed in correct order
	seeders.SeedUsers(db)

	fmt.Println("\n✅ Seeding completed!")
}
