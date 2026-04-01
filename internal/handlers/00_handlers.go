package handlers

import (
	"github.com/reshap0318/go-boilerplate/internal/services"
)

// Handlers holds all HTTP handlers.
type Handlers struct {
	svcs *services.Services
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(svcs *services.Services) *Handlers {
	return &Handlers{svcs: svcs}
}
