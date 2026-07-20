package references

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/reshap0318/go-boilerplate/internal/dtos"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/repositories"
)

// NOTE: ALL service errors MUST be handled via helpers.HandleError().
// Do NOT manually check errors (e.g., err == helpers.ErrNotFound).
// HandleError automatically handles: FieldError, CustomError, ErrNotFound,
// ErrForbidden, ErrInvalidToken, etc. Fallback message for unknown errors.

// ============================================================
// CREATE
// ============================================================
func (h *Handlers) PermissionCreate(c *gin.Context) {
	var req dtos.PermissionRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.PermissionCreate(c.Request.Context(), req)
	if helpers.HandleError(c, err, "Failed to create permission") {
		return
	}

	helpers.Created(c, "Permission created successfully", result)
}

// ============================================================
// GET ALL — page_size omitted or -1 returns all records; page_size>0 paginates
// ============================================================
func (h *Handlers) PermissionGetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "-1"))

	if pageSize < 0 {
		page = 1
	}

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
		SortBy:   c.DefaultQuery("sort", "id"),
		Order:    c.DefaultQuery("order", "ASC"),
	}

	// Optional: keyword search across multiple fields using ConditionGroups
	if search := c.Query("search"); search != "" {
		opts.ConditionGroups = []repositories.ConditionGroup{
			{
				Logic: "OR",
				Conditions: []repositories.QueryCondition{
					{Column: "name", Operator: "LIKE", Value: "%" + search + "%"},
					{Column: "description", Operator: "LIKE", Value: "%" + search + "%"},
				},
			},
		}
	}

	result, err := h.svcs.PermissionGetAll(c.Request.Context(), opts)
	if helpers.HandleError(c, err, "Failed to fetch permissions") {
		return
	}

	helpers.OKWithMetadata(c, "Permissions retrieved successfully", result)
}

// ============================================================
// GET BY ID
// ============================================================
func (h *Handlers) PermissionGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	result, err := h.svcs.PermissionGetByID(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to fetch permission") {
		return
	}

	helpers.OK(c, "Permission retrieved successfully", result)
}

// ============================================================
// UPDATE
// ============================================================
func (h *Handlers) PermissionUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	var req dtos.PermissionRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.PermissionUpdate(c.Request.Context(), uint(id), req)
	if helpers.HandleError(c, err, "Failed to update permission") {
		return
	}

	helpers.OK(c, "Permission updated successfully", result)
}

// ============================================================
// DELETE
// ============================================================
func (h *Handlers) PermissionDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	err = h.svcs.PermissionDelete(c.Request.Context(), uint(id))
	if helpers.HandleError(c, err, "Failed to delete permission") {
		return
	}

	helpers.OK(c, "Permission deleted successfully", nil)
}
