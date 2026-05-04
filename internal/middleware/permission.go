package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/user"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
	"gorm.io/gorm"
)

// CachedPermissions stores user permissions in Redis for performance
type CachedPermissions struct {
	IsSuperAdmin bool     `json:"is_super_admin"`
	Permissions  []string `json:"permissions"`
}

func Permission(db *gorm.DB, rdb *redis.Client, requiredPermission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDValue, exists := ctx.Get("user_id")
		if !exists {
			response.Error(ctx, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED", nil)
			ctx.Abort()
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			response.Error(ctx, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED", nil)
			ctx.Abort()
			return
		}

		cacheKey := fmt.Sprintf("user_permissions:%d", userID)
		var cached CachedPermissions

		// 1. Try to get permissions from Redis
		val, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			// Cache hit
			if err := json.Unmarshal([]byte(val), &cached); err == nil {
				if checkAccess(cached, requiredPermission) {
					ctx.Next()
					return
				}
				response.Error(ctx, http.StatusForbidden, "forbidden: insufficient permissions (cached)", "FORBIDDEN", nil)
				ctx.Abort()
				return
			}
		}

		// 2. Cache miss or unmarshal error -> Query database
		var u user.User
		if err := db.Preload("Roles.Permissions").First(&u, userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.Error(ctx, http.StatusUnauthorized, "user not found", "UNAUTHORIZED", nil)
			} else {
				response.Error(ctx, http.StatusInternalServerError, "internal server error", "INTERNAL_SERVER_ERROR", nil)
			}
			ctx.Abort()
			return
		}

		// Prepare data for caching
		cached.IsSuperAdmin = false
		cached.Permissions = []string{}
		for _, role := range u.Roles {
			if role.Slug == "superadmin" {
				cached.IsSuperAdmin = true
			}
			for _, perm := range role.Permissions {
				cached.Permissions = append(cached.Permissions, perm.Slug)
			}
		}

		// 3. Store in Redis with 15m TTL
		if cachedData, err := json.Marshal(cached); err == nil {
			rdb.Set(ctx, cacheKey, cachedData, 15*time.Minute)
		}

		// 4. Check access
		if checkAccess(cached, requiredPermission) {
			ctx.Next()
			return
		}

		response.Error(ctx, http.StatusForbidden, "forbidden: insufficient permissions", "FORBIDDEN", nil)
		ctx.Abort()
	}
}

func checkAccess(cached CachedPermissions, requiredPermission string) bool {
	if cached.IsSuperAdmin {
		return true
	}
	for _, p := range cached.Permissions {
		if p == requiredPermission {
			return true
		}
	}
	return false
}
