package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

func RBAC(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(allowedRoles) == 0 {
			ctx.Next()
			return
		}

		roleValue, exists := ctx.Get("user_role")
		if !exists {
			response.Error(ctx, http.StatusForbidden, "forbidden", "FORBIDDEN", nil)
			ctx.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			response.Error(ctx, http.StatusForbidden, "forbidden", "FORBIDDEN", nil)
			ctx.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				ctx.Next()
				return
			}
		}

		response.Error(ctx, http.StatusForbidden, "forbidden", "FORBIDDEN", nil)
		ctx.Abort()
	}
}
