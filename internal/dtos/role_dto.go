package dtos

import "github.com/reshap0318/go-boilerplate/internal/models"

// RoleRequest represents the request to create or update a role.
type RoleRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty,max=255"`
	Permissions []uint  `json:"permissions" binding:"required"`
}

// RoleDTO represents role data transfer object.
type RoleDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// ToRoleDTO converts Role model to RoleDTO.
func ToRoleDTO(r *models.Role) RoleDTO {
	return RoleDTO{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}
}

// ToRoleDTOList converts a slice of Role models to RoleDTOs.
func ToRoleDTOList(roles []models.Role) []RoleDTO {
	result := make([]RoleDTO, len(roles))
	for i, r := range roles {
		result[i] = ToRoleDTO(&r)
	}
	return result
}
