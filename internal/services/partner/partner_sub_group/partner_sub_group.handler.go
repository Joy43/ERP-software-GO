package partner_sub_group

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreatePartnerSubGroupRequest struct {
	Name           string `json:"name" binding:"required,min=2,max=150"`
	PartnerGroupID uint   `json:"partner_group_id" binding:"required"`
}

type UpdatePartnerSubGroupRequest struct {
	Name           string `json:"name" binding:"required,min=2,max=150"`
	PartnerGroupID uint   `json:"partner_group_id" binding:"required"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new partner sub group
// @Tags Partner Sub Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreatePartnerSubGroupRequest true "partner sub group payload"
// @Success 201 {object} response.APIResponse
// @Router /partners/partner-sub-groups [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreatePartnerSubGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	psg := &PartnerSubGroup{
		Name:           req.Name,
		PartnerGroupID: req.PartnerGroupID,
	}

	if err := h.service.Create(ctx.Request.Context(), psg); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create partner sub group", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "partner sub group created successfully", psg)
}

// List godoc
// @Summary List all partner sub groups
// @Tags Partner Sub Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param partner_group_id query int false "Partner Group ID to filter"
// @Success 200 {object} response.APIResponse
// @Router /partners/partner-sub-groups [get]
func (h *Handler) List(ctx *gin.Context) {
	groupIDStr := ctx.Query("partner_group_id")
	if groupIDStr != "" {
		groupID, _ := strconv.Atoi(groupIDStr)
		subGroups, err := h.service.FindByGroupID(ctx.Request.Context(), uint(groupID))
		if err != nil {
			response.Error(ctx, http.StatusInternalServerError, "failed to fetch partner sub groups", "INTERNAL_SERVER_ERROR", err.Error())
			return
		}
		response.Success(ctx, "partner sub groups fetched successfully", subGroups)
		return
	}

	subGroups, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch partner sub groups", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "partner sub groups fetched successfully", subGroups)
}

// Get godoc
// @Summary Get partner sub group by ID
// @Tags Partner Sub Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Partner Sub Group ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/partner-sub-groups/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	subGroup, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "partner sub group not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "partner sub group fetched successfully", subGroup)
}

// Update godoc
// @Summary Update partner sub group
// @Tags Partner Sub Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Partner Sub Group ID"
// @Param payload body UpdatePartnerSubGroupRequest true "partner sub group payload"
// @Success 200 {object} response.APIResponse
// @Router /partners/partner-sub-groups/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdatePartnerSubGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	psg := &PartnerSubGroup{
		ID:             uint(id),
		Name:           req.Name,
		PartnerGroupID: req.PartnerGroupID,
	}

	if err := h.service.Update(ctx.Request.Context(), psg); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update partner sub group", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "partner sub group updated successfully", psg)
}

// Delete godoc
// @Summary Delete partner sub group
// @Tags Partner Sub Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Partner Sub Group ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/partner-sub-groups/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete partner sub group", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "partner sub group deleted successfully", nil)
}
