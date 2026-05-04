package sales_tax_setup

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateSalesTaxSetupRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

type UpdateSalesTaxSetupRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new sales tax setup
// @Tags iteam-profile--------------------- SalesTaxSetups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateSalesTaxSetupRequest true "sales tax setup payload"
// @Success 201 {object} response.APIResponse
// @Router /profile-items/sales-tax-setups [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateSalesTaxSetupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &SalesTaxSetup{
		Name: req.Name,
	}

	if err := h.service.Create(ctx.Request.Context(), d); err != nil {
		// Check if it's a validation error
		if err.Error() == "sales tax setup cannot be nil" || err.Error() == "sales tax setup name is required" {
			response.Error(ctx, http.StatusBadRequest, "validation failed", "VALIDATION_ERROR", err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "failed to create sales tax setup", "INTERNAL_SERVER_ERROR", err.Error())
		}
		return
	}

	ctx.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "sales tax setup created successfully",
		Data:    d,
	})
}

// List godoc
// @Summary List all sales tax setups
// @Tags iteam-profile--------------------- SalesTaxSetups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sales-tax-setups [get]
func (h *Handler) List(ctx *gin.Context) {
	salesTaxSetups, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch sales tax setups", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	if salesTaxSetups == nil {
		salesTaxSetups = []SalesTaxSetup{}
	}

	response.Success(ctx,
		
		 "sales tax setups fetched successfully",
		  salesTaxSetups,
		)
}

// Get godoc
// @Summary Get sales tax setup by ID
// @Tags iteam-profile--------------------- SalesTaxSetups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sales Tax Setup ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sales-tax-setups/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	salesTaxSetup, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "invalid sales tax setup id" {
			response.Error(ctx, http.StatusBadRequest, "invalid id", "INVALID_ID", err.Error())
		} else {
			response.Error(ctx, http.StatusNotFound, "sales tax setup not found", "NOT_FOUND", err.Error())
		}
		return
	}

	response.Success(ctx, "sales tax setup fetched successfully", salesTaxSetup)
}

// Update godoc
// @Summary Update sales tax setup
// @Tags iteam-profile--------------------- SalesTaxSetups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sales Tax Setup ID"
// @Param payload body UpdateSalesTaxSetupRequest true "sales tax setup payload"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sales-tax-setups/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	var req UpdateSalesTaxSetupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &SalesTaxSetup{
		ID:   uint(id),
		Name: req.Name,
	}

	if err := h.service.Update(ctx.Request.Context(), d); err != nil {
		if err.Error() == "sales tax setup cannot be nil" || 
		   err.Error() == "sales tax setup id is required for update" || 
		   err.Error() == "sales tax setup name is required" {
			response.Error(ctx, http.StatusBadRequest, "validation failed", "VALIDATION_ERROR", err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "failed to update sales tax setup", "INTERNAL_SERVER_ERROR", err.Error())
		}
		return
	}

	response.Success(ctx, "sales tax setup updated successfully", d)
}

// Delete godoc
// @Summary Delete sales tax setup
// @Tags iteam-profile--------------------- SalesTaxSetups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sales Tax Setup ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sales-tax-setups/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		// ----------Check if it's a validation error------------
		if err.Error() == "invalid sales tax setup id" {
			response.Error(ctx, http.StatusBadRequest, "invalid id", "INVALID_ID", err.Error())
		} else {
			response.Error(ctx, http.StatusNotFound, "sales tax setup not found", "NOT_FOUND", err.Error())
		}
		return
	}

	response.Success(ctx, "sales tax setup deleted successfully", nil)
}
