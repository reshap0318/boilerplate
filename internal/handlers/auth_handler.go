package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/go-boilerplate/internal/dtos"
	"github.com/reshap0318/go-boilerplate/internal/helpers"
)

// AuthLogin handles user login.
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.LoginRequest true "Login credentials"
// @Success 200 {object} dtos.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func (h *Handlers) AuthLogin(c *gin.Context) {
	var req dtos.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := h.svcs.AuthLogin(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// AuthRefreshToken handles token refresh.
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dtos.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/refresh [post]
func (h *Handlers) AuthRefreshToken(c *gin.Context) {
	var req dtos.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := h.svcs.AuthRefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// AuthLogout handles user logout.
// @Summary User logout
// @Description Logout user (client should clear token)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /api/auth/logout [post]
func (h *Handlers) AuthLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful. Please clear your token on the client side.",
	})
}

// AuthForgetPassword handles forget password request.
// @Summary Forget password
// @Description Send reset password email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.ForgetPasswordRequest true "Email and frontend URL"
// @Success 200 {object} dtos.MessageResponse
// @Router /api/auth/forgot-password [post]
func (h *Handlers) AuthForgetPassword(c *gin.Context) {
	var req dtos.ForgetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.svcs.AuthForgetPassword(c.Request.Context(), req.Email, req.FrontendURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dtos.MessageResponse{
		Message: "If the email exists, a reset link has been sent",
	})
}

// AuthResetPassword handles reset password request.
// @Summary Reset password
// @Description Reset password with token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.ResetPasswordRequest true "Token and new password"
// @Success 200 {object} dtos.MessageResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/reset-password [post]
func (h *Handlers) AuthResetPassword(c *gin.Context) {
	var req dtos.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.svcs.AuthResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		if err == helpers.ErrTokenInvalid || err == helpers.ErrTokenExpired || err == helpers.ErrTokenUsed {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err == helpers.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dtos.MessageResponse{
		Message: "Password has been reset successfully",
	})
}
