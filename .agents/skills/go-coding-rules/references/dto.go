package references

// ============================================================
// DTO — Request/Response structs
// ============================================================
type PermissionRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,max=255"`
	UserID      uint   `json:"-"` // from JWT context, not from request body
}

type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
