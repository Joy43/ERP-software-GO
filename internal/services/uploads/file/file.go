package file

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"time"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/uploads/upload_config"
)

func ValidateFile(file *multipart.FileHeader) error {
	if file.Size > upload_config.Config.MaxFileSize {
		return errors.New("file too large")
	}

	contentType := file.Header.Get("Content-Type")
	if !upload_config.Config.AllowedMimeTypes[contentType] {
		return errors.New("file type not allowed")
	}

	return nil
}

func GenerateFilename(original string) string {
	ext := filepath.Ext(original)
	return  filepath.Base(
		time.Now().Format("20060102150405") + "_" + original,
	) + ext
}