package sales_representative

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

type CreateSalesRepresentativeRequest struct {
	Name           string  `json:"name" binding:"required"`
	Mobile         string  `json:"mobile" binding:"required"`
	Email          string  `json:"email"`
	CommissionRate float64 `json:"commission_rate"`
	IsActive       *bool   `json:"is_active"`
}

type UpdateSalesRepresentativeRequest struct {
	Name           string  `json:"name" binding:"required"`
	Mobile         string  `json:"mobile" binding:"required"`
	Email          string  `json:"email"`
	CommissionRate float64 `json:"commission_rate"`
	IsActive       *bool   `json:"is_active"`
}

// Create godoc
// @Summary Create a new sales representative
// @Tags Sales Representatives
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateSalesRepresentativeRequest true "sales representative payload"
// @Success 201 {object} response.APIResponse
// @Router /partners/sales-representatives [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateSalesRepresentativeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	sr := &SalesRepresentative{
		Name:           req.Name,
		Mobile:         req.Mobile,
		Email:          req.Email,
		CommissionRate: req.CommissionRate,
	}

	if req.IsActive != nil {
		sr.IsActive = *req.IsActive
	} else {
		sr.IsActive = true
	}

	if err := h.service.Create(ctx.Request.Context(), sr); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create sales representative", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "sales representative created successfully", sr)
}

// List godoc
// @Summary List all sales representatives
// @Tags Sales Representatives
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /partners/sales-representatives [get]
func (h *Handler) List(ctx *gin.Context) {
	srs, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch sales representatives", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "sales representatives fetched successfully", srs)
}

// Get godoc
// @Summary Get sales representative by ID
// @Tags Sales Representatives
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sales Representative ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/sales-representatives/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	sr, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "sales representative not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "sales representative fetched successfully", sr)
}

// Update godoc
// @Summary Update sales representative
// @Tags Sales Representatives
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sales Representative ID"
// @Param payload body UpdateSalesRepresentativeRequest true "sales representative payload"
// @Success 200 {object} response.APIResponse
// @Router /partners/sales-representatives/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateSalesRepresentativeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	sr, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "sales representative not found", "NOT_FOUND", err.Error())
		return
	}

	sr.Name = req.Name
	sr.Mobile = req.Mobile
	sr.Email = req.Email
	sr.CommissionRate = req.CommissionRate
	if req.IsActive != nil {
		sr.IsActive = *req.IsActive
	}

	if err := h.service.Update(ctx.Request.Context(), sr); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update sales representative", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "sales representative updated successfully", sr)
}

// Delete godoc
// @Summary Delete sales representative
// @Tags Sales Representatives
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sales Representative ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/sales-representatives/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete sales representative", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "sales representative deleted successfully", nil)
}
