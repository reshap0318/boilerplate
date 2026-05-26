package references

import "github.com/reshap0318/go-project/internal/models"

// ============================================================
// Repository registration — add to 00_repository.go
// ============================================================
type Repositories struct {
	DB         *gorm.DB
	Permission *PermissionRepository
	// ... other repos
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		DB:         db,
		Permission: NewPermissionRepository(db),
	}
}

// ============================================================
// Custom repository — extends GenericRepository
// ============================================================
type PermissionRepository struct {
	*GenericRepository[models.Permission]
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		GenericRepository: NewGenericRepository[models.Permission](db),
	}
}

// Add custom methods here if needed (beyond GenericRepository):
// func (r *PermissionRepository) FindByRole(tx *gorm.DB, roleID uint) ([]models.Permission, error) {
//     var permissions []models.Permission
//     err := r.DB.Model(&models.Permission{}).
//         Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
//         Where("role_permissions.role_id = ?", roleID).
//         Find(&permissions).Error
//     return permissions, err
// }
