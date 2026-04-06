package services

import (
	"context"

	"gorm.io/gorm"

	"github.com/reshap0318/go-boilerplate/internal/dtos"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/models"
)

// CreateRole creates a new role with permissions.
func (s *Services) CreateRole(ctx context.Context, req dtos.RoleRequest) (*dtos.RoleDTO, error) {
	s.Logger.LogStart("CreateRole", "Creating role: %s", req.Name)

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	var result *models.Role
	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.repo.Role.Create(tx, role)
		if err != nil {
			return err
		}

		for _, permID := range req.PermissionIDs {
			rolePerm := &models.RoleHasPermission{
				RoleID:       result.ID,
				PermissionID: permID,
			}
			if _, err := s.repo.RoleHasPerm.Create(tx, rolePerm); err != nil {
				s.Logger.LogStep("CreateRole", "Skipping duplicate or error for perm ID %d: %v", permID, err)
			}
		}
		return nil
	}); err != nil {
		s.Logger.LogEndWithError("CreateRole", "Failed to create role: %v", err)
		return nil, err
	}

	dto := dtos.ToRoleDTO(result)
	s.Logger.LogEnd("CreateRole", "Role created: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

// GetAllRoles returns all roles.
func (s *Services) GetAllRoles(ctx context.Context) ([]dtos.RoleDTO, error) {
	roles, err := s.repo.Role.FindAll(nil)
	if err != nil {
		return nil, err
	}

	return dtos.ToRoleDTOList(roles), nil
}

// GetRoleByID returns a role by ID.
func (s *Services) GetRoleByID(ctx context.Context, id uint) (*dtos.RoleDTO, error) {
	role, err := s.repo.Role.FindByID(nil, id)
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToRoleDTO(role)
	return &dto, nil
}

// UpdateRole updates an existing role with permissions.
func (s *Services) UpdateRole(ctx context.Context, id uint, req dtos.RoleRequest) (*dtos.RoleDTO, error) {
	s.Logger.LogStart("UpdateRole", "Updating role ID: %d", id)

	var result *models.Role
	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		role, err := s.repo.Role.FindByID(tx, id)
		if err != nil {
			return helpers.ErrNotFound
		}

		if req.Name != "" {
			role.Name = req.Name
		}
		if req.Description != nil {
			role.Description = req.Description
		}

		result, err = s.repo.Role.Update(tx, role, &models.Role{
			Name:        role.Name,
			Description: role.Description,
		})
		if err != nil {
			return err
		}

		// Replace all permissions
		if err := tx.Where("role_id = ?", result.ID).Delete(&models.RoleHasPermission{}).Error; err != nil {
			return err
		}

		for _, permID := range req.PermissionIDs {
			rolePerm := &models.RoleHasPermission{
				RoleID:       result.ID,
				PermissionID: permID,
			}
			if _, err := s.repo.RoleHasPerm.Create(tx, rolePerm); err != nil {
				s.Logger.LogStep("UpdateRole", "Skipping duplicate or error for perm ID %d: %v", permID, err)
			}
		}

		return nil
	}); err != nil {
		s.Logger.LogEndWithError("UpdateRole", "Failed to update role: %v", err)
		return nil, err
	}

	dto := dtos.ToRoleDTO(result)
	s.Logger.LogEnd("UpdateRole", "Role updated: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

// DeleteRole soft deletes a role.
func (s *Services) DeleteRole(ctx context.Context, id uint) error {
	s.Logger.LogStart("DeleteRole", "Deleting role ID: %d", id)

	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&models.RoleHasPermission{}).Error; err != nil {
			return err
		}
		_, err := s.repo.Role.Delete(tx, id)
		return err
	}); err != nil {
		s.Logger.LogEndWithError("DeleteRole", "Failed to delete role: %v", err)
		return err
	}

	s.Logger.LogEnd("DeleteRole", "Role deleted: ID: %d", id)
	return nil
}

// GetRolePermissions returns all permissions for a role.
func (s *Services) GetRolePermissions(ctx context.Context, roleID uint) ([]dtos.PermissionDTO, error) {
	var rolePerms []models.RoleHasPermission
	if err := s.repo.RoleHasPerm.DB.Where("role_id = ?", roleID).Find(&rolePerms).Error; err != nil {
		return nil, err
	}

	if len(rolePerms) == 0 {
		return []dtos.PermissionDTO{}, nil
	}

	permIDs := make([]uint, len(rolePerms))
	for i, rp := range rolePerms {
		permIDs[i] = rp.PermissionID
	}

	var permissions []models.Permission
	if err := s.repo.Permission.DB.Where("id IN ?", permIDs).Find(&permissions).Error; err != nil {
		return nil, err
	}

	return dtos.ToPermissionDTOList(permissions), nil
}
