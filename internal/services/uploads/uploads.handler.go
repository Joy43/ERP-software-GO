package uploads

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/config"
)

type Handler struct {
	service IService
	cfg     config.Config
}

func NewHandler(service IService, cfg config.Config) *Handler {
	return &Handler{service: service, cfg: cfg}
}

func getUserID(c *gin.Context) (uint64, bool) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	switch v := userIDInterface.(type) {
	case uint:
		return uint64(v), true
	case uint64:
		return v, true
	case float64:
		return uint64(v), true
	case int:
		return uint64(v), true
	case int64:
		return uint64(v), true
	default:
		return 0, false
	}
}

// ==================== FOLDER HANDLERS ====================

// CreateFolder godoc
// @Summary Create a new folder No need parent id if you want to create root folder
// @Description Create a new folder in the specified parent or root
// @Tags Uploads----------Folders
// @Accept json
// @Produce json
// @Param request body CreateFolderRequest true "Folder details"
// @Success 201 {object} FolderResponse
// @Failure 400 {object} ErrorResponse
// @Router /uploads/folders [post]
// @Security BearerAuth
func (h *Handler) CreateFolder(c *gin.Context) {
	var req CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request body", err.Error())
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		Unauthorized(c, "User ID not found in context")
		return
	}

	folder, err := h.service.CreateFolder(req.Name, req.ParentID, userID)
	if err != nil {
		BadRequest(c, "Failed to create folder", err.Error())
		return
	}

	folderResp := h.folderToResponse(folder)
	Success(c, http.StatusCreated, folderResp, "Folder created successfully")
}

// GetRootFolders godoc
// @Summary Get root folders with subfolders
// @Description Get all root folders (folders without parent) for the user with their child folders nested
// @Tags Uploads----------Folders
// @Produce json
// @Success 200 {object} []FolderResponseWithChildren
// @Router /uploads/folders [get]
// @Security BearerAuth
func (h *Handler) GetRootFolders(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		Unauthorized(c, "User ID not found in context")
		return
	}

	folders, err := h.service.GetRootFolders(userID)
	if err != nil {
		InternalServerError(c, "Failed to fetch folders", err.Error())
		return
	}

	var folderResps []FolderResponseWithChildren
	for _, folder := range folders {
		folderResps = append(folderResps, h.folderToResponseWithChildren(&folder))
	}

	Success(c, http.StatusOK, folderResps, "Folders retrieved successfully")
}
// GetFolderContents godoc
// @Summary Get folder contents like files 
// @Description Get subfolders and files in a folder with pagination. If no folder_id is provided, returns root folder contents.
// @Tags Uploads----------Folders
// @Produce json
// @Param folder_id path integer false "Folder ID (optional, defaults to root)"
// @Param page query integer false "Page number" default(1)
// @Param limit query integer false "Items per page" default(20)
// @Success 200 {object} FolderContentsResponse
// @Router /uploads/folders/{folder_id}/contents [get]
// @Router /uploads/folders/contents [get]
// @Security BearerAuth
func (h *Handler) GetFolderContents(c *gin.Context) {
	var folderID *uint64
	
	folderIDStr := c.Param("folder_id")

	if folderIDStr != "" && folderIDStr != "contents" && folderIDStr != "undefined" {
		id, err := strconv.ParseUint(folderIDStr, 10, 64)
		if err != nil {
			BadRequest(c, "Invalid folder ID", err.Error())
			return
		}
		folderID = &id
	}

	userID, ok := getUserID(c)
	if !ok {
		Unauthorized(c, "User ID not found in context")
		return
	}

	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	if l := c.Query("limit"); l != "" {
		if limitNum, err := strconv.Atoi(l); err == nil && limitNum > 0 && limitNum <= 100 {
			limit = limitNum
		}
	}

	// If no folder_id provided, default to the first root folder for the user (if any)
	if folderID == nil {
		rootFolders, err := h.service.GetRootFolders(userID)
		if err == nil && len(rootFolders) > 0 {
			id := rootFolders[0].ID
			folderID = &id
		}
	}

	folders, folderTotal, files, fileTotal, err := h.service.GetFolderContents(folderID, userID, limit, page)
	if err != nil {
		InternalServerError(c, "Failed to fetch folder contents", err.Error())
		return
	}

	// Initialize slices as empty instead of nil to avoid null in JSON response
	folderResps := make([]FolderResponse, 0)
	fileResps := make([]FileResponse, 0)

	for _, folder := range folders {
		folderResps = append(folderResps, h.folderToResponse(&folder))
	}

	for _, file := range files {
		fileResps = append(fileResps, h.fileToResponse(&file))
	}

	resp := FolderContentsResponse{
		Folders:     folderResps,
		Files:       fileResps,
		FolderTotal: folderTotal,
		FileTotal:   fileTotal,
		Page:        page,
		Limit:       limit,
	}

	Success(c, http.StatusOK, resp, "Contents retrieved successfully")
}


