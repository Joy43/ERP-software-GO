package item_type

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateItemTypeRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

type UpdateItemTypeRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new item type
// @Tags iteam-profile--------------------- ItemTypes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateItemTypeRequest true "item type payload"
// @Success 201 {object} response.APIResponse
// @Router /profile-items/item-types [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateItemTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &ItemType{
		Name: req.Name,
	}

	if err := h.service.Create(ctx.Request.Context(), d); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create item type", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "item type created successfully", d)
}

// List godoc
// @Summary List all item types
// @Tags iteam-profile--------------------- ItemTypes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /profile-items/item-types [get]
func (h *Handler) List(ctx *gin.Context) {
	itemTypes, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch item types", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "item types fetched successfully", itemTypes)
}

// Get godoc
// @Summary Get item type by ID
// @Tags iteam-profile--------------------- ItemTypes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Item Type ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/item-types/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	itemType, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "item type not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "item type fetched successfully", itemType)
}

// Update godoc
// @Summary Update item type
// @Tags iteam-profile--------------------- ItemTypes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Item Type ID"
// @Param payload body UpdateItemTypeRequest true "item type payload"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/item-types/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateItemTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &ItemType{
		ID:   uint(id),
		Name: req.Name,
	}

	if err := h.service.Update(ctx.Request.Context(), d); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update item type", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "item type updated successfully", d)
}

// Delete godoc
// @Summary Delete item type
// @Tags iteam-profile--------------------- ItemTypes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Item Type ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/item-types/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete item type", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "item type deleted successfully", nil)
}
