package minor_category

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateMinorCategoryRequest struct {
	SubCategoryID uint   `json:"sub_category_id" binding:"required"`
	Name          string `json:"name" binding:"required,min=2,max=150"`
}

type UpdateMinorCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new minor category
// @Tags iteam-profile--------------------- MinorCategories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateMinorCategoryRequest true "minor category payload"
// @Success 201 {object} response.APIResponse
// @Router /profile-items/minor-categories [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateMinorCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &MinorCategory{
		SubCategoryID: req.SubCategoryID,
		Name:          req.Name,
	}

	if err := h.service.Create(ctx.Request.Context(), d); err != nil {
		// ------Check if it's a validation error------
		if err.Error() == "minor category cannot be nil" || err.Error() == "minor category name is required" {
			response.Error(ctx, http.StatusBadRequest, "validation failed", "VALIDATION_ERROR", err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "failed to create minor category", "INTERNAL_SERVER_ERROR", err.Error())
		}
		return
	}

	ctx.JSON(http.StatusCreated, response.APIResponse{
		Message: "minor category created successfully",
		Data:    d,
	})
}

// List godoc
// @Summary List all minor categories
// @Tags iteam-profile--------------------- MinorCategories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /profile-items/minor-categories [get]
func (h *Handler) List(ctx *gin.Context) {
	minorCategories, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch minor categories", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	// ------Check if data is nil and Return empty array if no data instead of null
	if minorCategories == nil {
		minorCategories = []MinorCategory{}
	}

	response.Success(ctx, "minor categories fetched successfully", minorCategories)
}

// Get godoc
// @Summary Get minor category by ID
// @Tags iteam-profile--------------------- MinorCategories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Minor Category ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/minor-categories/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	minorCategory, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		// ------Check if it's a not found error------
		if err.Error() == "invalid minor category id" {
			response.Error(ctx, http.StatusBadRequest, "invalid id", "INVALID_ID", err.Error())
		} else if err.Error() == "record not found" || err.Error() == "sql: no rows in result set" {
			response.Error(ctx, http.StatusNotFound, "minor category not found", "NOT_FOUND", err.Error())
		} else {
			response.Error(ctx, http.StatusNotFound, "minor category not found", "NOT_FOUND", err.Error())
		}
		return
	}

	response.Success(ctx, "minor category fetched successfully", minorCategory)
}

// Update godoc
// @Summary Update minor category
// @Tags iteam-profile--------------------- MinorCategories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Minor Category ID"
// @Param payload body UpdateMinorCategoryRequest true "minor category payload"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/minor-categories/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	var req UpdateMinorCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &MinorCategory{
		ID:   uint(id),
		Name: req.Name,
	}

	if err := h.service.Update(ctx.Request.Context(), d); err != nil {
		// ------Check if it's a validation error------
		if err.Error() == "minor category cannot be nil" || 
		   err.Error() == "minor category id is required for update" || 
		   err.Error() == "minor category name is required" {
			response.Error(ctx, http.StatusBadRequest, "validation failed", "VALIDATION_ERROR", err.Error())
		} else if err.Error() == "record not found" || err.Error() == "sql: no rows in result set" {
			response.Error(ctx, http.StatusNotFound, "minor category not found", "NOT_FOUND", err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "failed to update minor category", "INTERNAL_SERVER_ERROR", err.Error())
		}
		return
	}

	response.Success(ctx, "minor category updated successfully", d)
}

// Delete godoc
// @Summary Delete minor category
// @Tags iteam-profile--------------------- MinorCategories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Minor Category ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/minor-categories/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
	
		if err.Error() == "invalid minor category id" {
			response.Error(ctx, http.StatusBadRequest, "invalid id", "INVALID_ID", err.Error())
		} else if err.Error() == "record not found" || err.Error() == "sql: no rows in result set" {
			response.Error(ctx, http.StatusNotFound, "minor category not found", "NOT_FOUND", err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "failed to delete minor category", "INTERNAL_SERVER_ERROR", err.Error())
		}
		return
	}

	response.Success(ctx, "minor category deleted successfully", nil)
}
