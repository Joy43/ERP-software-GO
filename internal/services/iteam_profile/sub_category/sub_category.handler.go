package sub_category

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateSubCategoryRequest struct {
	CategoryID uint   `json:"category_id" binding:"required"`
	Name       string `json:"name" binding:"required,min=2,max=150"`
}

type UpdateSubCategoryRequest struct {
	CategoryID uint   `json:"category_id" binding:"required"`
	Name       string `json:"name" binding:"required,min=2,max=150"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new sub-category
// @Tags iteam-profile -----------------------SubCategories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateSubCategoryRequest true "sub-category payload"
// @Success 201 {object} response.APIResponse
// @Router /profile-items/sub-categories [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateSubCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &SubCategory{
		CategoryID: req.CategoryID,
		Name:       req.Name,
	}

	if err := h.service.Create(ctx.Request.Context(), d); err != nil {
		// -------- Check if it's a validation error ----------
		if err.Error() == "sub-category cannot be nil" || 
		   err.Error() == "sub-category name is required" || 
		   err.Error() == "category id is required" {
			response.Error(ctx, http.StatusBadRequest, "validation failed", "VALIDATION_ERROR", err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "failed to create sub-category", "INTERNAL_SERVER_ERROR", err.Error())
		}
		return
	}

	response.Success(ctx, "sub-category created successfully", d)
}

// List godoc
// @Summary List all sub-categories
// @Tags iteam-profile -----------------------SubCategories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sub-categories [get]
func (h *Handler) List(ctx *gin.Context) {
	subCategories, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch sub-categories", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	// Return empty array if no data instead of null
	if subCategories == nil {
		subCategories = []SubCategory{}
	}

	response.Success(ctx, "sub-categories fetched successfully", subCategories)
}

// Get godoc
// @Summary Get sub-category by ID
// @Tags iteam-profile -----------------------SubCategories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "SubCategory ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sub-categories/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	subCategory, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "invalid sub-category id" {
			response.Error(ctx, http.StatusBadRequest, "invalid id", "INVALID_ID", err.Error())
		} else {
			response.Error(ctx, http.StatusNotFound, "sub-category not found", "NOT_FOUND", err.Error())
		}
		return
	}

	response.Success(ctx, "sub-category fetched successfully", subCategory)
}

// Update godoc
// @Summary Update sub-category
// @Tags iteam-profile -----------------------SubCategories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "SubCategory ID"
// @Param payload body UpdateSubCategoryRequest true "sub-category payload"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sub-categories/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	var req UpdateSubCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &SubCategory{
		ID:         uint(id),
		CategoryID: req.CategoryID,
		Name:       req.Name,
	}

	if err := h.service.Update(ctx.Request.Context(), d); err != nil {
		// Check if it's a validation error
		if err.Error() == "sub-category cannot be nil" || 
		   err.Error() == "sub-category id is required for update" || 
		   err.Error() == "sub-category name is required" || 
		   err.Error() == "category id is required" {
			response.Error(ctx, http.StatusBadRequest, "validation failed", "VALIDATION_ERROR", err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "failed to update sub-category", "INTERNAL_SERVER_ERROR", err.Error())
		}
		return
	}

	response.Success(ctx, "sub-category updated successfully", d)
}

// Delete godoc
// @Summary Delete sub-category
// @Tags iteam-profile -----------------------SubCategories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "SubCategory ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sub-categories/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		// Check if it's a validation error
		if err.Error() == "invalid sub-category id" {
			response.Error(ctx, http.StatusBadRequest, "invalid id", "INVALID_ID", err.Error())
		} else {
			response.Error(ctx, http.StatusNotFound, "sub-category not found", "NOT_FOUND", err.Error())
		}
		return
		return
	}

	response.Success(ctx, "sub-category deleted successfully", nil)
}
