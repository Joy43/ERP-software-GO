package uploads

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/uploads/upload_config"
)

// IService defines the upload service interface
type IService interface {
	CreateFolder(name string, parentID *uint64, userID uint64) (*Folder, error)
	GetRootFolders(userID uint64) ([]Folder, error)
	GetFolderContents(folderID *uint64, userID uint64, limit int, page int) ([]Folder, int64, []File, int64, error)
	GetFolderPath(folderID uint64) ([]Folder, error)
	UpdateFolder(folderID uint64, name string) (*Folder, error)
	DeleteFolder(folderID uint64) error
	RestoreFolder(folderID uint64) error

	//------------ File operations ---------------
	UploadFile(file *multipart.FileHeader, folderID *uint64, userID uint64) (*File, error)
	UploadMultiple(files []*multipart.FileHeader, folderID *uint64, userID uint64) ([]File, []string, error)
	GetFile(fileID uint64) (*File, error)
	GetUserFiles(userID uint64, limit int, page int) ([]File, int64, error)
	GetFolderFiles(folderID, userID uint64, limit, page int) ([]File, int64, error)
	SearchFiles(query string, userID uint64, limit int, page int) ([]File, int64, error)
	UpdateFileMetadata(fileID uint64, name string) (*File, error)
	DeleteFile(fileID uint64) error
	RestoreFile(fileID uint64) error
	MoveFile(fileID uint64, newFolderID *uint64) (*File, error)
	DownloadFile(fileID uint64) (string, error)

	// -----------Storage operations----------------
	GetStorageUsage(userID uint64) (*StorageUsage, error)

	// -----------History operations ----------------
	GetFileHistory(fileID uint64, limit int) ([]FileHistory, error)
}

//--------- Service implements IService---------
type Service struct {
	repo IUploadRepository
	baseURL string
}

// NewService creates a new upload service
func NewService(repo IUploadRepository, baseURL string) IService {
	return &Service{repo: repo, baseURL: baseURL}
}

// ==================== FOLDER OPERATIONS ====================

func (s *Service) CreateFolder(name string, parentID *uint64, userID uint64) (*Folder, error) {

	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("folder name cannot be empty")
	}

	// Check for duplicate name
	exists, err := s.repo.CheckDuplicateFolderName(parentID, name, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("folder with this name already exists")
	}

	level := 0
	if parentID != nil {
		parent, err := s.repo.GetFolderByID(*parentID)
		if err != nil {
			return nil, fmt.Errorf("parent folder not found")
		}
		level = parent.Level + 1
	}

	folder := &Folder{
		Name:      name,
		ParentID:  parentID,
		CreatedBy: userID,
		Level:     level,
		IsDeleted: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateFolder(folder); err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *Service) GetRootFolders(userID uint64) ([]Folder, error) {
	return s.repo.GetRootFolders(userID)
}

func (s *Service) GetFolderContents(folderID *uint64, userID uint64, limit int, page int) ([]Folder, int64, []File, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetFolderContents(folderID, userID, limit, offset)
}

func (s *Service) GetFolderPath(folderID uint64) ([]Folder, error) {
	return s.repo.GetFolderPath(folderID)
}

