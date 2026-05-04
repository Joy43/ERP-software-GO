package thana

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateThanaRequest struct {
	Name       string `json:"name" binding:"required,min=2,max=150"`
	DistrictID uint   `json:"district_id" binding:"required"`
}

type UpdateThanaRequest struct {
	Name       string `json:"name" binding:"required,min=2,max=150"`
	DistrictID uint   `json:"district_id" binding:"required"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new thana
// @Tags Thanas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateThanaRequest true "thana payload"
// @Success 201 {object} response.APIResponse
// @Router /partners/thanas [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateThanaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	t := &Thana{
		Name:       req.Name,
		DistrictID: req.DistrictID,
	}

	if err := h.service.Create(ctx.Request.Context(), t); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create thana", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "thana created successfully", t)
}

// List godoc
// @Summary List all thanas
// @Tags Thanas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param district_id query int false "District ID to filter"
// @Success 200 {object} response.APIResponse
// @Router /partners/thanas [get]
func (h *Handler) List(ctx *gin.Context) {
	districtIDStr := ctx.Query("district_id")
	if districtIDStr != "" {
		districtID, _ := strconv.Atoi(districtIDStr)
		thanas, err := h.service.FindByDistrictID(ctx.Request.Context(), uint(districtID))
		if err != nil {
			response.Error(ctx, http.StatusInternalServerError, "failed to fetch thanas", "INTERNAL_SERVER_ERROR", err.Error())
			return
		}
		response.Success(ctx, "thanas fetched successfully", thanas)
		return
	}

	thanas, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch thanas", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "thanas fetched successfully", thanas)
}

// Get godoc
// @Summary Get thana by ID
// @Tags Thanas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Thana ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/thanas/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	thana, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "thana not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "thana fetched successfully", thana)
}

// Update godoc
// @Summary Update thana
// @Tags Thanas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Thana ID"
// @Param payload body UpdateThanaRequest true "thana payload"
// @Success 200 {object} response.APIResponse
// @Router /partners/thanas/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateThanaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	t := &Thana{
		ID:         uint(id),
		Name:       req.Name,
		DistrictID: req.DistrictID,
	}

	if err := h.service.Update(ctx.Request.Context(), t); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update thana", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "thana updated successfully", t)
}

// Delete godoc
// @Summary Delete thana
// @Tags Thanas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Thana ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/thanas/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete thana", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "thana deleted successfully", nil)
}
