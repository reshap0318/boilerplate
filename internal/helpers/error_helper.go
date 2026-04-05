package helpers

import "errors"

// Common application errors.
var (
	ErrNotFound          = errors.New("record not found")
	ErrInvalidToken      = errors.New("invalid token")
	ErrExpiredToken      = errors.New("token expired")
	ErrInvalidCredential = errors.New("invalid email or password")
	ErrUserExists        = errors.New("user already exists")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrTokenExpired      = errors.New("reset token has expired")
	ErrTokenUsed         = errors.New("reset token has already been used")
	ErrTokenInvalid      = errors.New("invalid reset token")
)
