package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PaginationMeta represents pagination metadata in response.
type PaginationMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponseInterface defines the contract for paginated responses.
type PaginatedResponseInterface interface {
	GetData() interface{}
	GetMetadata() PaginationMeta
}

// Response represents a standard HTTP response structure
type Response struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    interface{}         `json:"data,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
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

// OKWithMetadata sends 200 OK response with pagination metadata extracted from PaginatedResponse
func OKWithMetadata(c *gin.Context, message string, paginated PaginatedResponseInterface) {
	c.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  message,
		"data":     paginated.GetData(),
		"metadata": paginated.GetMetadata(),
	})
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

// ValidationErrorWithMap sends 422 Unprocessable Entity response with pre-formatted errors map.
func ValidationErrorWithMap(c *gin.Context, errorsMap map[string][]string) {
	resp := Response{
		Code:    http.StatusUnprocessableEntity,
		Message: "The given data was invalid.",
		Errors:  errorsMap,
	}
	c.JSON(http.StatusUnprocessableEntity, resp)
}

// ValidationErrorWithField sends 422 Unprocessable Entity response for a single field error.
func ValidationErrorWithField(c *gin.Context, field string, message string) {
	errorsMap := map[string][]string{
		field: {message},
	}
	ValidationErrorWithMap(c, errorsMap)
}
