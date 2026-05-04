package audit

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(
	v1 *gin.RouterGroup,
	db *gorm.DB,
	rdb *redis.Client,
	h *Handlers,
	authMiddleware gin.HandlerFunc,
) {
	auditGroup := v1.Group("/audit", authMiddleware)
	{
		// Transaction History
		history := auditGroup.Group("/transaction-history")
		{
			history.GET("", middleware.Permission(db, rdb, "transaction_history.view"), h.TransactionHistory.List)
			history.POST("", h.TransactionHistory.Create) // Internal logging mostly
			history.DELETE("/:id", middleware.Permission(db, rdb, "transaction_history.delete"), h.TransactionHistory.Delete)
		}
	}
}