// UpdateFolder godoc
// @Summary Update folder
// @Description Update folder details
// @Tags Uploads----------Folders
// @Accept json
// @Produce json
// @Param folder_id path integer true "Folder ID"
// @Param request body UpdateFolderRequest true "Update details"
// @Success 200 {object} FolderResponse
// @Router /uploads/folders/{folder_id} [patch]
// @Security BearerAuth
func (h *Handler) UpdateFolder(c *gin.Context) {
	folderID, err := strconv.ParseUint(c.Param("folder_id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid folder ID", err.Error())
		return
	}

	var req UpdateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request body", err.Error())
		return
	}

	folder, err := h.service.UpdateFolder(folderID, req.Name)
	if err != nil {
		BadRequest(c, "Failed to update folder", err.Error())
		return
	}

	Success(c, http.StatusOK, h.folderToResponse(folder), "Folder updated successfully")
}

// DeleteFolder godoc
// @Summary Delete folder
// @Description Soft delete a folder and its contents
// @Tags Uploads----------Folders
// @Produce json
// @Param folder_id path integer true "Folder ID"
// @Success 200 {object} MessageResponse
// @Router /uploads/folders/{folder_id} [delete]
// @Security BearerAuth
func (h *Handler) DeleteFolder(c *gin.Context) {
	folderID, err := strconv.ParseUint(c.Param("folder_id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid folder ID", err.Error())
		return
	}

	if err := h.service.DeleteFolder(folderID); err != nil {
		InternalServerError(c, "Failed to delete folder", err.Error())
		return
	}

	Success(c, http.StatusOK, nil, "Folder deleted successfully")
}

// RestoreFolder godoc
// @Summary Restore folder
// @Description Restore a deleted folder
// @Tags Uploads----------Folders
// @Produce json
// @Param folder_id path integer true "Folder ID"
// @Success 200 {object} MessageResponse
// @Router /uploads/folders/{folder_id}/restore [post]
// @Security BearerAuth
func (h *Handler) RestoreFolder(c *gin.Context) {
	folderID, err := strconv.ParseUint(c.Param("folder_id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid folder ID", err.Error())
		return
	}

	if err := h.service.RestoreFolder(folderID); err != nil {
		InternalServerError(c, "Failed to restore folder", err.Error())
		return
	}

	Success(c, http.StatusOK, nil, "Folder restored successfully")
}

// ==================== FILE HANDLERS ====================

// UploadFile godoc
// @Summary Upload a single file
// @Description Upload a single file to a specific folder or root
// @Tags Uploads----------Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param folder_id formData integer false "Folder ID"
// @Success 201 {object} FileResponse
// @Router /uploads/files [post]
// @Security BearerAuth
func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		BadRequest(c, "File is required", err.Error())
		return
	}

	userID, ok := getUserID(c)
if !ok {
Unauthorized(c, "User ID not found in context")
return
}

	//------- Parse folder ID if provided---------
	var folderID *uint64
	if folderIDStr := c.PostForm("folder_id"); folderIDStr != "" {
		if id, err := strconv.ParseUint(folderIDStr, 10, 64); err == nil {
			folderID = &id
		}
	}

	uploadedFile, err := h.service.UploadFile(file, folderID, userID)
	if err != nil {
		BadRequest(c, "File upload failed", err.Error())
		return
	}

	Success(c, http.StatusCreated, h.fileToResponse(uploadedFile), "File uploaded successfully")
}

// UploadMultiple godoc
// @Summary Upload multiple files
// @Description Upload multiple files at once
// @Tags Uploads----------Files
// @Accept multipart/form-data
// @Produce json
// @Param files formData []file true "Files to upload"
// @Param folder_id formData integer false "Folder ID"
// @Success 201 {object} UploadResponse
// @Router /uploads/files/batch [post]
// @Security BearerAuth
func (h *Handler) UploadMultiple(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		BadRequest(c, "Failed to parse multipart form", err.Error())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		BadRequest(c, "No files provided", "")
		return
	}

	userID, ok := getUserID(c)
