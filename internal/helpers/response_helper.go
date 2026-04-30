package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

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

// ValidationError sends 422 Unprocessable Entity response for validation errors, or 400 for JSON syntax errors
func ValidationError(c *gin.Context, err error) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		resp := Response{
			Code:    http.StatusUnprocessableEntity,
			Message: "The given data was invalid.",
			Errors:  formatValidationError(validationErrs),
		}
		c.JSON(http.StatusUnprocessableEntity, resp)
		return
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		field := unmarshalTypeError.Field
		if field == "" {
			field = "payload"
		} else {
			// Convert to lowercase if needed, although Field from unmarshalTypeError is often derived from the json tag directly.
			field = strings.ToLower(field)
		}
		
		errorsMap := map[string][]string{
			field: {fmt.Sprintf("The %s field must be of type %s.", field, unmarshalTypeError.Type.String())},
		}

		resp := Response{
			Code:    http.StatusUnprocessableEntity,
			Message: "The given data was invalid.",
			Errors:  errorsMap,
		}
		c.JSON(http.StatusUnprocessableEntity, resp)
		return
	}

	// Jika bukan error validasi (misalnya JSON syntax error atau tipe error lainnya)
	resp := Response{
		Code:    http.StatusBadRequest,
		Message: "Invalid JSON payload: " + err.Error(),
	}
	c.JSON(http.StatusBadRequest, resp)
}

func formatValidationError(validationErrs validator.ValidationErrors) map[string][]string {
	errorsMap := make(map[string][]string)

	for _, e := range validationErrs {
		field := strings.ToLower(e.Field())
		errorsMap[field] = append(errorsMap[field], getErrorMessage(e))
	}

	return errorsMap
}

func getErrorMessage(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	switch e.Tag() {
	case "required":
		return "The " + field + " field is required."
	case "email":
		return "The " + field + " must be a valid email address."
	case "min":
		return "The " + field + " must be at least " + e.Param() + " characters."
	case "max":
		return "The " + field + " may not be greater than " + e.Param() + " characters."
	default:
		return "The " + field + " field is invalid."
	}
}
