package sales_setup

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateSalesSetupRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

type UpdateSalesSetupRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// respondWithServiceError handles service layer errors and returns appropriate HTTP response
func (h *Handler) respondWithServiceError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	errMsg := err.Error()

	// Map service error messages to HTTP status codes and error codes
	// Validation errors (400)
	if errMsg == "sales setup cannot be nil" ||
		errMsg == "sales setup name is required" ||
		errMsg == "sales setup id is required for update" ||
		errMsg == "invalid sales setup id" ||
		errMsg == "sales setup name must not exceed 150 characters" ||
		errMsg == "sales setup name must have at least 2 characters" {
		response.Error(ctx, http.StatusBadRequest, "validation failed", "VALIDATION_ERROR", errMsg)
		return
	}

	// Not found errors (404)
	if contains(errMsg, "not found") {
		response.Error(ctx, http.StatusNotFound, "resource not found", "NOT_FOUND", errMsg)
		return
	}

	// Default to internal server error (500)
	response.Error(ctx, http.StatusInternalServerError, "internal server error", "INTERNAL_SERVER_ERROR", errMsg)
}

// contains is a helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Create godoc
// @Summary Create a new sales setup
// @Tags iteam-profile--------------------- SalesSetups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateSalesSetupRequest true "sales setup payload"
// @Success 201 {object} response.APIResponse
// @Router /profile-items/sales-setups [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateSalesSetupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &SalesSetup{
		Name: req.Name,
	}

	if err := h.service.Create(ctx.Request.Context(), d); err != nil {
		h.respondWithServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, response.APIResponse{
		Message: "sales setup created successfully",
		Data:    d,
		Success: true,
	})
}

// List godoc
// @Summary List all sales setups
// @Tags iteam-profile--------------------- SalesSetups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sales-setups [get]
func (h *Handler) List(ctx *gin.Context) {
	salesSetups, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		h.respondWithServiceError(ctx, err)
		return
	}

	response.Success(ctx, "sales setups fetched successfully", salesSetups)
}

// Get godoc
// @Summary Get sales setup by ID
// @Tags iteam-profile--------------------- SalesSetups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sales Setup ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sales-setups/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	salesSetup, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		h.respondWithServiceError(ctx, err)
		return
	}

	response.Success(ctx, "sales setup fetched successfully", salesSetup)
}

// Update godoc
// @Summary Update sales setup
// @Tags iteam-profile--------------------- SalesSetups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sales Setup ID"
// @Param payload body UpdateSalesSetupRequest true "sales setup payload"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sales-setups/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	var req UpdateSalesSetupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &SalesSetup{
		ID:   uint(id),
		Name: req.Name,
	}

	if err := h.service.Update(ctx.Request.Context(), d); err != nil {
		h.respondWithServiceError(ctx, err)
		return
	}

	response.Success(ctx, "sales setup updated successfully", d)
}

// Delete godoc
// @Summary Delete sales setup
// @Tags iteam-profile--------------------- SalesSetups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Sales Setup ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/sales-setups/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		h.respondWithServiceError(ctx, err)
		return
	}

	response.Success(ctx, "sales setup deleted successfully", nil)
}
