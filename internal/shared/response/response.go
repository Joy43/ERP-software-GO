package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
	Details   interface{} `json:"details,omitempty"`
}

func Success(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(200, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(ctx *gin.Context, status int, message, errorCode string, details interface{}) {
	ctx.JSON(status, APIResponse{
		Success:   false,
		Message:   message,
		ErrorCode: errorCode,
		Details:   details,
	})
}
