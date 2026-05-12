package dtos

import (
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/models"
)

// UserCreateRequest represents the request to create a user.
type UserCreateRequest struct {
	Name                 string `json:"name" validate:"required,min=2,max=100"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	Avatar               string `json:"avatar"`
	Roles                []uint `json:"roles"`
}

// UserUpdateRequest represents the request to update a user.
type UserUpdateRequest struct {
	Name                 string `json:"name" validate:"required,min=2,max=100"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"omitempty,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"omitempty,eqfield=Password"`
	Avatar               string `json:"avatar"`
	Roles                []uint `json:"roles"`
}

// ToUserDTO converts User model to UserDTO.
func ToUserDTO(u *models.User) UserDTO {
	dto := UserDTO{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Avatar:      helpers.GetFileURL(u.Avatar),
		CreatedAt:   u.CreatedAt,
		Roles:       []RoleMiniDTO{},
		Permissions: []PermissionDTO{},
	}

	permSet := make(map[uint]bool)
	for _, r := range u.Roles {
		dto.Roles = append(dto.Roles, ToRoleMiniDTO(&r))
		for _, p := range r.Permissions {
			if !permSet[p.ID] {
				permSet[p.ID] = true
				dto.Permissions = append(dto.Permissions, ToPermissionDTO(&p))
			}
		}
	}

	return dto
}

// ToUserDTOList converts a slice of User models to UserDTOs.
func ToUserDTOList(users []models.User) []UserDTO {
	result := make([]UserDTO, len(users))
	for i, u := range users {
		result[i] = ToUserDTO(&u)
	}
	return result
}
