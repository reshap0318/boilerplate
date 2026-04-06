package services

import (
	"context"

	"gorm.io/gorm"

	"github.com/reshap0318/go-boilerplate/internal/dtos"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/models"
)

// RoleCreate creates a new role with permissions.
func (s *Services) RoleCreate(ctx context.Context, req dtos.RoleRequest) (*dtos.RoleDTO, error) {
	s.Logger.LogStart("RoleCreate", "Creating role: %s", req.Name)

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

		for _, permID := range req.Permissions {
			rolePerm := &models.RoleHasPermission{
				RoleID:       result.ID,
				PermissionID: permID,
			}
			if _, err := s.repo.RoleHasPerm.Create(tx, rolePerm); err != nil {
				s.Logger.LogStep("RoleCreate", "Skipping duplicate or error for perm ID %d: %v", permID, err)
			}
		}
		return nil
	}); err != nil {
		s.Logger.LogEndWithError("RoleCreate", "Failed to create role: %v", err)
		return nil, err
	}

	dto := dtos.ToRoleDTO(result)
	s.Logger.LogEnd("RoleCreate", "Role created: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

// RoleGetAll returns all roles.
func (s *Services) RoleGetAll(ctx context.Context) ([]dtos.RoleDTO, error) {
	roles, err := s.repo.Role.FindAll(nil)
	if err != nil {
		return nil, err
	}

	return dtos.ToRoleDTOList(roles), nil
}

// RoleGetByID returns a role by ID.
func (s *Services) RoleGetByID(ctx context.Context, id uint) (*dtos.RoleDTO, error) {
	role, err := s.repo.Role.FindByID(nil, id)
	if err != nil {
		return nil, helpers.ErrNotFound
	}

	dto := dtos.ToRoleDTO(role)
	return &dto, nil
}

// RoleUpdate updates an existing role with permissions.
func (s *Services) RoleUpdate(ctx context.Context, id uint, req dtos.RoleRequest) (*dtos.RoleDTO, error) {
	s.Logger.LogStart("RoleUpdate", "Updating role ID: %d", id)

	role := &models.Role{ID: id}
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != nil {
		role.Description = req.Description
	}

	var result *models.Role
	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.repo.Role.Update(tx, &models.Role{ID: id}, role)
		if err != nil {
			return err
		}

		// Replace all permissions
		if err := s.repo.RoleHasPerm.DeleteByRoleID(tx, result.ID); err != nil {
			return err
		}

		for _, permID := range req.Permissions {
			rolePerm := &models.RoleHasPermission{
				RoleID:       result.ID,
				PermissionID: permID,
			}
			if _, err := s.repo.RoleHasPerm.Create(tx, rolePerm); err != nil {
				s.Logger.LogStep("RoleUpdate", "Skipping duplicate or error for perm ID %d: %v", permID, err)
			}
		}

		return nil
	}); err != nil {
		s.Logger.LogEndWithError("RoleUpdate", "Failed to update role: %v", err)
		return nil, err
	}

	dto := dtos.ToRoleDTO(result)
	s.Logger.LogEnd("RoleUpdate", "Role updated: %s (ID: %d)", dto.Name, dto.ID)
	return &dto, nil
}

// RoleDelete soft deletes a role.
func (s *Services) RoleDelete(ctx context.Context, id uint) error {
	s.Logger.LogStart("RoleDelete", "Deleting role ID: %d", id)

	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		if err := s.repo.RoleHasPerm.DeleteByRoleID(tx, id); err != nil {
			return err
		}
		_, err := s.repo.Role.Delete(tx, id)
		return err
	}); err != nil {
		s.Logger.LogEndWithError("RoleDelete", "Failed to delete role: %v", err)
		return err
	}

	s.Logger.LogEnd("RoleDelete", "Role deleted: ID: %d", id)
	return nil
}

// RoleGetPermissions returns all permissions for a role.
func (s *Services) RoleGetPermissions(ctx context.Context, roleID uint) ([]dtos.PermissionDTO, error) {
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
