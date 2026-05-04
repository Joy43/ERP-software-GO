package uploads

import (
	"time"

	"gorm.io/gorm"
)

type IUploadRepository interface {
	//-------- Folder operations------------
	CreateFolder(folder *Folder) error
	GetFolderByID(id uint64) (*Folder, error)
	GetRootFolders(userID uint64) ([]Folder, error)
	GetFolderContents(folderID *uint64, userID uint64, limit int, offset int) ([]Folder, int64, []File, int64, error)
	GetFolderPath(folderID uint64) ([]Folder, error)
	UpdateFolder(folder *Folder) error
	DeleteFolder(id uint64) error
	RestoreFolder(id uint64) error
	CheckFolderExists(folderID uint64) (bool, error)
	CheckDuplicateFolderName(parentID *uint64, name string, userID uint64) (bool, error)

	// File operations
	UploadFile(file *File) error
	GetFileByID(id uint64) (*File, error)
	GetFilesByFolder(folderID uint64) ([]File, error)
	GetUserFiles(userID uint64, limit int, offset int) ([]File, int64, error)
	GetFolderFiles(folderID, userID uint64, limit int, offset int) ([]File, int64, error)
	UpdateFile(file *File) error
	DeleteFile(id uint64) error
	RestoreFile(id uint64) error
	CheckDuplicateFileName(folderID *uint64, name string, uploadedBy uint64) (bool, error)
	SearchFiles(query string, userID uint64, limit int, offset int) ([]File, int64, error)

	// Storage usage
	GetStorageUsage(userID uint64) (*StorageUsage, error)
	UpdateStorageUsage(userID uint64, fileSize uint64, increment bool) error

	// File history
	AddHistory(history *FileHistory) error
	GetFileHistory(fileID uint64, limit int) ([]FileHistory, error)
	GetUserHistory(userID uint64, limit int) ([]FileHistory, error)

	// Utility operations
	CheckUserExists(userID uint64) (bool, error)
}

// Repository implements IUploadRepository
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new upload repository
func NewRepository(db *gorm.DB) IUploadRepository {
	return &Repository{db: db}
}

// ==================== FOLDER OPERATIONS ====================

func (r *Repository) CreateFolder(folder *Folder) error {
	return r.db.Create(folder).Error
}

func (r *Repository) GetFolderByID(id uint64) (*Folder, error) {
	var folder Folder
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&folder).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *Repository) GetRootFolders(userID uint64) ([]Folder, error) {
	var folders []Folder
	err := r.db.Where("parent_id IS NULL AND created_by = ? AND is_deleted = ?", userID, false).
		Order("name ASC").
		Preload("Children", "is_deleted = ?", false).
		Preload("Children.Children", "is_deleted = ?", false).
		Preload("Children.Children.Children", "is_deleted = ?", false).
		Preload("Children.Children.Children.Children", "is_deleted = ?", false).
		Find(&folders).Error
	return folders, err
}

