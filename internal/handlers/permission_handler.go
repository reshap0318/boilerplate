package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/dtos"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
)

// CreatePermission handles POST /api/permissions
func (h *Handlers) CreatePermission(c *gin.Context) {
	var req dtos.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	dto, err := h.svcs.CreatePermission(c.Request.Context(), req)
	if err != nil {
		helpers.InternalServerError(c, "Failed to create permission")
		return
	}

	helpers.Created(c, "Permission created successfully", dto)
}

// GetAllPermissions handles GET /api/permissions
func (h *Handlers) GetAllPermissions(c *gin.Context) {
	dtos, err := h.svcs.GetAllPermissions(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, "Failed to fetch permissions")
		return
	}

	helpers.OK(c, "Permissions fetched successfully", dtos)
}

// GetPermissionByID handles GET /api/permissions/:id
func (h *Handlers) GetPermissionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	dto, err := h.svcs.GetPermissionByID(c.Request.Context(), uint(id))
	if err != nil {
		helpers.NotFound(c, "Permission not found")
		return
	}

	helpers.OK(c, "Permission fetched successfully", dto)
}

// UpdatePermission handles PUT /api/permissions/:id
func (h *Handlers) UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	var req dtos.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	dto, err := h.svcs.UpdatePermission(c.Request.Context(), uint(id), req)
	if err != nil {
		helpers.NotFound(c, "Permission not found")
		return
	}

	helpers.OK(c, "Permission updated successfully", dto)
}

// DeletePermission handles DELETE /api/permissions/:id
func (h *Handlers) DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	err = h.svcs.DeletePermission(c.Request.Context(), uint(id))
	if err != nil {
		helpers.NotFound(c, "Permission not found")
		return
	}

	helpers.OK(c, "Permission deleted successfully", nil)
}
