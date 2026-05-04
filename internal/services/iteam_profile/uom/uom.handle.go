package umo_measurement

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreatePackageUom struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

type UpdatePackageUom struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new uom
// @Tags iteam-profile---------------------Uom
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreatePackageUom true "package uom payload"
// @Success 201 {object} response.APIResponse
// @Router /profile-items/uom [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreatePackageUom
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &Uom{
		Name: req.Name,
	}

	if err := h.service.Create(ctx.Request.Context(), d); err != nil {
		if err.Error() == "uom base cannot be nil" || err.Error() == "uom base name is required" {
			response.Error(ctx, http.StatusBadRequest, "validation failed", "VALIDATION_ERROR", err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "failed to create uom base", "INTERNAL_SERVER_ERROR", err.Error())
		}
		return
	}

	ctx.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "uom created successfully",
		Data:    d,
	})
}

// List godoc
// @Summary List all package uom
// @Tags iteam-profile---------------------Uom
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /profile-items/uom [get]
func (h *Handler) List(ctx *gin.Context) {
	salesSetups, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch uoms", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}
	if salesSetups == nil {
		salesSetups = []Uom{}
	}

	response.Success(ctx, "uoms fetched successfully", salesSetups)
}

// Get godoc
// @Summary Get uom by ID
// @Tags iteam-profile---------------------Uom
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "UOM ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/uom/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	salesSetup, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "invalid uom id" {
			response.Error(ctx, http.StatusBadRequest, "invalid id", "INVALID_ID", err.Error())
		} else if err.Error() == "record not found" || err.Error() == "sql: no rows in result set" {
			response.Error(ctx, http.StatusNotFound, "uom not found", "NOT_FOUND", err.Error())
		} else {
			response.Error(ctx, http.StatusNotFound, "uom not found", "NOT_FOUND", err.Error())
		}
		return
	}

	response.Success(ctx, "uom fetched successfully", salesSetup)
}

// Update godoc
// @Summary Update uom
// @Tags iteam-profile---------------------Uom
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "uom ID"
// @Param payload body UpdatePackageUom true "uom payload"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/uom/{id} [patch]
func (h *Handler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	var req UpdatePackageUom
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request payload", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &Uom{
		ID:   uint(id),
		Name: req.Name,
	}

	if err := h.service.Update(ctx.Request.Context(), d); err != nil {
	
		if err.Error() == "uom cannot be nil" || 
		   err.Error() == "uom id is required for update" || 
		   err.Error() == "uom name is required" {
			response.Error(ctx, http.StatusBadRequest, "validation failed", "VALIDATION_ERROR", err.Error())
		} else {
			response.Error(ctx, http.StatusInternalServerError, "failed to update uom", "INTERNAL_SERVER_ERROR", err.Error())
		}
		return
	}

	response.Success(ctx, "uom updated successfully", d)
}

// Delete godoc
// @Summary Delete uom
// @Tags iteam-profile---------------------Uom
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "uom ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/uom/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.Error(ctx, http.StatusBadRequest, "invalid id format", "INVALID_ID", "id must be a positive integer")
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete uom", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "uom deleted successfully", nil)
}
