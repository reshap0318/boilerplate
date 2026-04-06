package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/dtos"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
)

// PermissionCreate handles POST /api/permissions
func (h *Handlers) PermissionCreate(c *gin.Context) {
	var req dtos.PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	dto, err := h.svcs.PermissionCreate(c.Request.Context(), req)
	if err != nil {
		helpers.InternalServerError(c, "Failed to create permission")
		return
	}

	helpers.Created(c, "Permission created successfully", dto)
}

// PermissionGetAll handles GET /api/permissions
func (h *Handlers) PermissionGetAll(c *gin.Context) {
	dtos, err := h.svcs.PermissionGetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, "Failed to fetch permissions")
		return
	}

	helpers.OK(c, "Permissions fetched successfully", dtos)
}

// PermissionGetByID handles GET /api/permissions/:id
func (h *Handlers) PermissionGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	dto, err := h.svcs.PermissionGetByID(c.Request.Context(), uint(id))
	if err != nil {
		helpers.NotFound(c, "Permission not found")
		return
	}

	helpers.OK(c, "Permission fetched successfully", dto)
}

// PermissionUpdate handles PUT /api/permissions/:id
func (h *Handlers) PermissionUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	var req dtos.PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	dto, err := h.svcs.PermissionUpdate(c.Request.Context(), uint(id), req)
	if err != nil {
		helpers.NotFound(c, "Permission not found")
		return
	}

	helpers.OK(c, "Permission updated successfully", dto)
}

// PermissionDelete handles DELETE /api/permissions/:id
func (h *Handlers) PermissionDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid permission ID")
		return
	}

	err = h.svcs.PermissionDelete(c.Request.Context(), uint(id))
	if err != nil {
		helpers.NotFound(c, "Permission not found")
		return
	}

	helpers.OK(c, "Permission deleted successfully", nil)
}