func (s *Service) UpdateFolder(folderID uint64, name string) (*Folder, error) {
	folder, err := s.repo.GetFolderByID(folderID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(name) != "" {
		folder.Name = name
	}
	folder.UpdatedAt = time.Now()

	if err := s.repo.UpdateFolder(folder); err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *Service) DeleteFolder(folderID uint64) error {
	return s.repo.DeleteFolder(folderID)
}

func (s *Service) RestoreFolder(folderID uint64) error {
	return s.repo.RestoreFolder(folderID)
}

// ==================== FILE OPERATIONS ====================

func (s *Service) UploadFile(file *multipart.FileHeader, folderID *uint64, userID uint64) (*File, error) {
	if file == nil {
		return nil, fmt.Errorf("file is required")
	}

	if file.Size > upload_config.Config.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum limit of %d bytes", upload_config.Config.MaxFileSize)
	}

	if folderID != nil {
		exists, err := s.repo.CheckFolderExists(*folderID)
		if err != nil || !exists {
			return nil, fmt.Errorf("folder not found")
		}
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	if !upload_config.Config.AllowedMimeTypes[mimeType] && mimeType != "application/octet-stream" {
		return nil, fmt.Errorf("file type not allowed: %s", mimeType)
	}

	//------------ Generate unique filename--------------
	filename := generateUniqueFilename(file.Filename)
	storagePath := filepath.Join(upload_config.Config.UploadDir, filename)

	// ----------Create upload directory if not exists-----------
	if err := os.MkdirAll(upload_config.Config.UploadDir, os.ModePerm); err != nil {
		return nil, err
	}

	// -------------- Save file to storage ---------------
	dst, err := os.Create(storagePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(storagePath)
		return nil, err
	}

	// Determine file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	fileType := getFileType(ext, mimeType)

	// Create file record
	fileRecord := &File{
		FolderID:     folderID,
		Name:         strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename)),
		OriginalName: file.Filename,
		Extension:    ext,
		MimeType:     mimeType,
		Size:         uint64(file.Size),
		StoragePath:  storagePath,
		UploadedBy:   userID,
		FileType:     fileType,
		IsPublic:     false,
		IsDeleted:    false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// If service has baseURL configured, store full URL in DB instead of raw filesystem path
	if s.baseURL != "" {
		fileRecord.StoragePath = buildFullURL(s.baseURL, storagePath)
	}

	if err := s.repo.UploadFile(fileRecord); err != nil {
		os.Remove(storagePath)
		return nil, err
	}

	// -----------Update storage usage--------
	s.repo.UpdateStorageUsage(userID, uint64(file.Size), true)

	//-------------  Add history ----------------
	s.repo.AddHistory(&FileHistory{
		FileID:      fileRecord.ID,
		Action:      FileActionUpload,
		PerformedBy: userID,
		Details:     toJSON(map[string]interface{}{"filename": file.Filename, "size": file.Size}),
		CreatedAt:   time.Now(),
	})

	return fileRecord, nil
}

func (s *Service) UploadMultiple(files []*multipart.FileHeader, folderID *uint64, userID uint64) ([]File, []string, error) {
	if len(files) == 0 {
		return nil, nil, fmt.Errorf("no files provided")
	}

	if len(files) > upload_config.Config.MaxFiles {
		return nil, nil, fmt.Errorf("maximum %d files can be uploaded at once", upload_config.Config.MaxFiles)
	}

	var uploaded []File
	var failed []string

	for _, file := range files {
		uploadedFile, err := s.UploadFile(file, folderID, userID)
		if err != nil {
			failed = append(failed, file.Filename+": "+err.Error())
		} else {
			uploaded = append(uploaded, *uploadedFile)
		}
	}

	return uploaded, failed, nil
}

func (s *Service) GetFile(fileID uint64) (*File, error) {
	return s.repo.GetFileByID(fileID)
}

func (s *Service) GetUserFiles(userID uint64, limit int, page int) ([]File, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	return s.repo.GetUserFiles(userID, limit, offset)
}

func (s *Service) GetFolderFiles(folderID, userID uint64, limit, page int) ([]File, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	return s.repo.GetFolderFiles(folderID, userID, limit, offset)
}

