package tax_bracket

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateTaxBracketRequest struct {
	Name       string  `json:"name" binding:"required,min=2,max=150"`
	Percentage float64 `json:"percentage"`
}

type UpdateTaxBracketRequest struct {
	Name       string  `json:"name" binding:"required,min=2,max=150"`
	Percentage float64 `json:"percentage"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new tax bracket
// @Tags Tax Brackets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateTaxBracketRequest true "tax bracket payload"
// @Success 201 {object} response.APIResponse
// @Router /partners/tax-brackets [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateTaxBracketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	tb := &TaxBracket{
		Name:       req.Name,
		Percentage: req.Percentage,
	}

	if err := h.service.Create(ctx.Request.Context(), tb); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create tax bracket", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "tax bracket created successfully", tb)
}

// List godoc
// @Summary List all tax brackets
// @Tags Tax Brackets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /partners/tax-brackets [get]
func (h *Handler) List(ctx *gin.Context) {
	brackets, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch tax brackets", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "tax brackets fetched successfully", brackets)
}

// Get godoc
// @Summary Get tax bracket by ID
// @Tags Tax Brackets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tax Bracket ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/tax-brackets/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	bracket, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "tax bracket not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "tax bracket fetched successfully", bracket)
}

// Update godoc
// @Summary Update tax bracket
// @Tags Tax Brackets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tax Bracket ID"
// @Param payload body UpdateTaxBracketRequest true "tax bracket payload"
// @Success 200 {object} response.APIResponse
// @Router /partners/tax-brackets/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateTaxBracketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	tb := &TaxBracket{
		ID:         uint(id),
		Name:       req.Name,
		Percentage: req.Percentage,
	}

	if err := h.service.Update(ctx.Request.Context(), tb); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update tax bracket", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "tax bracket updated successfully", tb)
}

// Delete godoc
// @Summary Delete tax bracket
// @Tags Tax Brackets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tax Bracket ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/tax-brackets/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete tax bracket", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "tax bracket deleted successfully", nil)
}
