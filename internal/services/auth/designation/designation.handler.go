package designation

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateDesignationRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateDesignationRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Description string `json:"description"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new designation
// @Tags auth--------------------------------- Designations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateDesignationRequest true "designation payload"
// @Success 201 {object} response.APIResponse
// @Router /auth/designations [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateDesignationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &Designation{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.Create(ctx.Request.Context(), d); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create designation", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "designation created successfully", d)
}

// List godoc
// @Summary List all designations
// @Tags auth--------------------------------- Designations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /auth/designations [get]
func (h *Handler) List(ctx *gin.Context) {
	designations, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch designations", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "designations fetched successfully", designations)
}

// Get godoc
// @Summary Get designation by ID
// @Tags auth--------------------------------- Designations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Designation ID"
// @Success 200 {object} response.APIResponse
// @Router /auth/designations/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	designation, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "designation not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "designation fetched successfully", designation)
}

// Update godoc
// @Summary Update designation
// @Tags auth--------------------------------- Designations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Designation ID"
// @Param payload body UpdateDesignationRequest true "designation payload"
// @Success 200 {object} response.APIResponse
// @Router /auth/designations/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateDesignationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &Designation{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.Update(ctx.Request.Context(), d); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update designation", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "designation updated successfully", d)
}

// Delete godoc
// @Summary Delete designation
// @Tags auth--------------------------------- Designations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Designation ID"
// @Success 200 {object} response.APIResponse
// @Router /auth/designations/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete designation", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "designation deleted successfully", nil)
}
