package dtos

import (
	"github.com/reshap0318/go-boilerplate/internal/models"
	"github.com/reshap0318/go-boilerplate/internal/repositories"
)

// UserRequest represents the request to create or update a user.
type UserRequest struct {
	Name                 string `json:"name" binding:"required,min=2,max=100"`
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
	Roles                []uint `json:"roles"`
}

// UserPageResult represents paginated user response.
type UserPageResult struct {
	Data       []UserDTO `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
}

// ToUserPageResult converts PagedResult[User] to UserPageResult with roles.
func ToUserPageResult(pageResult *repositories.PagedResult[models.User]) UserPageResult {
	return UserPageResult{
		Data:       ToUserDTOList(pageResult.Data),
		Total:      pageResult.Total,
		Page:       pageResult.Page,
		PageSize:   pageResult.PageSize,
		TotalPages: pageResult.TotalPages,
	}
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

	if len(u.UserHasRoles) > 0 {
		dto.Roles = make([]RoleDTO, len(u.UserHasRoles))
		for i, ur := range u.UserHasRoles {
			dto.Roles[i] = ToRoleDTO(&ur.Role)
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
