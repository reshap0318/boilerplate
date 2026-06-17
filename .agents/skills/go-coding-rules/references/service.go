package references

import (
	"context"
	"fmt"

	"github.com/reshap0318/go-project/internal/dtos"
	"github.com/reshap0318/go-project/internal/helpers"
	"github.com/reshap0318/go-project/internal/models"
	"github.com/reshap0318/go-project/internal/pkg/email"
	"github.com/reshap0318/go-project/internal/repositories"
)

// ============================================================
// CREATE — with logging + notification + transaction
// ============================================================
func (s *Services) PermissionCreate(ctx context.Context, req *dtos.PermissionRequest) (*models.Permission, error) {
	s.Logger.LogStart("PermissionCreate", "Creating permission: %s", req.Name)

	var result *models.Permission
	err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		permission := &models.Permission{
			Name:        req.Name,
			Description: req.Description,
		}
		var err error
		result, err = s.repo.Permission.Create(tx, permission)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.Logger.LogEndWithError("PermissionCreate", "Failed: %v", err)
		return nil, err
	}

	_ = s.NotificationCreate(ctx, &services.NotificationCreateParams{
		Type:    "success",
		Title:   "Permission Created",
		Message: fmt.Sprintf("New permission created: %s", result.Name),
		Data: map[string]interface{}{
			"id":   result.ID,
			"name": result.Name,
		},
	})

	s.Logger.LogEnd("PermissionCreate", "Permission created: %s (ID: %d)", req.Name, result.ID)
	return result, nil
}

// ============================================================
// GET ALL PAGINATED — with DTO conversion, preloads, nil opts handling
// ============================================================
func (s *Services) PermissionGetAllPaginated(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.PermissionDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}

	opts.Preloads = []string{"Roles"}

	result, err := s.repo.Permission.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	permissionDTOs := make([]dtos.PermissionDTO, len(result.Data))
	for i, p := range result.Data {
		permissionDTOs[i] = dtos.ToPermissionDTO(&p)
	}

	return &repositories.PagedResult[dtos.PermissionDTO]{
		Data:       permissionDTOs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// ============================================================
// GET ALL — returns all with DTO conversion (no pagination)
// ============================================================
func (s *Services) PermissionGetAll(ctx context.Context) ([]dtos.PermissionDTO, error) {
	permissions, err := s.repo.Permission.FindAll(nil, "Roles")
	if err != nil {
		return nil, err
	}

	permissionDTOs := make([]dtos.PermissionDTO, len(permissions))
	for i, p := range permissions {
		permissionDTOs[i] = dtos.ToPermissionDTO(&p)
	}
	return permissionDTOs, nil
}

// ============================================================
// GET BY ID — simple read, NO LOG, returns DTO
// ============================================================
func (s *Services) PermissionGetByID(ctx context.Context, id uint) (*dtos.PermissionDTO, error) {
	permission, err := s.repo.Permission.FindByID(nil, id)
	if err != nil {
		return nil, err
	}
	dto := dtos.ToPermissionDTO(permission)
	return &dto, nil
}

// ============================================================
// UPDATE — with logging + notification + transaction
// ============================================================
func (s *Services) PermissionUpdate(ctx context.Context, id uint, req *dtos.PermissionRequest) (*models.Permission, error) {
	s.Logger.LogStart("PermissionUpdate", "Updating permission %d", id)

	var result *models.Permission
	err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.repo.Permission.Update(tx, &models.Permission{ID: id}, &models.Permission{
			Name:        req.Name,
			Description: req.Description,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.Logger.LogEndWithError("PermissionUpdate", "Failed: %v", err)
		return nil, err
	}

	_ = s.NotificationCreate(ctx, &services.NotificationCreateParams{
		Type:    "info",
		Title:   "Permission Updated",
		Message: fmt.Sprintf("Permission updated: %s", result.Name),
		Data: map[string]interface{}{
			"id":   result.ID,
			"name": result.Name,
		},
	})

	s.Logger.LogEnd("PermissionUpdate", "Permission %d updated", id)
	return result, nil
}

// ============================================================
// DELETE — with logging + notification + transaction
// ============================================================
func (s *Services) PermissionDelete(ctx context.Context, id uint) error {
	s.Logger.LogStart("PermissionDelete", "Deleting permission %d", id)

	var permission *models.Permission
	err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		var err error
		permission, err = s.repo.Permission.FindByID(nil, id, "Roles")
		if err != nil {
			return err
		}

		_, err = s.repo.Permission.Delete(tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.Logger.LogEndWithError("PermissionDelete", "Failed: %v", err)
		return err
	}

	_ = s.NotificationCreate(ctx, &services.NotificationCreateParams{
		Type:    "warning",
		Title:   "Permission Deleted",
		Message: fmt.Sprintf("Permission deleted: %s", permission.Name),
		Data: map[string]interface{}{
			"id":   permission.ID,
			"name": permission.Name,
		},
	})

	s.Logger.LogEnd("PermissionDelete", "Permission %d deleted", id)
	return nil
}