if !ok {
Unauthorized(c, "User ID not found in context")
return
}

	// --------- Parse folder ID if provided ---------
	var folderID *uint64
	if folderIDStr := c.PostForm("folder_id"); folderIDStr != "" {
		if id, err := strconv.ParseUint(folderIDStr, 10, 64); err == nil {
			folderID = &id
		}
	}

	uploaded, failed, err := h.service.UploadMultiple(files, folderID, userID)
	if err != nil {
		BadRequest(c, "Batch upload failed", err.Error())
		return
	}

	var uploadedResps []FileResponse
	for _, file := range uploaded {
		uploadedResps = append(uploadedResps, h.fileToResponse(&file))
	}

	resp := UploadResponse{
		Files:  uploadedResps,
		Failed: failed,
	}

	Success(c, http.StatusCreated, resp, "Files uploaded successfully")
}


// GetUserFiles godoc
// @Summary Get user files
// @Description Get all files uploaded by the user with optional search and pagination
// @Tags Uploads----------Files
// @Produce json
// @Param q query string false "Search query (optional)"
// @Param page query integer false "Page number" default(1)
// @Param limit query integer false "Items per page" default(20)
// @Success 200 {object} ListResponse
// @Router /uploads/files [get]
// @Security BearerAuth
func (h *Handler) GetUserFiles(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		Unauthorized(c, "User ID not found in context")
		return
	}

	page, limit := 1, 20
	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	if l := c.Query("limit"); l != "" {
		if limitNum, err := strconv.Atoi(l); err == nil && limitNum > 0 && limitNum <= 100 {
			limit = limitNum
		}
	}

	var (files []File; total int64; err error)
	if q := c.Query("q"); q != "" {
		files, total, err = h.service.SearchFiles(q, userID, limit, page)
	} else {
		files, total, err = h.service.GetUserFiles(userID, limit, page)
	}
	if err != nil {
		InternalServerError(c, "Failed to fetch files", err.Error())
		return
	}

	var fileResps []FileResponse
	for _, file := range files {
		fileResps = append(fileResps, h.fileToResponse(&file))
	}
	SuccessList(c, http.StatusOK, fileResps, total, page, limit, "Files retrieved successfully")
}

//===============================
// Soft delete file (move to trash)
//===============================


// DeleteFile godoc
// @Summary Delete file
// @Description Soft delete a file
// @Tags Uploads----------Files
// @Produce json
// @Param file_id path integer true "File ID"
// @Success 200 {object} MessageResponse
// @Router /uploads/files/{file_id} [delete]
// @Security BearerAuth
func (h *Handler) DeleteFile(c *gin.Context) {
	fileID, err := strconv.ParseUint(c.Param("file_id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid file ID", err.Error())
		return
	}

	if err := h.service.DeleteFile(fileID); err != nil {
		InternalServerError(c, "Failed to delete file", err.Error())
		return
	}

	Success(c, http.StatusOK, nil, "File deleted successfully")
}
//===============================
// restore file from trash (soft deleted)
//===============================

// RestoreFile godoc
// @Summary Restore file
// @Description Restore a deleted file
// @Tags Uploads----------Files
// @Produce json
// @Param file_id path integer true "File ID"
// @Success 200 {object} MessageResponse
// @Router /uploads/files/{file_id}/restore [post]
// @Security BearerAuth
func (h *Handler) RestoreFile(c *gin.Context) {
	fileID, err := strconv.ParseUint(c.Param("file_id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid file ID", err.Error())
		return
	}

	if err := h.service.RestoreFile(fileID); err != nil {
		InternalServerError(c, "Failed to restore file", err.Error())
		return
	}

	Success(c, http.StatusOK, nil, "File restored successfully")
}
//===============================
// Move file to another folder
//===============================

// MoveFile godoc
// @Summary Move file to another folder
// @Description Move a file to a different folder
// @Tags Uploads----------Files
// @Accept json
// @Produce json
// @Param file_id path integer true "File ID"
// @Param request body MoveFileRequest true "Destination folder"
// @Success 200 {object} FileResponse
// @Router /uploads/files/{file_id}/move [post]
// @Security BearerAuth
func (h *Handler) MoveFile(c *gin.Context) {
	fileID, err := strconv.ParseUint(c.Param("file_id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid file ID", err.Error())
		return
	}

	var req MoveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request body", err.Error())
		return
	}

	file, err := h.service.MoveFile(fileID, req.FolderID)
	if err != nil {
		BadRequest(c, "Failed to move file", err.Error())
		return
	}

	Success(c, http.StatusOK, h.fileToResponse(file), "File moved successfully")
}

