package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/go-boilerplate/internal/models"
	"gorm.io/gorm"
)

// SeedPermissions inserts default permission data
func SeedPermissions(db *gorm.DB) map[string]uint {
	fmt.Println("Seeding permissions...")

	permissions := []string{
		// User
		"user.create", "user.read", "user.update", "user.delete",
		// Role
		"role.create", "role.read", "role.update", "role.delete",
		// Permission
		"permission.create", "permission.read", "permission.update", "permission.delete",
	}

	resultMap := make(map[string]uint)

	for _, name := range permissions {
		var existing models.Permission
		err := db.Where("name = ?", name).First(&existing).Error
		if err == nil {
			resultMap[name] = existing.ID
			fmt.Printf("  ⊘ Permission %s already exists, skipping\n", name)
			continue
		}

		perm := models.Permission{
			Name:        name,
			Description: strPtr("Permission for " + name),
		}

		if err := db.Create(&perm).Error; err != nil {
			log.Printf("Failed to create permission %s: %v", name, err)
		} else {
			resultMap[name] = perm.ID
		}
	}

	fmt.Printf("✓ Seeded %d permissions\n", len(resultMap))
	return resultMap
}

func strPtr(s string) *string {
	return &s
}
