package uploads

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

//-------- PaginatedResponse represents a paginated response---------
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total,omitempty"`
	Page    int         `json:"page,omitempty"`
	Limit   int         `json:"limit,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// -----------Success sends a success response ----------------
func Success(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Code:    statusCode,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string, errorCode string, details interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   errorCode,
		Data:    details,
		Code:    statusCode,
	})
}

// SuccessList sends a paginated success response
func SuccessList(c *gin.Context, statusCode int, data interface{}, total int64, page int, limit int, message string) {
	c.JSON(statusCode, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}

//-----------  ErrorList sends a paginated error response ------------
func ErrorList(c *gin.Context, statusCode int, message string, error string) {
	c.JSON(statusCode, PaginatedResponse{
		Success: false,
		Message: message,
		Error:   error,
	})
}

// NotFound sends a 404 response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Message: message,
		Code:    http.StatusNotFound,
	})
}

//--------- BadRequest sends a 400 response---------------
func BadRequest(c *gin.Context, message string, errorCode string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: message,
		Error:   errorCode,
		Code:    http.StatusBadRequest,
	})
}

// Unauthorized sends a 401 response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Message: message,
		Code:    http.StatusUnauthorized,
	})
}

//--------  Forbidden sends a 403 response ------------
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, APIResponse{
		Success: false,
		Message: message,
		Code:    http.StatusForbidden,
	})
}

// InternalServerError sends a 500 response
func InternalServerError(c *gin.Context, message string, errorCode string) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Message: message,
		Error:   errorCode,
		Code:    http.StatusInternalServerError,
	})
}