//================== view file inline (for supported types) ==================

// DownloadFile godoc
// @Summary Download file
// @Description Download a file
// @Tags Uploads----------Files
// @Produce octet-stream
// @Param file_id path integer true "File ID"
// @Success 200
// @Router /uploads/files/{file_id}/download [get]
// @Security BearerAuth
func (h *Handler) DownloadFile(c *gin.Context) {
	fileID, err := strconv.ParseUint(c.Param("file_id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid file ID", err.Error())
		return
	}

	storagePath, err := h.service.DownloadFile(fileID)
	if err != nil {
		NotFound(c, "File not found")
		return
	}

	file, err := h.service.GetFile(fileID)
	if err != nil {
		NotFound(c, "File not found")
		return
	}

	c.FileAttachment(storagePath, file.OriginalName)
}



// ==================== HELPER METHODS ====================

func (h *Handler) folderToResponse(folder *Folder) FolderResponse {
	return FolderResponse{
		ID:          folder.ID,
		Name:        folder.Name,
		ParentID:    folder.ParentID,
		Level:       folder.Level,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
	}
}

func (h *Handler) fileToResponse(file *File) FileResponse {
	var folderName string
	if file.Folder != nil {
		folderName = file.Folder.Name
	}

	// Build full URL only if storage path is not already a URL
	storagePath := file.StoragePath
	if storagePath != "" && h.cfg.BaseURL != "" {
		if !(strings.HasPrefix(storagePath, "http://") || strings.HasPrefix(storagePath, "https://")) {
			// Trim leading ./ or / from path when joining
			p := strings.TrimPrefix(storagePath, "./")
			p = strings.TrimLeft(p, "/")
			storagePath = strings.TrimRight(h.cfg.BaseURL, "/") + "/" + p
		}
	}

	return FileResponse{
		ID:           file.ID,
		Name:         file.Name,
		OriginalName: file.OriginalName,
		Extension:    file.Extension,
		MimeType:     file.MimeType,
		Size:         file.Size,
		FileType:     string(file.FileType),
		IsPublic:     file.IsPublic,
		StoragePath:  storagePath,
		CreatedAt:    file.CreatedAt,
		UpdatedAt:    file.UpdatedAt,
		FolderName:   folderName,
	}
}

func (h *Handler) fileDetailToResponse(file *File) FileDetailResponse {
	var folderResp *FolderResponse
	if file.Folder != nil {
		resp := h.folderToResponse(file.Folder)
		folderResp = &resp
	}

	return FileDetailResponse{
		ID:           file.ID,
		Name:         file.Name,
		OriginalName: file.OriginalName,
		Extension:    file.Extension,
		MimeType:     file.MimeType,
		Size:         file.Size,
		FileType:     string(file.FileType),
		IsPublic:     file.IsPublic,
		StoragePath:  file.StoragePath,
		UploadedBy:   file.UploadedBy,
		CreatedAt:    file.CreatedAt,
		UpdatedAt:    file.UpdatedAt,
		Folder:       folderResp,
	}
}

func (h *Handler) historyToResponse(h2 *FileHistory) FileHistoryResponse {
	return FileHistoryResponse{
		ID:          h2.ID,
		FileID:      h2.FileID,
		Action:      string(h2.Action),
		PerformedBy: h2.PerformedBy,
		Details:     h2.Details,
		CreatedAt:   h2.CreatedAt,
	}
}

func (h *Handler) folderToResponseWithChildren(folder *Folder) FolderResponseWithChildren {
	resp := FolderResponseWithChildren{
		ID:          folder.ID,
		Name:        folder.Name,
		ParentID:    folder.ParentID,
		Level:       folder.Level,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
	}

	// --------- Recursively convert children------------
	if len(folder.Children) > 0 {
		resp.Children = make([]FolderResponseWithChildren, 0, len(folder.Children))
		for _, child := range folder.Children {
			resp.Children = append(resp.Children, h.folderToResponseWithChildren(&child))
		}
	}

	return resp
}