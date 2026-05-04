package district

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateDistrictRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

type UpdateDistrictRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new district
// @Tags Districts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateDistrictRequest true "district payload"
// @Success 201 {object} response.APIResponse
// @Router /partners/districts [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateDistrictRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &District{
		Name: req.Name,
	}

	if err := h.service.Create(ctx.Request.Context(), d); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create district", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "district created successfully", d)
}

// List godoc
// @Summary List all districts
// @Tags Districts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /partners/districts [get]
func (h *Handler) List(ctx *gin.Context) {
	districts, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch districts", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "districts fetched successfully", districts)
}

// Get godoc
// @Summary Get district by ID
// @Tags Districts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "District ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/districts/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	district, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "district not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "district fetched successfully", district)
}

// Update godoc
// @Summary Update district
// @Tags Districts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "District ID"
// @Param payload body UpdateDistrictRequest true "district payload"
// @Success 200 {object} response.APIResponse
// @Router /partners/districts/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateDistrictRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &District{
		ID:   uint(id),
		Name: req.Name,
	}

	if err := h.service.Update(ctx.Request.Context(), d); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update district", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "district updated successfully", d)
}

// Delete godoc
// @Summary Delete district
// @Tags Districts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "District ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/districts/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete district", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "district deleted successfully", nil)
}
