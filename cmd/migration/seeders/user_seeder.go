package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/models"
	"gorm.io/gorm"
)

// SeedUsers inserts default user data
func SeedUsers(db *gorm.DB) {
	fmt.Println("Seeding users...")

	// Default users to seed
	defaultUsers := []struct {
		Email    string
		Password string
		Name     string
	}{
		{
			Email:    "admin@example.com",
			Password: "admin123",
			Name:     "Admin User",
		},
		{
			Email:    "user@example.com",
			Password: "user123",
			Name:     "Regular User",
		},
		{
			Email:    "test@example.com",
			Password: "test123",
			Name:     "Test User",
		},
	}

	count := 0
	for _, userData := range defaultUsers {
		// Check if user already exists
		var existing models.User
		result := db.Where("email = ?", userData.Email).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			// User doesn't exist, create it
			hashedPassword, err := helpers.HashString(userData.Password)
			if err != nil {
				log.Printf("Failed to hash password for %s: %v", userData.Email, err)
				continue
			}

			user := models.User{
				Email:    userData.Email,
				Password: hashedPassword,
				Name:     userData.Name,
			}

			if err := db.Create(&user).Error; err != nil {
				log.Printf("Failed to create user %s: %v", userData.Email, err)
			} else {
				count++
			}
		} else if result.Error != nil {
			log.Printf("Failed to check user %s: %v", userData.Email, result.Error)
		} else {
			// User already exists, skip
			fmt.Printf("  ⊘ User %s already exists, skipping\n", userData.Email)
		}
	}

	fmt.Printf("✓ Seeded %d users\n", count)
}
