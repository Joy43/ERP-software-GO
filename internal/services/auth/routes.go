package auth

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
	authGroup := v1.Group("/auth")
	authGroup.POST("/login", h.User.Login)
	authGroup.POST("/refresh", h.User.RefreshToken)

	//-------------- User Management ----------------
	userManagement := authGroup.Group("/users", authMiddleware)
	{
		userManagement.GET("/me", h.User.Me)
		userManagement.GET("", middleware.Permission(db, rdb, "users.view"), h.User.List)
		userManagement.GET("/:id", middleware.Permission(db, rdb, "users.view"), h.User.Get)
		userManagement.POST("/create-user", middleware.Permission(db, rdb, "users.create"), h.User.CreateUser)
		userManagement.PUT("/:id", middleware.Permission(db, rdb, "users.edit"), h.User.UpdateUser)
		userManagement.DELETE("/:id", middleware.Permission(db, rdb, "users.delete"), h.User.DeleteUser)
	}

	

	//------------- Role routes------------
	roles := authGroup.Group("/roles", authMiddleware)
	{
		roles.POST("", middleware.Permission(db, rdb, "roles.create"), h.Role.Create)
		roles.GET("", middleware.Permission(db, rdb, "roles.view"), h.Role.List)
		roles.GET("/:id", middleware.Permission(db, rdb, "roles.view"), h.Role.Get)
		roles.PUT("/:id", middleware.Permission(db, rdb, "roles.edit"), h.Role.Update)
		roles.DELETE("/:id", middleware.Permission(db, rdb, "roles.delete"), h.Role.Delete)
	}

	// Permission routes
	permissions := authGroup.Group("/permissions", authMiddleware)
	{
		permissions.GET("", middleware.Permission(db, rdb, "permissions.view"), h.Permission.List)
	}

	// Designation routes
	designations := authGroup.Group("/designations", authMiddleware)
	{
		designations.POST("", middleware.Permission(db, rdb, "designations.create"), h.Designation.Create)
		designations.GET("", middleware.Permission(db, rdb, "designations.view"), h.Designation.List)
		designations.GET("/:id", middleware.Permission(db, rdb, "designations.view"), h.Designation.Get)
		designations.PUT("/:id", middleware.Permission(db, rdb, "designations.edit"), h.Designation.Update)
		designations.DELETE("/:id", middleware.Permission(db, rdb, "designations.delete"), h.Designation.Delete)
	}

	// Department routes
	departments := authGroup.Group("/departments", authMiddleware)
	{
		departments.POST("", middleware.Permission(db, rdb, "departments.create"), h.Department.Create)
		departments.GET("", middleware.Permission(db, rdb, "departments.view"), h.Department.List)
		departments.GET("/:id", middleware.Permission(db, rdb, "departments.view"), h.Department.Get)
		departments.PUT("/:id", middleware.Permission(db, rdb, "departments.edit"), h.Department.Update)
		departments.DELETE("/:id", middleware.Permission(db, rdb, "departments.delete"), h.Department.Delete)
	}

	// ----------------Office routes-----------------
	offices := authGroup.Group("/offices", authMiddleware)
	{
		offices.POST("", middleware.Permission(db, rdb, "offices.create"), h.Office.Create)
		offices.GET("", middleware.Permission(db, rdb, "offices.view"), h.Office.List)
		offices.GET("/:id", middleware.Permission(db, rdb, "offices.view"), h.Office.Get)
		offices.PUT("/:id", middleware.Permission(db, rdb, "offices.edit"), h.Office.Update)
		offices.DELETE("/:id", middleware.Permission(db, rdb, "offices.delete"), h.Office.Delete)
	}

	// ---------- Payment Mode routes (GET all only) -----------
	paymentModes := authGroup.Group("/payment-modes", authMiddleware)
	{
		paymentModes.GET("", h.PaymentMode.List)
	}

	//--------- Responsibility Transfer routes. ------------	
	transfers := authGroup.Group("/responsibility-transfers", authMiddleware)
	{
		transfers.GET("", middleware.Permission(db, rdb, "responsibility_transfers.view"), h.ResponsibilityTransfer.List)
		transfers.GET("/:id", middleware.Permission(db, rdb, "responsibility_transfers.view"), h.ResponsibilityTransfer.Get)
		transfers.POST("", middleware.Permission(db, rdb, "responsibility_transfers.create"), h.ResponsibilityTransfer.Create)
		transfers.POST("/:id/approve", middleware.Permission(db, rdb, "responsibility_transfers.approve"), h.ResponsibilityTransfer.Approve)
		transfers.DELETE("/:id", middleware.Permission(db, rdb, "responsibility_transfers.delete"), h.ResponsibilityTransfer.Delete)
	}

	// ------------Wallet routes-------------
	wallets := authGroup.Group("/wallets", authMiddleware)
	{
		wallets.POST("", middleware.Permission(db, rdb, "wallets.create"), h.Wallet.Create)
		wallets.GET("", middleware.Permission(db, rdb, "wallets.view"), h.Wallet.List)
		wallets.GET("/:id", middleware.Permission(db, rdb, "wallets.view"), h.Wallet.Get)
		wallets.PUT("/:id", middleware.Permission(db, rdb, "wallets.edit"), h.Wallet.Update)
		wallets.DELETE("/:id", middleware.Permission(db, rdb, "wallets.delete"), h.Wallet.Delete)
	}
}
