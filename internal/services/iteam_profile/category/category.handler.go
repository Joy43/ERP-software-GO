package category

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateCategoryRequest struct {
	Name         string `json:"name" binding:"required,min=2,max=150"`
	DepartmentID *uint  `json:"department_id"`
}

type UpdateCategoryRequest struct {
	Name         string `json:"name" binding:"required,min=2,max=150"`
	DepartmentID *uint  `json:"department_id"`
}

// MinorCategoryDTO for hierarchical response
type MinorCategoryDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// SubCategoryDTO for hierarchical response with minor categories
type SubCategoryDTO struct {
	ID             uint                `json:"id"`
	Name           string              `json:"name"`
	MinorCategories []MinorCategoryDTO `json:"minor_categories,omitempty"`
	CreatedAt      string              `json:"created_at"`
	UpdatedAt      string              `json:"updated_at"`
}

// CategoryHierarchyResponse for full hierarchy
type CategoryHierarchyResponse struct {
	ID             uint              `json:"id"`
	Name           string            `json:"name"`
	DepartmentID   *uint             `json:"department_id,omitempty"`
	Department     interface{}       `json:"department,omitempty"`
	SubCategories  []SubCategoryDTO  `json:"sub_categories,omitempty"`
	CreatedAt      string            `json:"created_at"`
	UpdatedAt      string            `json:"updated_at"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new category
// @Tags iteam-profile--------------------- Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateCategoryRequest true "category payload"
// @Success 201 {object} response.APIResponse
// @Router /profile-items/categories [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &Category{
		Name:         req.Name,
		DepartmentID: req.DepartmentID,
	}

	if err := h.service.Create(ctx.Request.Context(), d); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create category", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "category created successfully", d)
}

// List godoc
// @Summary List all categories with full hierarchy (Department>Category>SubCategory>MinorCategory)
// @Tags iteam-profile--------------------- Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /profile-items/categories [get]
func (h *Handler) List(ctx *gin.Context) {
	categories, err := h.service.FindAllWithHierarchy(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch categories", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "categories fetched successfully with full hierarchy", categories)
}

// Get godoc
// @Summary Get category by ID with full hierarchy (Department>Category>SubCategory>MinorCategory)
// @Tags iteam-profile--------------------- Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/categories/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	hierarchy, err := h.service.FindByIDWithHierarchy(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "category not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "category fetched successfully with full hierarchy", hierarchy)
}

// Update godoc
// @Summary Update category
// @Tags iteam-profile--------------------- Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param payload body UpdateCategoryRequest true "category payload"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/categories/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &Category{
		ID:           uint(id),
		Name:         req.Name,
		DepartmentID: req.DepartmentID,
	}

	if err := h.service.Update(ctx.Request.Context(), d); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update category", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "category updated successfully", d)
}

// Delete godoc
// @Summary Delete category
// @Tags iteam-profile--------------------- Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/categories/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete category", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "category deleted successfully", nil)
}
