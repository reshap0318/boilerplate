package repositories

import (
	"github.com/reshap0318/go-boilerplate/internal/models"
	"gorm.io/gorm"
)

// RoleHasPermissionRepository provides database operations for RoleHasPermission model.
type RoleHasPermissionRepository struct {
	*GenericRepository[models.RoleHasPermission]
}

// NewRoleHasPermissionRepository creates a new RoleHasPermissionRepository.
func NewRoleHasPermissionRepository(db *gorm.DB) *RoleHasPermissionRepository {
	return &RoleHasPermissionRepository{
		GenericRepository: NewGenericRepository(db, &models.RoleHasPermission{}),
	}
}
