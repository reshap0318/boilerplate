package dtos

import "github.com/reshap0318/go-boilerplate/internal/models"

// UserRequest represents the request to create or update a user.
type UserRequest struct {
	Name                 string `json:"name" binding:"required,min=2,max=100"`
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
	Roles                []uint `json:"roles"`
}

// ToUserDTO converts User model to UserDTO.
func ToUserDTO(u *models.User) UserDTO {
	dto := UserDTO{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		Roles:     []RoleDTO{},
	}

	for _, r := range u.Roles {
		dto.Roles = append(dto.Roles, ToRoleDTO(&r))
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
