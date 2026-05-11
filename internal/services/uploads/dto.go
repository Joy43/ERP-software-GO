package uploads

import "time"

// ==================== FOLDER DTOs ====================

//------- CreateFolderRequest represents folder creation request---------
type CreateFolderRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	ParentID    *uint64 `json:"parent_id"`
}

//--------- FolderResponse represents folder response-----------
type FolderResponse struct {
	ID          uint64            `json:"id"`
	Name        string            `json:"name"`
	ParentID    *uint64           `json:"parent_id,omitempty"`
    StoragePath  string    `json:"storage_path,omitempty"`
	Level       int               `json:"level"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Files       int64             `json:"file_count,omitempty"`
}

// FolderResponseWithChildren represents folder response with nested children
type FolderResponseWithChildren struct {
	ID          uint64                        `json:"id"`
	Name        string                        `json:"name"`
	ParentID    *uint64                       `json:"parent_id"`
	StoragePath string                        `json:"storage_path,omitempty"`
	Level       int                           `json:"level"`
	CreatedAt   time.Time                     `json:"created_at"`
	UpdatedAt   time.Time                     `json:"updated_at"`
	Files       int64                         `json:"file_count,omitempty"`
	Children    []FolderResponseWithChildren  `json:"children,omitempty"`
}

//---------- UpdateFolderRequest represents folder update request---------
type UpdateFolderRequest struct {
	Name        string `json:"name"`

}

//---------- FolderContentsResponse represents folder contents-----------
type FolderContentsResponse struct {
	Folders      []FolderResponse `json:"folders"`
	Files        []FileResponse   `json:"files"`
	FolderTotal  int64            `json:"folder_total"`
	FileTotal    int64            `json:"file_total"`
	Page         int              `json:"page"`
	Limit        int              `json:"limit"`
}

// FolderPathResponse represents breadcrumb path
type FolderPathResponse struct {
	Path []FolderResponse `json:"path"`
}

// ==================== FILE DTOs ====================

// FileUploadRequest represents file upload metadata
type FileUploadRequest struct {
	FolderID *uint64 `form:"folder_id"`
}

// FileResponse represents file response
type FileResponse struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	OriginalName string    `json:"original_name"`
	Extension    string    `json:"extension"`
	MimeType     string    `json:"mime_type,omitempty"`
	Size         uint64    `json:"size"`
	FileType     string    `json:"file_type"`
	IsPublic     bool      `json:"is_public"`
	StoragePath  string    `json:"storage_path,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	FolderName   string    `json:"folder_name,omitempty"`
}

// FileDetailResponse represents detailed file information
type FileDetailResponse struct {
	ID           uint64          `json:"id"`
	Name         string          `json:"name"`
	OriginalName string          `json:"original_name"`
	Extension    string          `json:"extension"`
	MimeType     string          `json:"mime_type"`
	Size         uint64          `json:"size"`
	FileType     string          `json:"file_type"`
	IsPublic     bool            `json:"is_public"`
	StoragePath  string          `json:"storage_path"`
	UploadedBy   uint64          `json:"uploaded_by"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Folder       *FolderResponse `json:"folder,omitempty"`
}

// FileSearchResponse represents search result
type FileSearchResponse struct {
	Files []FileResponse `json:"files"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// UpdateFileMetadataRequest represents metadata update request
type UpdateFileMetadataRequest struct {
	Name string `json:"name"`
}

// MoveFileRequest represents move file request
type MoveFileRequest struct {
	FolderID *uint64 `json:"folder_id" binding:"required"`
}

// ==================== STORAGE DTOs ====================

// StorageUsageResponse represents storage usage information
type StorageUsageResponse struct {
	TotalFiles uint64  `json:"total_files"`
	TotalBytes uint64  `json:"total_bytes"`
	UsedPercent float64 `json:"used_percent"`
	MaxBytes   uint64  `json:"max_bytes"`
}

// ==================== HISTORY DTOs ====================

type FileHistoryResponse struct {
	ID          uint64    `json:"id"`
	FileID      uint64    `json:"file_id"`
	Action      string    `json:"action"`
	PerformedBy uint64    `json:"performed_by"`
	Details     string    `json:"details,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// ==================== PAGINATION DTOs ====================

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
}

// ==================== GENERIC RESPONSE DTOs ====================

// UploadResponse represents upload response
type UploadResponse struct {
	Files []FileResponse `json:"files"`
	Failed []string        `json:"failed,omitempty"`
}

// MessageResponse represents simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// ==================== HANDLER RESPONSE WRAPPERS ====================

// SuccessResponse wraps successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ListResponse wraps list API response
type ListResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total,omitempty"`
	Page    int         `json:"page,omitempty"`
	Limit   int         `json:"limit,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ResponseMeta contains metadata for responses
type ResponseMeta struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}