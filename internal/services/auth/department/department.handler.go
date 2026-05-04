package department

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateDepartmentRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Description string `json:"description"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new department
// @Tags auth--------------------------------- Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateDepartmentRequest true "department payload"
// @Success 201 {object} response.APIResponse
// @Router /auth/departments [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &Department{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.Create(ctx.Request.Context(), d); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create department", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "department created successfully", d)
}

// List godoc
// @Summary List all departments
// @Tags auth--------------------------------- Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /auth/departments [get]
func (h *Handler) List(ctx *gin.Context) {
	departments, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch departments", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "departments fetched successfully", departments)
}

// Get godoc
// @Summary Get department by ID
// @Tags auth--------------------------------- Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Success 200 {object} response.APIResponse
// @Router /auth/departments/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	department, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "department not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "department fetched successfully", department)
}

// Update godoc
// @Summary Update department
// @Tags auth--------------------------------- Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Param payload body UpdateDepartmentRequest true "department payload"
// @Success 200 {object} response.APIResponse
// @Router /auth/departments/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &Department{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.Update(ctx.Request.Context(), d); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update department", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "department updated successfully", d)
}

// Delete godoc
// @Summary Delete department
// @Tags auth--------------------------------- Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Success 200 {object} response.APIResponse
// @Router /auth/departments/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete department", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "department deleted successfully", nil)
}
