package filter

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	ErrorCode string      `json:"error_code,omitempty"`
	Details   interface{} `json:"details,omitempty"`
}


type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

//------------ HTTPException represents an HTTP exception----------
type HTTPException struct {
	Status    int
	Message   string
	ErrorCode string
	Details   interface{}
}

//----------- Error implements the error interface-------------
func (e *HTTPException) Error() string {
	return e.Message
}

//------------ NewHTTPException creates a new HTTP exception---------------
func NewHTTPException(status int, message string, errorCode string, details interface{}) *HTTPException {
	return &HTTPException{
		Status:    status,
		Message:   message,
		ErrorCode: errorCode,
		Details:   details,
	}
}

//------------ ExceptionHandler is the main exception handling middleware-------------
func ExceptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				handlePanic(c, r)
			}
		}()

		c.Next()

		// ------------Handle errors that were added to context------------
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			handleError(c, err)
		}
	}
}

// ------------ handleError handles different types of errors--------------
func handleError(c *gin.Context, err error) {
	var httpException *HTTPException
	var validationErrors []ValidationError
	var status int = http.StatusInternalServerError
	var message string = "Internal server error"
	var errorCode string = "INTERNAL_SERVER_ERROR"
	var details interface{}

	// --------------- Check if it's our custom HTTPException-------------
	if errors.As(err, &httpException) {
		status = httpException.Status
		message = httpException.Message
		errorCode = httpException.ErrorCode
		details = httpException.Details
	} else if err != nil {
		// ---------------Check for validation errors---------------
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			status = http.StatusBadRequest
			message = "Validation failed"
			errorCode = "VALIDATION_ERROR"
			validationErrors = parseValidationErrors(ve)
			details = validationErrors
		} else {
			//---------  Generic error ---------------
			message = err.Error()
			errorCode = "ERROR"
		}
	}

	response := ErrorResponse{
		Success:   false,
		Message:   message,
		ErrorCode: errorCode,
		Details:   details,
	}

	c.JSON(status, response)
}

// handlePanic handles panic recovery
func handlePanic(c *gin.Context, r interface{}) {
	status := http.StatusInternalServerError
	message := "Internal server error"
	errorCode := "PANIC_ERROR"
	var details interface{}

	if err, ok := r.(error); ok {
		details = err.Error()
	} else if str, ok := r.(string); ok {
		details = str
	} else {
		details = r
	}

	response := ErrorResponse{
		Success:   false,
		Message:   message,
		ErrorCode: errorCode,
		Details:   details,
	}

	c.JSON(status, response)
}

// -----------parseValidationErrors converts validator errors to ValidationError slice------------
func parseValidationErrors(ve validator.ValidationErrors) []ValidationError {
	var errors []ValidationError

	for _, err := range ve {
		validationErr := ValidationError{
			Field:   err.Field(),
			Message: getValidationMessage(err),
		}
		errors = append(errors, validationErr)
	}

	return errors
}

//------------ getValidationMessage returns a human-readable validation message--------------
func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters"
	case "max":
		return err.Field() + " must not exceed " + err.Param() + " characters"
	case "len":
		return err.Field() + " must be exactly " + err.Param() + " characters"
	case "numeric":
		return err.Field() + " must be numeric"
	case "url":
		return err.Field() + " must be a valid URL"
	case "gte":
		return err.Field() + " must be greater than or equal to " + err.Param()
	case "lte":
		return err.Field() + " must be less than or equal to " + err.Param()
	case "gt":
		return err.Field() + " must be greater than " + err.Param()
	case "lt":
		return err.Field() + " must be less than " + err.Param()
	default:
		return err.Field() + " validation failed for " + err.Tag()
	}
}

// Custom error constructors for common HTTP statuses

// BadRequestError creates a 400 Bad Request error
func BadRequestError(message string, details interface{}) error {
	exc := NewHTTPException(http.StatusBadRequest, message, "BAD_REQUEST", details)
	return exc
}

//-------------- UnauthorizedError creates a 401 Unauthorized error---------
func UnauthorizedError(message string, details interface{}) error {
	exc := NewHTTPException(http.StatusUnauthorized, message, "UNAUTHORIZED", details)
	return exc
}

// ---------------ForbiddenError creates a 403 Forbidden error----------------
func ForbiddenError(message string, details interface{}) error {
	exc := NewHTTPException(http.StatusForbidden, message, "FORBIDDEN", details)
	return exc
}

// ----------------- NotFoundError creates a 404 Not Found error ---------------------
func NotFoundError(message string, details interface{}) error {
	exc := NewHTTPException(http.StatusNotFound, message, "NOT_FOUND", details)
	return exc
}

//---------------- ConflictError creates a 409 Conflict error --------------
func ConflictError(message string, details interface{}) error {
	exc := NewHTTPException(http.StatusConflict, message, "CONFLICT", details)
	return exc
}

// InternalServerError creates a 500 Internal Server Error
func InternalServerError(message string, details interface{}) error {
	exc := NewHTTPException(http.StatusInternalServerError, message, "INTERNAL_SERVER_ERROR", details)
	return exc
}

// ------------ServiceUnavailableError creates a 503 Service Unavailable error--------------
func ServiceUnavailableError(message string, details interface{}) error {
	exc := NewHTTPException(http.StatusServiceUnavailable, message, "SERVICE_UNAVAILABLE", details)
	return exc
}
