package references

import (
	"github.com/gin-gonic/gin"
	"github.com/reshap0318/go-project/internal/dtos"
	"github.com/reshap0318/go-project/internal/helpers"
)

// ============================================================
// CREATE
// ============================================================
func (h *Handlers) PermissionCreate(c *gin.Context) {
	var req dtos.PermissionRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.PermissionCreate(c.Request.Context(), &req)
	if helpers.HandleError(c, err, "Failed to create permission") {
		return
	}

	helpers.Created(c, "Permission created successfully", result)
}

// ============================================================
// GET ALL
// ============================================================
func (h *Handlers) PermissionGetAll(c *gin.Context) {
	opts := &helpers.QueryOptions{
		Page:         helpers.ParseInt(c.Query("page"), 1),
		PageSize:     helpers.ParseInt(c.Query("per_page"), 10),
		SortBy:       c.DefaultQuery("sort", "id"),
		Order:        c.DefaultQuery("order", "ASC"),
		Search:       c.Query("search"),
		SearchFields: []string{"name", "description"},
	}

	result, err := h.svcs.PermissionGetAll(c.Request.Context(), opts)
	if helpers.HandleError(c, err, "Failed to fetch permissions") {
		return
	}

	helpers.OK(c, "Permissions retrieved successfully", result)
}

// ============================================================
// GET BY ID
// ============================================================
func (h *Handlers) PermissionGetByID(c *gin.Context) {
	id := helpers.ParseUint(c.Param("id"))
	if id == 0 {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	result, err := h.svcs.PermissionGetByID(c.Request.Context(), id)
	if helpers.HandleError(c, err, "Failed to fetch permission") {
		return
	}

	helpers.OK(c, "Permission retrieved successfully", result)
}

// ============================================================
// UPDATE
// ============================================================
func (h *Handlers) PermissionUpdate(c *gin.Context) {
	id := helpers.ParseUint(c.Param("id"))
	if id == 0 {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	var req dtos.PermissionRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.PermissionUpdate(c.Request.Context(), id, &req)
	if helpers.HandleError(c, err, "Failed to update permission") {
		return
	}

	helpers.OK(c, "Permission updated successfully", result)
}

// ============================================================
// DELETE
// ============================================================
func (h *Handlers) PermissionDelete(c *gin.Context) {
	id := helpers.ParseUint(c.Param("id"))
	if id == 0 {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	err := h.svcs.PermissionDelete(c.Request.Context(), id)
	if helpers.HandleError(c, err, "Failed to delete permission") {
		return
	}

	helpers.OK(c, "Permission deleted successfully", nil)
}