func (r *Repository) GetFolderContents(folderID *uint64, userID uint64, limit int, offset int) ([]Folder, int64, []File, int64, error) {
	var folders []Folder
	var files []File
	var folderTotal int64
	var fileTotal int64

	// Get total folders count
	var errFCount error
	if folderID == nil {
		// Root folder contents
		errFCount = r.db.Model(&Folder{}).
			Where("parent_id IS NULL AND created_by = ? AND is_deleted = ?", userID, false).
			Count(&folderTotal).Error
	} else {
		// Specific folder contents
		errFCount = r.db.Model(&Folder{}).
			Where("parent_id = ? AND created_by = ? AND is_deleted = ?", *folderID, userID, false).
			Count(&folderTotal).Error
	}

	// Get subfolders with pagination
	var errF error
	if folderID == nil {
		// Root folder contents
		errF = r.db.Where("parent_id IS NULL AND created_by = ? AND is_deleted = ?", userID, false).
			Order("name ASC").
			Limit(limit).
			Offset(offset).
			Find(&folders).Error
	} else {
		// Specific folder contents
		errF = r.db.Where("parent_id = ? AND created_by = ? AND is_deleted = ?", *folderID, userID, false).
			Order("name ASC").
			Limit(limit).
			Offset(offset).
			Find(&folders).Error
	}

	// Get total files count
	var errFiCount error
	if folderID == nil {
		// Root folder files - get all user files (both with NULL and assigned folder_id)
		errFiCount = r.db.Model(&File{}).
			Where("uploaded_by = ? AND is_deleted = ?", userID, false).
			Count(&fileTotal).Error
	} else {
		// Specific folder files
		errFiCount = r.db.Model(&File{}).
			Where("folder_id = ? AND uploaded_by = ? AND is_deleted = ?", *folderID, userID, false).
			Count(&fileTotal).Error
	}

	// Get files with pagination
	var errFi error
	if folderID == nil {
		// Root folder files - get all user files (both with NULL and assigned folder_id)
		errFi = r.db.Where("uploaded_by = ? AND is_deleted = ?", userID, false).
			Order("created_at DESC").
			Limit(limit).
			Offset(offset).
			Find(&files).Error
	} else {
		// Specific folder files
		errFi = r.db.Where("folder_id = ? AND uploaded_by = ? AND is_deleted = ?", *folderID, userID, false).
			Order("created_at DESC").
			Limit(limit).
			Offset(offset).
			Find(&files).Error
	}

	if errFCount != nil {
		return nil, 0, nil, 0, errFCount
	}
	if errF != nil {
		return nil, 0, nil, 0, errF
	}
	if errFiCount != nil {
		return nil, 0, nil, 0, errFiCount
	}
	if errFi != nil {
		return nil, 0, nil, 0, errFi
	}
	return folders, folderTotal, files, fileTotal, nil
}

func (r *Repository) GetFolderPath(folderID uint64) ([]Folder, error) {
	
	var result []Folder
	var current *Folder
	var err error

	for {
		current, err = r.GetFolderByID(folderID)
		if err != nil || current == nil {
			break
		}
		result = append(result, *current)
		if current.ParentID == nil {
			break
		}
		folderID = *current.ParentID
	}

	return result, nil
}

func (r *Repository) UpdateFolder(folder *Folder) error {
	return r.db.Model(folder).Updates(folder).Error
}

func (r *Repository) DeleteFolder(id uint64) error {
	// Soft delete folder and all subfolders
	return r.db.Model(&Folder{}).
		Where("id = ? OR parent_id = ?", id, id).
		Update("is_deleted", true).Error
}

func (r *Repository) RestoreFolder(id uint64) error {
	return r.db.Model(&Folder{}).
		Where("id = ?", id).
		Update("is_deleted", false).Error
}

func (r *Repository) CheckFolderExists(folderID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&Folder{}).
		Where("id = ? AND is_deleted = ?", folderID, false).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) CheckDuplicateFolderName(parentID *uint64, name string, userID uint64) (bool, error) {
	var count int64
	query := r.db.Model(&Folder{}).
		Where("created_by = ? AND name = ? AND is_deleted = ?", userID, name, false)
	
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", parentID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// ==================== FILE OPERATIONS ====================

func (r *Repository) UploadFile(file *File) error {
	return r.db.Create(file).Error
}

func (r *Repository) GetFileByID(id uint64) (*File, error) {
	var file File
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).
		Preload("Folder").
		First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *Repository) GetFilesByFolder(folderID uint64) ([]File, error) {
	var files []File
	err := r.db.Where("folder_id = ? AND is_deleted = ?", folderID, false).
		Preload("Folder").
		Order("created_at DESC").
		Find(&files).Error
	return files, err
}

