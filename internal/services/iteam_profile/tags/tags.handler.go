package tags

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreateTagsRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

type UpdateTagsRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new tag
// @Tags iteam-profile----------------------- Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateTagsRequest true "tag payload"
// @Success 201 {object} response.APIResponse
// @Router /profile-items/tags [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateTagsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &Tags{
		Name: req.Name,
	}

	if err := h.service.Create(ctx.Request.Context(), d); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create tag", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "tag created successfully", d)
}

// List godoc
// @Summary List all tags
// @Tags iteam-profile----------------------- Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /profile-items/tags [get]
func (h *Handler) List(ctx *gin.Context) {
	tags, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch tags", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}
	response.Success(ctx, "tags fetched successfully", tags)
}

// Get godoc
// @Summary Get tag by ID
// @Tags iteam-profile----------------------- Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/tags/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	tag, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "tag not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "tag fetched successfully", tag)
}

// Update godoc
// @Summary Update tag
// @Tags iteam-profile----------------------- Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Param payload body UpdateTagsRequest true "tag payload"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/tags/{id} [patch]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateTagsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	d := &Tags{
		ID:   uint(id),
		Name: req.Name,
	}

	if err := h.service.Update(ctx.Request.Context(), d); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update tag", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "tag updated successfully", d)
}

// Delete godoc
// @Summary Delete tag
// @Tags iteam-profile----------------------- Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Tag ID"
// @Success 200 {object} response.APIResponse
// @Router /profile-items/tags/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete tag", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "tag deleted successfully", nil)
}
