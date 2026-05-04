package uploads

import (
	"time"
)

// Folder represents a directory in the file system
type Folder struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"column:name;index" json:"name"`
	ParentID    *uint64        `gorm:"column:parent_id;index" json:"parent_id,omitempty"`
	CreatedBy   uint64         `gorm:"column:created_by;index" json:"created_by"`
	Level       int            `gorm:"column:level" json:"level"`
	IsDeleted   bool           `gorm:"column:is_deleted;index:idx_parent_deleted" json:"is_deleted"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	
	//----------  Relations -------------------
	Parent   *Folder `gorm:"foreignKey:ParentID;references:ID" json:"parent,omitempty"`
	Files    []File  `gorm:"foreignKey:FolderID;references:ID" json:"files,omitempty"`
	Children []Folder `gorm:"foreignKey:ParentID;references:ID" json:"children,omitempty"`
}

// -----------  File represents a file in the system --------------
type File struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	FolderID     *uint64        `gorm:"column:folder_id;index" json:"folder_id,omitempty"`
	Name         string         `gorm:"column:name;index" json:"name"`
	OriginalName string         `gorm:"column:original_name" json:"original_name"`
	Extension    string         `gorm:"column:extension" json:"extension,omitempty"`
	MimeType     string         `gorm:"column:mime_type" json:"mime_type,omitempty"`
	Size         uint64         `gorm:"column:size" json:"size"`
	StoragePath  string         `gorm:"column:storage_path" json:"storage_path"`
	UploadedBy   uint64         `gorm:"column:uploaded_by;index" json:"uploaded_by"`
	FileType     FileTypeEnum   `gorm:"column:file_type;type:enum('image','document','video','audio','archive','zip','other')" json:"file_type"`
	IsPublic     bool           `gorm:"column:is_public" json:"is_public"`
	IsDeleted    bool           `gorm:"column:is_deleted;index" json:"is_deleted"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	
	//----------  Relations -------------------
	Folder *Folder `gorm:"foreignKey:FolderID;references:ID" json:"folder,omitempty"`
}


// StorageUsage tracks user storage utilization
type StorageUsage struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	UserID    uint64    `gorm:"column:user_id;uniqueIndex" json:"user_id"`
	TotalFiles uint64   `gorm:"column:total_files" json:"total_files"`
	TotalBytes uint64   `gorm:"column:total_bytes" json:"total_bytes"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// FileHistory represents audit trail for file operations
type FileHistory struct {
	ID          uint64          `gorm:"primaryKey" json:"id"`
	FileID      uint64          `gorm:"column:file_id;index" json:"file_id"`
	Action      FileActionEnum  `gorm:"column:action;type:enum('upload','move','rename','delete','restore','share')" json:"action"`
	PerformedBy uint64          `gorm:"column:performed_by;index" json:"performed_by"`
	Details     string          `gorm:"column:details;type:json" json:"details,omitempty"`
	CreatedAt   time.Time       `gorm:"column:created_at;index" json:"created_at"`
	
	// Relations
	File *File `gorm:"foreignKey:FileID;references:ID" json:"file,omitempty"`
}

// Enums
type FileTypeEnum string
type SharePermEnum string
type FileActionEnum string

const (
	FileTypeImage    FileTypeEnum = "image"
	FileTypeDocument FileTypeEnum = "document"
	FileTypeVideo    FileTypeEnum = "video"
	FileTypeAudio    FileTypeEnum = "audio"
	FileTypeArchive  FileTypeEnum = "archive"
	FileTypeOther    FileTypeEnum = "other"
)

const (
	SharePermView     SharePermEnum = "view"
	SharePermDownload SharePermEnum = "download"
	SharePermEdit     SharePermEnum = "edit"
)

const (
	FileActionUpload  FileActionEnum = "upload"
	FileActionMove    FileActionEnum = "move"
	FileActionRename  FileActionEnum = "rename"
	FileActionDelete  FileActionEnum = "delete"
	FileActionRestore FileActionEnum = "restore"
	FileActionShare   FileActionEnum = "share"
)

// Table names
func (Folder) TableName() string {
	return "folders"
}

func (File) TableName() string {
	return "files"
}

func (StorageUsage) TableName() string {
	return "storage_usage"
}

func (FileHistory) TableName() string {
	return "file_history"
}