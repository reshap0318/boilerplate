package dtos

// LoginRequest represents the login request payload.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse represents the login response payload.
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         UserDTO `json:"user"`
}

// RefreshTokenRequest represents the refresh token request payload.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RegisterRequest represents the register request payload.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required,min=3"`
}

// UserDTO represents user data transfer object.
type UserDTO struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ToUserDTO converts User model to UserDTO.
func ToUserDTO(email, name string, id uint) UserDTO {
	return UserDTO{
		ID:    id,
		Email: email,
		Name:  name,
	}
}

// ForgetPasswordRequest represents the forget password request payload.
type ForgetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	FrontendURL string `json:"frontend_url" binding:"required"`
}

// ResetPasswordRequest represents the reset password request payload.
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// MessageResponse represents a simple message response.
type MessageResponse struct {
	Message string `json:"message"`
}
