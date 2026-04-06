package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/dtos"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/repositories"
)

// UserCreate handles POST /api/users
func (h *Handlers) UserCreate(c *gin.Context) {
	var req dtos.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	dto, err := h.svcs.UserCreate(c.Request.Context(), req)
	if err != nil {
		if err == helpers.ErrUserExists {
			helpers.BadRequest(c, "Email already exists")
			return
		}
		helpers.InternalServerError(c, "Failed to create user")
		return
	}

	helpers.Created(c, "User created successfully", dto)
}

// UserGetAll handles GET /api/users with pagination
func (h *Handlers) UserGetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.UserGetAll(c.Request.Context(), opts)
	if err != nil {
		helpers.InternalServerError(c, "Failed to fetch users")
		return
	}

	helpers.OK(c, "Users fetched successfully", result)
}

// UserGetByID handles GET /api/users/:id
func (h *Handlers) UserGetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	dto, err := h.svcs.UserGetByID(c.Request.Context(), uint(id))
	if err != nil {
		helpers.NotFound(c, "User not found")
		return
	}

	helpers.OK(c, "User fetched successfully", dto)
}

// UserUpdate handles PUT /api/users/:id
func (h *Handlers) UserUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	var req dtos.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, err.Error())
		return
	}

	dto, err := h.svcs.UserUpdate(c.Request.Context(), uint(id), req)
	if err != nil {
		if err == helpers.ErrNotFound {
			helpers.NotFound(c, "User not found")
			return
		}
		if err == helpers.ErrUserExists {
			helpers.BadRequest(c, "Email already exists")
			return
		}
		helpers.InternalServerError(c, "Failed to update user")
		return
	}

	helpers.OK(c, "User updated successfully", dto)
}

// UserDelete handles DELETE /api/users/:id
func (h *Handlers) UserDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	err = h.svcs.UserDelete(c.Request.Context(), uint(id))
	if err != nil {
		helpers.NotFound(c, "User not found")
		return
	}

	helpers.OK(c, "User deleted successfully", nil)
}
