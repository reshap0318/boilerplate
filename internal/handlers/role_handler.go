package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/dtos"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
)

// RoleCreate handles POST /api/roles
func (h *Handlers) RoleCreate(c *gin.Context) {
	var req dtos.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	dto, err := h.svcs.RoleCreate(c.Request.Context(), req)
	if err != nil {
		helpers.InternalServerError(c, "Failed to create role")
		return
	}

	helpers.Created(c, "Role created successfully", dto)
}

// RoleGetAll handles GET /api/roles
func (h *Handlers) RoleGetAll(c *gin.Context) {
	dtos, err := h.svcs.RoleGetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, "Failed to fetch roles")
		return
	}

	helpers.OK(c, "Roles fetched successfully", dtos)
}

// RoleGetByID handles GET /api/roles/:id
func (h *Handlers) RoleGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid role ID")
		return
	}

	dto, err := h.svcs.RoleGetByID(c.Request.Context(), uint(id))
	if err != nil {
		helpers.NotFound(c, "Role not found")
		return
	}

	helpers.OK(c, "Role fetched successfully", dto)
}

// RoleUpdate handles PUT /api/roles/:id
func (h *Handlers) RoleUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid role ID")
		return
	}

	var req dtos.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	dto, err := h.svcs.RoleUpdate(c.Request.Context(), uint(id), req)
	if err != nil {
		helpers.NotFound(c, "Role not found")
		return
	}

	helpers.OK(c, "Role updated successfully", dto)
}

// RoleDelete handles DELETE /api/roles/:id
func (h *Handlers) RoleDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid role ID")
		return
	}

	err = h.svcs.RoleDelete(c.Request.Context(), uint(id))
	if err != nil {
		helpers.NotFound(c, "Role not found")
		return
	}

	helpers.OK(c, "Role deleted successfully", nil)
}

// RoleGetPermissions handles GET /api/roles/:id/permissions
func (h *Handlers) RoleGetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid role ID")
		return
	}

	perms, err := h.svcs.RoleGetPermissions(c.Request.Context(), uint(id))
	if err != nil {
		helpers.InternalServerError(c, "Failed to fetch role permissions")
		return
	}

	helpers.OK(c, "Role permissions fetched successfully", perms)
}