func (s *Service) SearchFiles(query string, userID uint64, limit int, page int) ([]File, int64, error) {
	if strings.TrimSpace(query) == "" {
		return nil, 0, fmt.Errorf("search query cannot be empty")
	}

	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	return s.repo.SearchFiles(query, userID, limit, offset)
}
//--------------------- update metadata------
func (s *Service) UpdateFileMetadata(fileID uint64, name string) (*File, error) {
	file, err := s.repo.GetFileByID(fileID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(name) != "" {
		file.Name = name
	}

	file.UpdatedAt = time.Now()

	if err := s.repo.UpdateFile(file); err != nil {
		return nil, err
	}

	return file, nil
}

func (s *Service) DeleteFile(fileID uint64) error {
	file, err := s.repo.GetFileByID(fileID)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteFile(fileID); err != nil {
		return err
	}

	//--------- Update storage usage-----------
	s.repo.UpdateStorageUsage(file.UploadedBy, file.Size, false)

	// ------------Add history-------------
	s.repo.AddHistory(&FileHistory{
		FileID:      fileID,
		Action:      FileActionDelete,
		PerformedBy: file.UploadedBy,
		Details:     toJSON(map[string]interface{}{"filename": file.OriginalName}),
		CreatedAt:   time.Now(),
	})

	return nil
}

func (s *Service) RestoreFile(fileID uint64) error {
	file, err := s.repo.GetFileByID(fileID)
	if err != nil {
		return err
	}

	if err := s.repo.RestoreFile(fileID); err != nil {
		return err
	}

	// Update storage usage
	s.repo.UpdateStorageUsage(file.UploadedBy, file.Size, true)

	// Add history
	s.repo.AddHistory(&FileHistory{
		FileID:      fileID,
		Action:      FileActionRestore,
		PerformedBy: file.UploadedBy,
		Details:     toJSON(map[string]interface{}{"filename": file.OriginalName}),
		CreatedAt:   time.Now(),
	})

	return nil
}

func (s *Service) MoveFile(fileID uint64, newFolderID *uint64) (*File, error) {
	file, err := s.repo.GetFileByID(fileID)
	if err != nil {
		return nil, err
	}

	// Validate new folder if specified
	if newFolderID != nil {
		exists, err := s.repo.CheckFolderExists(*newFolderID)
		if err != nil || !exists {
			return nil, fmt.Errorf("destination folder not found")
		}
	}

	oldFolderID := file.FolderID
	file.FolderID = newFolderID
	file.UpdatedAt = time.Now()

	if err := s.repo.UpdateFile(file); err != nil {
		return nil, err
	}

	//-------- Add history----------
	details := map[string]interface{}{
		"filename": file.OriginalName,
		"from":     oldFolderID,
		"to":       newFolderID,
	}
	s.repo.AddHistory(&FileHistory{
		FileID:      fileID,
		Action:      FileActionMove,
		PerformedBy: file.UploadedBy,
		Details:     toJSON(details),
		CreatedAt:   time.Now(),
	})

	return file, nil
}

func (s *Service) DownloadFile(fileID uint64) (string, error) {
	file, err := s.repo.GetFileByID(fileID)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(file.StoragePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found on storage")
	}

	return file.StoragePath, nil
}



// ==================== STORAGE OPERATIONS ====================

func (s *Service) GetStorageUsage(userID uint64) (*StorageUsage, error) {
	return s.repo.GetStorageUsage(userID)
}

// ==================== HISTORY OPERATIONS ====================

func (s *Service) GetFileHistory(fileID uint64, limit int) ([]FileHistory, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.GetFileHistory(fileID, limit)
}

// ==================== HELPER FUNCTIONS ====================

func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	name := strings.TrimSuffix(originalFilename, ext)
	hash := md5.Sum([]byte(name + time.Now().String()))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr + ext
}

func getFileType(ext string, mimeType string) FileTypeEnum {
	ext = strings.ToLower(ext)

	// -----------Image types--------------
	imageExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".bmp": true}
	if imageExts[ext] {
		return FileTypeImage
	}

	// Document types
	docExts := map[string]bool{".pdf": true, ".doc": true, ".docx": true, ".txt": true, ".xlsx": true, ".xls": true, ".ppt": true, ".pptx": true}
	if docExts[ext] {
		return FileTypeDocument
	}

	// Video types
	videoExts := map[string]bool{".mp4": true, ".avi": true, ".mov": true, ".mkv": true, ".webm": true}
	if videoExts[ext] {
		return FileTypeVideo
	}

	// -----------Audio types----------
	audioExts := map[string]bool{".mp3": true, ".wav": true, ".flac": true, ".aac": true, ".m4a": true}
	if audioExts[ext] {
		return FileTypeAudio
	}

	//--------- Archive types---------
	archiveExts := map[string]bool{".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true}
	if archiveExts[ext] {
		return FileTypeArchive
	}

	return FileTypeOther
}

func toJSON(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

// buildFullURL ensures baseURL and path are joined correctly
func buildFullURL(baseURL, path string) string {
	// If path already looks like a URL, return as is
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	// Trim leading ./ or / from path when joining
	p := strings.TrimPrefix(path, "./")
	p = strings.TrimLeft(p, "/")

	base := strings.TrimRight(baseURL, "/")
	return base + "/" + p
}