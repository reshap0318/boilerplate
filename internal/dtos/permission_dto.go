package dtos

import "github.com/reshap0318/go-boilerplate/internal/models"

// CreatePermissionRequest represents the request to create a permission.
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"max=255"`
}

// UpdatePermissionRequest represents the request to update a permission.
type UpdatePermissionRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=100"`
	Description string `json:"description" binding:"max=255"`
}

// PermissionDTO represents permission data transfer object.
type PermissionDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ToPermissionDTO converts Permission model to PermissionDTO.
func ToPermissionDTO(p *models.Permission) PermissionDTO {
	return PermissionDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
	}
}

// ToPermissionDTOList converts a slice of Permission models to PermissionDTOs.
func ToPermissionDTOList(permissions []models.Permission) []PermissionDTO {
	result := make([]PermissionDTO, len(permissions))
	for i, p := range permissions {
		result[i] = ToPermissionDTO(&p)
	}
	return result
}
