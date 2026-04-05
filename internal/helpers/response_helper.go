package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard HTTP response structure
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// SuccessResponse sends a success response with optional data
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	resp := Response{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, resp)
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	resp := Response{
		Code:    statusCode,
		Message: message,
	}
	c.JSON(statusCode, resp)
}

// OK sends 200 OK response
func OK(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusOK, message, data)
}

// Created sends 201 Created response
func Created(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusCreated, message, data)
}

// BadRequest sends 400 Bad Request response
func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// Unauthorized sends 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message)
}

// Forbidden sends 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message)
}

// NotFound sends 404 Not Found response
func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

// InternalServerError sends 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message)
}
