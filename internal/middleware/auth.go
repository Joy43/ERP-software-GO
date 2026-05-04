package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/config"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/utils"
)

func Auth(cfg config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("Authorization")
		if authorizationHeader == "" {
			response.Error(ctx, http.StatusUnauthorized, "missing authorization header", "UNAUTHORIZED", nil)
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authorizationHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(ctx, http.StatusUnauthorized, "invalid authorization format", "UNAUTHORIZED", nil)
			ctx.Abort()
			return
		}

		claims, err := utils.ParseAccessToken(parts[1], cfg.JWTAccessSecret)
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, "invalid token", "UNAUTHORIZED", nil)
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_role", claims.Role)
		ctx.Set("user_email", claims.Email)
		ctx.Next()
	}
}
