package projects

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
	rdb     *redis.Client
}

func NewHandler(service *Service, rdb *redis.Client) *Handler {
	return &Handler{service: service, rdb: rdb}
}

// =============================
// HANDLER METHODS
// =============================

// CreateProject godoc
// @Summary Create a new project
// @Description Creates a new project in the system
// @Tags purchase ----------------------- Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateProjectRequest true "create project payload"
// @Success 201 {object} ProjectResponse
// @Router /purchase/projects [post]
func (h *Handler) CreateProject(ctx *gin.Context) {
	var req CreateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	cmd, err := h.convertToCreateCommand(&req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid input data", "VALIDATION_ERROR", err.Error())
		return
	}

	project, err := h.service.CreateProject(ctx.Request.Context(), cmd)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.Success(ctx, "project created successfully", h.toResponse(project))
}



// ListProjects godoc
// @Summary List projects with pagination
// @Description Retrieves projects with pagination and filters
// @Tags purchase ----------------------- Projects
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param status query string false "Project status"
// @Param office_id query int false "Office ID"
// @Param manager_id query int false "Manager ID"
// @Success 200 {object} ListProjectsResponse
// @Router /purchase/projects [get]
func (h *Handler) ListProjects(ctx *gin.Context) {
	var req ListProjectsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid query parameters", "VALIDATION_ERROR", err.Error())
		return
	}

	filters := make(map[string]interface{})
	if req.Status != nil {
		filters["status"] = *req.Status
	}
	if req.OfficeID != nil {
		filters["office_id"] = *req.OfficeID
	}
	if req.ManagerID != nil {
		filters["manager_id"] = *req.ManagerID
	}

	projects, total, err := h.service.GetProjectsPaginated(ctx.Request.Context(), req.Page, req.Limit, filters)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	data := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		data[i] = h.toResponse(&p)
	}

	totalPages := (total + int64(req.Limit) - 1) / int64(req.Limit)
	resp := ListProjectsResponse{
		Data:       data,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "projects retrieved successfully",
		"data":    resp,
	})
}

// UpdateProject godoc
// @Summary Update a project
// @Description Updates an existing project
// @Tags purchase ----------------------- Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param payload body UpdateProjectRequest true "update project payload"
// @Success 200 {object} ProjectResponse
// @Router /purchase/projects/{id} [patch]
func (h *Handler) UpdateProject(ctx *gin.Context) {
	id, err := h.parseID(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid project id", "VALIDATION_ERROR", err.Error())
		return
	}

	var req UpdateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	cmd, err := h.convertToUpdateCommand(&req)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid input data", "VALIDATION_ERROR", err.Error())
		return
	}

	project, err := h.service.UpdateProject(ctx.Request.Context(), id, cmd)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.Success(ctx, "project updated successfully", h.toResponse(project))
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Soft deletes a project
// @Tags purchase ----------------------- Projects
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {object} response.APIResponse
// @Router /purchase/projects/{id} [delete]
func (h *Handler) DeleteProject(ctx *gin.Context) {
	id, err := h.parseID(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid project id", "VALIDATION_ERROR", err.Error())
		return
	}

	if err := h.service.DeleteProject(ctx.Request.Context(), id); err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.Success(ctx, "project deleted successfully", nil)
}

// ChangeProjectStatus godoc
// @Summary Change project status
// @Description Changes the status of a project
// @Tags purchase ----------------------- Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Param payload body ChangeStatusRequest true "status payload"
// @Success 200 {object} ProjectResponse
// @Router /purchase/projects/{id}/status [put]
func (h *Handler) ChangeProjectStatus(ctx *gin.Context) {
	id, err := h.parseID(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid project id", "VALIDATION_ERROR", err.Error())
		return
	}

	var req ChangeStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	project, err := h.service.ChangeProjectStatus(ctx.Request.Context(), id, req.Status)
	if err != nil {
		h.handleServiceError(ctx, err)
		return
	}

	response.Success(ctx, "project status changed successfully", h.toResponse(project))
}









// =============================
// HELPER METHODS
// =============================

