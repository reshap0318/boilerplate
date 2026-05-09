package handlers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/dtos"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
	"github.com/reshap0318/go-boilerplate/internal/repositories"
)

// UserCreate handles POST /api/users
func (h *Handlers) UserCreate(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		helpers.BadRequest(c, "Failed to parse form data")
		return
	}

	var req dtos.UserRequest
	req.Name = c.Request.FormValue("name")
	req.Email = c.Request.FormValue("email")
	req.Password = c.Request.FormValue("password")
	req.PasswordConfirmation = c.Request.FormValue("password_confirmation")

	if roles := c.Request.FormValue("roles"); roles != "" {
		// Parse roles from comma-separated string
		roleStrs := strings.Split(roles, ",")
		for _, r := range roleStrs {
			if id, err := strconv.ParseUint(strings.TrimSpace(r), 10, 64); err == nil {
				req.Roles = append(req.Roles, uint(id))
			}
		}
	}

	if err := helpers.ValidateStruct(req); err != nil {
		helpers.ValidationError(c, err)
		return
	}

	avatarPath := ""
	if _, _, err := c.Request.FormFile("avatar"); err == nil {
		var saveErr error
		avatarPath, saveErr = helpers.SaveUploadedFile(c, "avatar", "storage/avatars")
		if saveErr != nil {
			helpers.BadRequest(c, saveErr.Error())
			return
		}
	}

	dto, err := h.svcs.UserCreate(c.Request.Context(), req, avatarPath)
	if err != nil {
		if avatarPath != "" {
			helpers.DeleteFile(avatarPath)
		}
		if err == helpers.ErrUserExists {
			helpers.BadRequest(c, "Email already exists")
			return
		}
		helpers.InternalServerError(c, "Failed to create user")
		return
	}

	helpers.Created(c, "User created successfully", dto)
}

// UserGetAll handles GET /api/users with optional pagination
func (h *Handlers) UserGetAll(c *gin.Context) {
	pageStr := c.Query("page")

	if pageStr == "" {
		// No pagination - return all
		users, err := h.svcs.UserGetAll(c.Request.Context())
		if err != nil {
			helpers.InternalServerError(c, "Failed to fetch users")
			return
		}

		helpers.OK(c, "Users fetched successfully", users)
		return
	}

	// Pagination requested
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	opts := &repositories.QueryOptions{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.svcs.UserGetAllPaginated(c.Request.Context(), opts)
	if err != nil {
		helpers.InternalServerError(c, "Failed to fetch users")
		return
	}

	helpers.OKWithMetadata(c, "Users fetched successfully", result)
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

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		helpers.BadRequest(c, "Failed to parse form data")
		return
	}

	var req dtos.UserRequest
	req.Name = c.Request.FormValue("name")
	req.Email = c.Request.FormValue("email")
	req.Password = c.Request.FormValue("password")
	req.PasswordConfirmation = c.Request.FormValue("password_confirmation")

	if roles := c.Request.FormValue("roles"); roles != "" {
		roleStrs := strings.Split(roles, ",")
		for _, r := range roleStrs {
			if id, err := strconv.ParseUint(strings.TrimSpace(r), 10, 64); err == nil {
				req.Roles = append(req.Roles, uint(id))
			}
		}
	}

	if req.Password == "" {
		req.PasswordConfirmation = ""
	}

	if err := helpers.ValidateStruct(req); err != nil {
		helpers.ValidationError(c, err)
		return
	}

	avatarPath := ""
	if _, _, err := c.Request.FormFile("avatar"); err == nil {
		var saveErr error
		avatarPath, saveErr = helpers.SaveUploadedFile(c, "avatar", "storage/avatars")
		if saveErr != nil {
			helpers.BadRequest(c, saveErr.Error())
			return
		}
	}

	dto, oldAvatar, err := h.svcs.UserUpdate(c.Request.Context(), uint(id), req, avatarPath)
	if err != nil {
		if avatarPath != "" {
			helpers.DeleteFile(avatarPath)
		}
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

	if avatarPath != "" && oldAvatar != "" {
		helpers.DeleteFile(oldAvatar)
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
