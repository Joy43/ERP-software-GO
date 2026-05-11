package uploads

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// RegisterRoutes registers all upload routes
func RegisterRoutes(
	v1 *gin.RouterGroup,
	db *gorm.DB,
	rdb *redis.Client,
	h *Handler,
	authMiddleware gin.HandlerFunc,
) {
	uploadsGroup := v1.Group("/uploads", authMiddleware)

	// ==================== FOLDER ROUTES ====================
	folders := uploadsGroup.Group("/folders")
	{

		folders.POST("", h.CreateFolder)
	
		folders.GET("", h.GetRootFolders)
	
		// Get folder contents - supports both /folders/contents and /folders/:folder_id/contents
		folders.GET("/contents", h.GetFolderContents)
		folders.GET("/:folder_id/contents", h.GetFolderContents)
		folders.PUT("/:folder_id", h.UpdateFolder)
	
		folders.DELETE("/:folder_id", h.DeleteFolder)

		folders.POST("/:folder_id/restore", h.RestoreFolder)
	}

	// ==================== FILE ROUTES ====================
	files := uploadsGroup.Group("/files")
	{
	
		files.POST("", h.UploadFile)
		// Upload multiple files
		files.POST("/batch", h.UploadMultiple)
		// Get user files (with optional ?q= search)
		files.GET("", h.GetUserFiles)

	

		// Move file to folder
		files.POST("/:file_id/move", h.MoveFile)
		// Delete file
		files.DELETE("/:file_id", h.DeleteFile)
		// Restore file
		files.POST("/:file_id/restore", h.RestoreFile)
		// Download file
		files.GET("/:file_id/download", h.DownloadFile)
	
		
	}

}