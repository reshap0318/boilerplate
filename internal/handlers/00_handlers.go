package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/reshap0318/go-boilerplate/internal/services"
)

// Handlers holds all HTTP handlers.
type Handlers struct {
	svcs     *services.Services
	Validate *validator.Validate
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(svcs *services.Services, validate *validator.Validate) *Handlers {
	return &Handlers{
		svcs:     svcs,
		Validate: validate,
	}
}
