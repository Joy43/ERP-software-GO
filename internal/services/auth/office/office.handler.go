package office

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateOfficeRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Location string `json:"location"`
}

type UpdateOfficeRequest struct {
	Name     string `json:"name" binding:"omitempty,min=2,max=100"`
	Location string `json:"location"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new office
// @Tags  auth--------------------------------- Offices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateOfficeRequest true "office payload"
// @Success 201 {object} response.APIResponse
// @Router /auth/offices [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateOfficeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	o := &Office{
		Name:     req.Name,
		Location: req.Location,
	}

	if err := h.service.Create(ctx.Request.Context(), o); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create office", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "office created successfully", o)
}

// List godoc
// @Summary List all offices
// @Tags  auth--------------------------------- Offices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /auth/offices [get]
func (h *Handler) List(ctx *gin.Context) {
	offices, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch offices", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "offices fetched successfully", offices)
}

// Get godoc
// @Summary Get office by ID
// @Tags  auth--------------------------------- Offices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Office ID"
// @Success 200 {object} response.APIResponse
// @Router /auth/offices/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	office, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "office not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "office fetched successfully", office)
}

// Update godoc
// @Summary Update office
// @Tags  auth--------------------------------- Offices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Office ID"
// @Param payload body UpdateOfficeRequest true "office payload"
// @Success 200 {object} response.APIResponse
// @Router /auth/offices/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateOfficeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	o := &Office{
		ID:       uint(id),
		Name:     req.Name,
		Location: req.Location,
	}

	if err := h.service.Update(ctx.Request.Context(), o); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update office", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "office updated successfully", o)
}

// Delete godoc
// @Summary Delete office
// @Tags  auth--------------------------------- Offices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Office ID"
// @Success 200 {object} response.APIResponse
// @Router /auth/offices/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete office", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "office deleted successfully", nil)
}