func (r *Repository) GetUserFiles(userID uint64, limit int, offset int) ([]File, int64, error) {
	var files []File
	var total int64

	err := r.db.Model(&File{}).
		Where("uploaded_by = ? AND is_deleted = ?", userID, false).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("uploaded_by = ? AND is_deleted = ?", userID, false).
		Preload("Folder").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&files).Error

	return files, total, err
}

func (r *Repository) GetFolderFiles(folderID, userID uint64, limit int, offset int) ([]File, int64, error) {
	var files []File
	var total int64

	err := r.db.Model(&File{}).
		Where("folder_id = ? AND uploaded_by = ? AND is_deleted = ?", folderID, userID, false).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("folder_id = ? AND uploaded_by = ? AND is_deleted = ?", folderID, userID, false).
		Preload("Folder").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&files).Error

	return files, total, err
}

func (r *Repository) UpdateFile(file *File) error {
	return r.db.Model(file).Updates(file).Error
}

func (r *Repository) DeleteFile(id uint64) error {
	return r.db.Model(&File{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error
}

func (r *Repository) RestoreFile(id uint64) error {
	return r.db.Model(&File{}).
		Where("id = ?", id).
		Update("is_deleted", false).Error
}

func (r *Repository) CheckDuplicateFileName(folderID *uint64, name string, uploadedBy uint64) (bool, error) {
	var count int64
	query := r.db.Model(&File{}).
		Where("uploaded_by = ? AND name = ? AND is_deleted = ?", uploadedBy, name, false)
	
	if folderID == nil {
		query = query.Where("folder_id IS NULL")
	} else {
		query = query.Where("folder_id = ?", folderID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *Repository) SearchFiles(query string, userID uint64, limit int, offset int) ([]File, int64, error) {
	var files []File
	var total int64
	searchPattern := "%" + query + "%"

	err := r.db.Model(&File{}).
		Where("uploaded_by = ? AND is_deleted = ? AND (name LIKE ? OR original_name LIKE ?)", 
			userID, false, searchPattern, searchPattern).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("uploaded_by = ? AND is_deleted = ? AND (name LIKE ? OR original_name LIKE ?)", 
		userID, false, searchPattern, searchPattern).
		Preload("Folder").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&files).Error

	return files, total, err
}

// ==================== STORAGE USAGE ====================

func (r *Repository) GetStorageUsage(userID uint64) (*StorageUsage, error) {
	var usage StorageUsage
	err := r.db.Where("user_id = ?", userID).First(&usage).Error
	if err == gorm.ErrRecordNotFound {
		// Create new storage usage record
		usage.UserID = userID
		usage.TotalFiles = 0
		usage.TotalBytes = 0
		usage.UpdatedAt = time.Now()
		return &usage, r.db.Create(&usage).Error
	}
	return &usage, err
}

func (r *Repository) UpdateStorageUsage(userID uint64, fileSize uint64, increment bool) error {
	usage, err := r.GetStorageUsage(userID)
	if err != nil {
		return err
	}

	if increment {
		usage.TotalBytes += fileSize
		usage.TotalFiles++
	} else {
		if usage.TotalBytes >= fileSize {
			usage.TotalBytes -= fileSize
		}
		if usage.TotalFiles > 0 {
			usage.TotalFiles--
		}
	}
	usage.UpdatedAt = time.Now()

	return r.db.Save(usage).Error
}

// ==================== FILE HISTORY ====================

func (r *Repository) AddHistory(history *FileHistory) error {
	return r.db.Create(history).Error
}

func (r *Repository) GetFileHistory(fileID uint64, limit int) ([]FileHistory, error) {
	var histories []FileHistory
	err := r.db.Where("file_id = ?", fileID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

func (r *Repository) GetUserHistory(userID uint64, limit int) ([]FileHistory, error) {
	var histories []FileHistory
	err := r.db.Where("performed_by = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

// ==================== UTILITY OPERATIONS ====================

func (r *Repository) CheckUserExists(userID uint64) (bool, error) {
	var count int64
	err := r.db.Table("users").Where("id = ?", userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}