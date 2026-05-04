package price_type

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type CreatePriceTypeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type UpdatePriceTypeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// Create godoc
// @Summary Create a new price type
// @Tags Price Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreatePriceTypeRequest true "price type payload"
// @Success 201 {object} response.APIResponse
// @Router /partners/price-types [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreatePriceTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	pt := &PriceType{
		Name:        req.Name,
		Description: req.Description,
	}

	if req.IsActive != nil {
		pt.IsActive = *req.IsActive
	} else {
		pt.IsActive = true
	}

	if err := h.service.Create(ctx.Request.Context(), pt); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create price type", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "price type created successfully", pt)
}

// List godoc
// @Summary List all price types
// @Tags Price Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /partners/price-types [get]
func (h *Handler) List(ctx *gin.Context) {
	pts, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch price types", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "price types fetched successfully", pts)
}

// Get godoc
// @Summary Get price type by ID
// @Tags Price Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Price Type ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/price-types/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	pt, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "price type not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "price type fetched successfully", pt)
}

// Update godoc
// @Summary Update price type
// @Tags Price Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Price Type ID"
// @Param payload body UpdatePriceTypeRequest true "price type payload"
// @Success 200 {object} response.APIResponse
// @Router /partners/price-types/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdatePriceTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	pt, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "price type not found", "NOT_FOUND", err.Error())
		return
	}

	pt.Name = req.Name
	pt.Description = req.Description
	if req.IsActive != nil {
		pt.IsActive = *req.IsActive
	}

	if err := h.service.Update(ctx.Request.Context(), pt); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update price type", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "price type updated successfully", pt)
}

// Delete godoc
// @Summary Delete price type
// @Tags Price Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Price Type ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/price-types/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete price type", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "price type deleted successfully", nil)
}
