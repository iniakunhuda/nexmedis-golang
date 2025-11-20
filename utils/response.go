package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponseWithMessage sends an error response with a custom message
func ErrorResponseWithMessage(c echo.Context, statusCode int, message string, details string) error {
	return c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error:   message,
		Details: details,
	})
}

// BadRequestResponse sends a 400 Bad Request response
func BadRequestResponse(c echo.Context, message string) error {
	return ErrorResponseWithMessage(c, http.StatusBadRequest, message, "")
}

// UnauthorizedResponse sends a 401 Unauthorized response
func UnauthorizedResponse(c echo.Context, message string) error {
	return ErrorResponseWithMessage(c, http.StatusUnauthorized, message, "")
}

// ForbiddenResponse sends a 403 Forbidden response
func ForbiddenResponse(c echo.Context, message string) error {
	return ErrorResponseWithMessage(c, http.StatusForbidden, message, "")
}

// NotFoundResponse sends a 404 Not Found response
func NotFoundResponse(c echo.Context, message string) error {
	return ErrorResponseWithMessage(c, http.StatusNotFound, message, "")
}

// ConflictResponse sends a 409 Conflict response
func ConflictResponse(c echo.Context, message string) error {
	return ErrorResponseWithMessage(c, http.StatusConflict, message, "")
}

// InternalServerErrorResponse sends a 500 Internal Server Error response
func InternalServerErrorResponse(c echo.Context, message string, details string) error {
	return ErrorResponseWithMessage(c, http.StatusInternalServerError, message, details)
}

// TooManyRequestsResponse sends a 429 Too Many Requests response
func TooManyRequestsResponse(c echo.Context, message string) error {
	return ErrorResponseWithMessage(c, http.StatusTooManyRequests, message, "")
}

// CreatedResponse sends a 201 Created response
func CreatedResponse(c echo.Context, message string, data interface{}) error {
	return SuccessResponse(c, http.StatusCreated, message, data)
}

// OKResponse sends a 200 OK response
func OKResponse(c echo.Context, message string, data interface{}) error {
	return SuccessResponse(c, http.StatusOK, message, data)
}

func ServiceUnavailableResponse(c echo.Context, message string) error {
	return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
		"success": false,
		"message": message,
		"error":   "Service temporarily unavailable",
	})
}
