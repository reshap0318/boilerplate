package services

import (
	"context"

	"gorm.io/gorm"

	"github.com/reshap0318/go-boilerplate/internal/dtos"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/models"
)

// CreateRole creates a new role.
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
		return err
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

// UpdateRole updates an existing role.
func (s *Services) UpdateRole(ctx context.Context, id uint, req dtos.RoleRequest) (*dtos.RoleDTO, error) {
	s.Logger.LogStart("UpdateRole", "Updating role ID: %d", id)

	var result *models.Role
	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		role, err := s.repo.Role.FindByID(tx, id)
		if err != nil {
			s.Logger.LogEndWithError("UpdateRole", "Role not found: %v", err)
			return helpers.ErrNotFound
		}

		s.Logger.LogStep("UpdateRole", "Role found: %s", role.Name)

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
		return err
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
		_, err := s.repo.Role.FindByID(tx, id)
		if err != nil {
			s.Logger.LogEndWithError("DeleteRole", "Role not found: %v", err)
			return helpers.ErrNotFound
		}

		s.Logger.LogStep("DeleteRole", "Role found, deleting...")

		_, err = s.repo.Role.Delete(tx, id)
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

// AttachPermissions attaches permissions to a role.
func (s *Services) AttachPermissions(ctx context.Context, roleID uint, req dtos.AttachPermissionsRequest) error {
	s.Logger.LogStart("AttachPermissions", "Attaching %d permissions to role ID: %d", len(req.PermissionIDs), roleID)

	_, err := s.repo.Role.FindByID(nil, roleID)
	if err != nil {
		s.Logger.LogEndWithError("AttachPermissions", "Role not found: %v", err)
		return helpers.ErrNotFound
	}

	s.Logger.LogStep("AttachPermissions", "Role found, attaching permissions...")

	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		for _, permID := range req.PermissionIDs {
			rolePerm := &models.RoleHasPermission{
				RoleID:       roleID,
				PermissionID: permID,
			}
			if _, err := s.repo.RoleHasPerm.Create(tx, rolePerm); err != nil {
				s.Logger.LogStep("AttachPermissions", "Skipping duplicate or error for perm ID %d: %v", permID, err)
				continue
			}
		}
		return nil
	}); err != nil {
		s.Logger.LogEndWithError("AttachPermissions", "Failed to attach permissions: %v", err)
		return err
	}

	s.Logger.LogEnd("AttachPermissions", "Attached %d permissions to role ID: %d", len(req.PermissionIDs), roleID)
	return nil
}

// DetachPermissions detaches all permissions from a role.
func (s *Services) DetachPermissions(ctx context.Context, roleID uint) error {
	s.Logger.LogStart("DetachPermissions", "Detaching all permissions from role ID: %d", roleID)

	if err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		return tx.Where("role_id = ?", roleID).Delete(&models.RoleHasPermission{}).Error
	}); err != nil {
		s.Logger.LogEndWithError("DetachPermissions", "Failed to detach permissions: %v", err)
		return err
	}

	s.Logger.LogEnd("DetachPermissions", "All permissions detached from role ID: %d", roleID)
	return nil
}
