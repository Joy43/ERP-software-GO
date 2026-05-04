package sales_supply_type

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateSalesSupplyTypeRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

type UpdateSalesSupplyTypeRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new sales supply type
// @Tags iteam-profile--------------------- SalesSupplyTypes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateSalesSupplyTypeRequest true "sales supply type payload"
// @Success 201 {object} response.APIResponse
// @Router /profile-items/sales-supply-types [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateSalesSupplyTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &SalesSupplyType{
		Name: req.Name,
	}

	if err := h.service.Create(ctx.Request.Context(), d); err != nil {
		if err.Error() == "sales supply type cannot be nil" || err.Error() == "sales supply type name is required" {
			response.Error(ctx, http.StatusBadRequest, "validation failed", "VALIDATION_ERROR", err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "failed to create sales supply type", "INTERNAL_SERVER_ERROR", err.Error())
		}
		return
	}

	ctx.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "sales supply type created successfully",
		Data:    d,
	})
}

// List godoc
// @Summary List all sales supply types
// @Tags iteam-profile--------------------- SalesSupplyTypes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sales-supply-types [get]
func (h *Handler) List(ctx *gin.Context) {
	salesSupplyTypes, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch sales supply types", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	// Return empty array if no data instead of null
	if salesSupplyTypes == nil {
		salesSupplyTypes = []SalesSupplyType{}
	}

	response.Success(ctx, "sales supply types fetched successfully", salesSupplyTypes)
}

// Get godoc
// @Summary Get sales supply type by ID
// @Tags iteam-profile--------------------- SalesSupplyTypes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sales Supply Type ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sales-supply-types/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	salesSupplyType, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "invalid sales supply type id" {
			response.Error(ctx, http.StatusBadRequest, "invalid id", "INVALID_ID", err.Error())
		} else {
			response.Error(ctx, http.StatusNotFound, "sales supply type not found", "NOT_FOUND", err.Error())
		}
		return
	}

	response.Success(ctx, "sales supply type fetched successfully", salesSupplyType)
}

// Update godoc
// @Summary Update sales supply type
// @Tags iteam-profile--------------------- SalesSupplyTypes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sales Supply Type ID"
// @Param payload body UpdateSalesSupplyTypeRequest true "sales supply type payload"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sales-supply-types/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	var req UpdateSalesSupplyTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &SalesSupplyType{
		ID:   uint(id),
		Name: req.Name,
	}

	if err := h.service.Update(ctx.Request.Context(), d); err != nil {
		// Check if it's a validation error
		if err.Error() == "sales supply type cannot be nil" || 
		   err.Error() == "sales supply type id is required for update" || 
		   err.Error() == "sales supply type name is required" {
			response.Error(ctx, http.StatusBadRequest, "validation failed", "VALIDATION_ERROR", err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "failed to update sales supply type", "INTERNAL_SERVER_ERROR", err.Error())
		}
		return
	}

	response.Success(ctx, "sales supply type updated successfully", d)
}

// Delete godoc
// @Summary Delete sales supply type
// @Tags iteam-profile--------------------- SalesSupplyTypes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sales Supply Type ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sales-supply-types/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		// Check if it's a validation error
		if err.Error() == "invalid sales supply type id" {
			response.Error(ctx, http.StatusBadRequest, "invalid id", "INVALID_ID", err.Error())
		} else {
			response.Error(ctx, http.StatusNotFound, "sales supply type not found", "NOT_FOUND", err.Error())
		}
		return
	}

	response.Success(ctx, "sales supply type deleted successfully", nil)
}
