package services

import (
	"context"

	"gorm.io/gorm"

	"github.com/reshap0318/go-boilerplate/internal/dtos"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/models"
)

// CreatePermission creates a new permission.
func (s *Services) CreatePermission(ctx context.Context, req dtos.PermissionRequest) (*dtos.PermissionDTO, error) {
	s.Logger.LogStart("CreatePermission", "Creating permission: %s", req.Name)

	permission := &models.Permission{
		Name:        req.Name,
		Description: req.Description,
	}

	var result *models.Permission
	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.repo.Permission.Create(tx, permission)
		return err
	}); err != nil {
		s.Logger.LogEndWithError("CreatePermission", "Failed to create permission: %v", err)
		return nil, err
	}

	dto := dtos.ToPermissionDTO(result)
	s.Logger.LogEnd("CreatePermission", "Permission created: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

// GetAllPermissions returns all permissions.
func (s *Services) GetAllPermissions(ctx context.Context) ([]dtos.PermissionDTO, error) {
	permissions, err := s.repo.Permission.FindAll(nil)
	if err != nil {
		return nil, err
	}

	return dtos.ToPermissionDTOList(permissions), nil
}

// GetPermissionByID returns a permission by ID.
func (s *Services) GetPermissionByID(ctx context.Context, id uint) (*dtos.PermissionDTO, error) {
	permission, err := s.repo.Permission.FindByID(nil, id)
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToPermissionDTO(permission)
	return &dto, nil
}

// UpdatePermission updates an existing permission.
func (s *Services) UpdatePermission(ctx context.Context, id uint, req dtos.PermissionRequest) (*dtos.PermissionDTO, error) {
	s.Logger.LogStart("UpdatePermission", "Updating permission ID: %d", id)

	var result *models.Permission
	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		permission, err := s.repo.Permission.FindByID(tx, id)
		if err != nil {
			s.Logger.LogEndWithError("UpdatePermission", "Permission not found: %v", err)
			return helpers.ErrNotFound
		}

		s.Logger.LogStep("UpdatePermission", "Permission found: %s", permission.Name)

		if req.Name != "" {
			permission.Name = req.Name
		}
		if req.Description != nil {
			permission.Description = req.Description
		}

		result, err = s.repo.Permission.Update(tx, permission, &models.Permission{
			Name:        permission.Name,
			Description: permission.Description,
		})
		return err
	}); err != nil {
		s.Logger.LogEndWithError("UpdatePermission", "Failed to update permission: %v", err)
		return nil, err
	}

	dto := dtos.ToPermissionDTO(result)
	s.Logger.LogEnd("UpdatePermission", "Permission updated: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

// DeletePermission soft deletes a permission.
func (s *Services) DeletePermission(ctx context.Context, id uint) error {
	s.Logger.LogStart("DeletePermission", "Deleting permission ID: %d", id)

	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		_, err := s.repo.Permission.FindByID(tx, id)
		if err != nil {
			s.Logger.LogEndWithError("DeletePermission", "Permission not found: %v", err)
			return helpers.ErrNotFound
		}

		s.Logger.LogStep("DeletePermission", "Permission found, deleting...")

		_, err = s.repo.Permission.Delete(tx, id)
		return err
	}); err != nil {
		s.Logger.LogEndWithError("DeletePermission", "Failed to delete permission: %v", err)
		return err
	}

	s.Logger.LogEnd("DeletePermission", "Permission deleted: ID: %d", id)
	return nil
}
