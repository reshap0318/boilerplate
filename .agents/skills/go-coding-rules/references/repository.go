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

// ============================================================
// Priority: Use Generic Methods First
// ============================================================
// Prefer using existing GenericRepository methods over creating custom repository methods.
//
// Guidelines:
// 1. Check if generic method can solve the problem first (see SKILL.md table)
// 2. Create custom method when generic methods cannot handle the requirement
// 3. Common scenarios that may need custom methods:
//    - Complex JOIN queries
//    - Aggregate functions (SUM, COUNT with GROUP BY)
//    - Subqueries or nested conditions
//    - Database-specific features not covered by generic methods
//
// Examples:
//
// PREFERRED - Use generic method:
//     permissions, _ := s.repo.Permission.FindByFieldMap(nil, map[string]interface{}{
//         "resource": "user",
//         "action":   "create",
//     })
//
// ACCEPTABLE - Custom method needed (JOIN query):
//     func (r *PermissionRepository) FindByRoleID(tx *gorm.DB, roleID uint) ([]models.Permission, error) {
//         var permissions []models.Permission
//         err := r.DB.Model(&models.Permission{}).
//             Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
//             Where("role_permissions.role_id = ?", roleID).
//             Find(&permissions).Error
//         return permissions, err
//     }
//
// NOT RECOMMENDED - Custom method when generic can handle it:
//     func (r *PermissionRepository) FindByResource(tx *gorm.DB, resource string) ([]models.Permission, error) {
//         var permissions []models.Permission
//         err := r.DB.Where("resource = ?", resource).Find(&permissions).Error
//         return permissions, err
//     }
//     Better to use: s.repo.Permission.FindByFieldMap(nil, map[string]interface{}{"resource": resource})
