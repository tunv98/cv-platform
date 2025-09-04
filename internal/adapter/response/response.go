package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is the standard response wrapper for all API endpoints
type APIResponse struct {
	Success bool        `json:"success"`         // true if successful
	Data    interface{} `json:"data,omitempty"`  // response data (if successful)
	Error   *APIError   `json:"error,omitempty"` // error information (if failed)
}

// APIError contains error details for failed requests
type APIError struct {
	Code    string `json:"code"`    // internal or business error code (e.g., INVALID_REQUEST)
	Message string `json:"message"` // human-readable error message for client
}

// Common error codes
const (
	ErrorCodeInvalidRequest   = "INVALID_REQUEST"
	ErrorCodeNotFound         = "NOT_FOUND"
	ErrorCodeInternalError    = "INTERNAL_ERROR"
	ErrorCodeValidationFailed = "VALIDATION_FAILED"
	ErrorCodeUnauthorized     = "UNAUTHORIZED"
	ErrorCodeForbidden        = "FORBIDDEN"
)

// RespondSuccess creates a successful API response
func RespondSuccess(c *gin.Context, statusCode int, data any) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
	})
}

// RespondError creates an error API response
func RespondError(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

// RespondBadRequest creates a 400 bad request response
func RespondBadRequest(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, ErrorCodeInvalidRequest, message)
}

// RespondNotFound creates a 404 not found response
func RespondNotFound(c *gin.Context, message string) {
	RespondError(c, http.StatusNotFound, ErrorCodeNotFound, message)
}

// RespondInternalErr creates a 500 internal server error response
func RespondInternalErr(c *gin.Context, message string) {
	RespondError(c, http.StatusInternalServerError, ErrorCodeInternalError, message)
}

// RespondValidationErr creates a 400 validation error response
func RespondValidationErr(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, ErrorCodeValidationFailed, message)
}