// parseID parses string ID to uint
func (h *Handler) parseID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// convertToCreateCommand converts API DTO to service command
func (h *Handler) convertToCreateCommand(req *CreateProjectRequest) (CreateProjectCmd, error) {
	cmd := CreateProjectCmd{
		ProjectName: req.ProjectName,
		Description: req.Description,
		ProjectCode: req.ProjectCode,
		Budget:      req.Budget,
		Status:      req.Status,
		IsActive:    req.IsActive,
		ManagerID:   req.ManagerID,
		OfficeID:    req.OfficeID,
	}

	// Parse dates if provided
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return cmd, err
		}
		cmd.StartDate = &t
	}

	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return cmd, err
		}
		cmd.EndDate = &t
	}

	return cmd, nil
}

// convertToUpdateCommand converts API DTO to service command
func (h *Handler) convertToUpdateCommand(req *UpdateProjectRequest) (UpdateProjectCmd, error) {
	cmd := UpdateProjectCmd{
		ProjectName: req.ProjectName,
		Description: req.Description,
		ProjectCode: req.ProjectCode,
		Budget:      req.Budget,
		Status:      req.Status,
		IsActive:    req.IsActive,
	
	}

	// Parse dates if provided
	if req.StartDate != nil {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return cmd, err
		}
		cmd.StartDate = &t
	}

	if req.EndDate != nil {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return cmd, err
		}
		cmd.EndDate = &t
	}

	return cmd, nil
}

// toResponse converts Project model to ProjectResponse DTO
func (h *Handler) toResponse(project *Project) ProjectResponse {
	resp := ProjectResponse{
		ID:          project.ID,
		ProjectName: project.ProjectName,
		Description: project.Description,
		ProjectCode: project.ProjectCode,
		Status:      project.Status,
		IsActive:    project.IsActive,
		CreatedAt:   project.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   project.UpdatedAt.Format("2006-01-02 15:04:05"),
		Budget:      project.Budget,
	}

	if project.StartDate != nil {
		startDateStr := project.StartDate.Format("2006-01-02")
		resp.StartDate = &startDateStr
	}

	if project.EndDate != nil {
		endDateStr := project.EndDate.Format("2006-01-02")
		resp.EndDate = &endDateStr
	}

	// Map manager details if available
	if project.Manager != nil {
		resp.Manager = &ManagerResponse{
			ID:       project.Manager.ID,
			FullName: project.Manager.Name,
			Email:    project.Manager.Email,
		}
	}

	// Map office details if available
	if project.Office.ID > 0 {
		resp.Office = &OfficeResponse{
			ID:         project.Office.ID,
			OfficeName: project.Office.Name,
			Location:   project.Office.Location,
		}
	}

	return resp
}

// handleServiceError maps service errors to HTTP responses
func (h *Handler) handleServiceError(ctx *gin.Context, err error) {
	if errors.Is(err, ErrProjectNotFound) {
		response.Error(ctx, http.StatusNotFound, "project not found", "NOT_FOUND", nil)
		return
	}
	if errors.Is(err, ErrProjectCodeAlreadyExists) {
		response.Error(ctx, http.StatusConflict, "project code already exists", "CONFLICT", nil)
		return
	}
	if errors.Is(err, ErrCannotDeleteCompletedProject) {
		response.Error(ctx, http.StatusConflict, "cannot delete a completed project", "CONFLICT", nil)
		return
	}
	if errors.Is(err, ErrCannotDeleteCancelledProject) {
		response.Error(ctx, http.StatusConflict, "cannot delete a cancelled project", "CONFLICT", nil)
		return
	}
	if errors.Is(err, ErrInvalidDateRange) {
		response.Error(ctx, http.StatusBadRequest, "end date must be after start date", "VALIDATION_ERROR", nil)
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Error(ctx, http.StatusNotFound, "record not found", "NOT_FOUND", nil)
		return
	}

	response.Error(ctx, http.StatusInternalServerError, "internal server error", "INTERNAL_ERROR", err.Error())
}