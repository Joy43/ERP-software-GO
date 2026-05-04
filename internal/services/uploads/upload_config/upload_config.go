package upload_config

import (
	"os"
	"path/filepath"
)

type UploadConfig struct {
	UploadDir        string
	MaxFileSize      int64
	MaxFiles         int
	AllowedMimeTypes map[string]bool
}

var Config = UploadConfig{
	UploadDir: getUploadDir(),
	MaxFileSize: 1000 * 1024 * 1024,
	MaxFiles: 5,
	AllowedMimeTypes: map[string]bool{
	// ---------Images -----------
    "image/jpeg": true,
    "image/png":  true,
    "image/gif":  true,
    "image/webp": true,
    
    // ----------- Documents ---------------
    "application/pdf":  true,
    "text/plain":       true,
    "application/msword": true,
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
    
    // ------------Archives ----------------
    "application/zip":        true,
    "application/x-tar":      true,
    "application/gzip":       true,
    
    // ---------- Others  -------------
    "application/json":       true,
    "text/csv":               true,
	},
}

func getUploadDir() string {
	if os.Getenv("NODE_ENV") == "production" {
		return "/tmp/media"
	}
	dir := os.Getenv("UPLOAD_DIR")
	if dir != "" {
		return dir
	}
	return filepath.Join(".", "media")
}