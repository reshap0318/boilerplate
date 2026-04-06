package seeders

import (
	"fmt"
	"log"

	"github.com/reshap0318/go-boilerplate/internal/models"
	"gorm.io/gorm"
)

// SeedRolePermissions maps roles to their permissions
func SeedRolePermissions(db *gorm.DB, roleIDs map[string]uint, permIDs map[string]uint) {
	fmt.Println("Seeding role permissions...")

	rolePerms := map[string][]string{
		"Super Admin": {
			// All permissions
			"user.create", "user.read", "user.update", "user.delete",
			"role.create", "role.read", "role.update", "role.delete",
			"permission.create", "permission.read", "permission.update", "permission.delete",
		},
		"Admin": {
			"user.create", "user.read", "user.update", "user.delete",
			"role.create", "role.read", "role.update", "role.delete",
			"permission.read",
		},
		"Editor": {
			"user.create", "user.read", "user.update",
			"permission.read",
		},
		"Viewer": {
			"user.read",
			"role.read",
			"permission.read",
		},
	}

	count := 0
	for roleName, permNames := range rolePerms {
		roleID, roleOK := roleIDs[roleName]
		if !roleOK {
			log.Printf("Role %s not found, skipping", roleName)
			continue
		}

		for _, permName := range permNames {
			permID, permOK := permIDs[permName]
			if !permOK {
				log.Printf("Permission %s not found, skipping", permName)
				continue
			}

			var existing models.RoleHasPermission
			err := db.Where("role_id = ? AND permission_id = ?", roleID, permID).First(&existing).Error
			if err == nil {
				continue
			}

			rp := models.RoleHasPermission{
				RoleID:       roleID,
				PermissionID: permID,
			}

			if err := db.Create(&rp).Error; err != nil {
				log.Printf("Failed to create role_permission for %s-%s: %v", roleName, permName, err)
			} else {
				count++
			}
		}
	}

	fmt.Printf("✓ Seeded %d role permissions\n", count)
}
