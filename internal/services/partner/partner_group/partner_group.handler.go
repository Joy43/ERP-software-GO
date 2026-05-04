package partner_group

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

type CreatePartnerGroupRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

type UpdatePartnerGroupRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new partner group
// @Tags Partner Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreatePartnerGroupRequest true "partner group payload"
// @Success 201 {object} response.APIResponse
// @Router /partners/partner-groups [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreatePartnerGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	pg := &PartnerGroup{
		Name: req.Name,
	}

	if err := h.service.Create(ctx.Request.Context(), pg); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create partner group", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "partner group created successfully", pg)
}

// List godoc
// @Summary List all partner groups
// @Tags Partner Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /partners/partner-groups [get]
func (h *Handler) List(ctx *gin.Context) {
	groups, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch partner groups", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "partner groups fetched successfully", groups)
}

// Get godoc
// @Summary Get partner group by ID
// @Tags Partner Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Partner Group ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/partner-groups/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	group, err := h.service.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "partner group not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "partner group fetched successfully", group)
}

// Update godoc
// @Summary Update partner group
// @Tags Partner Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Partner Group ID"
// @Param payload body UpdatePartnerGroupRequest true "partner group payload"
// @Success 200 {object} response.APIResponse
// @Router /partners/partner-groups/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdatePartnerGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR", err.Error())
		return
	}

	pg := &PartnerGroup{
		ID:   uint(id),
		Name: req.Name,
	}

	if err := h.service.Update(ctx.Request.Context(), pg); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update partner group", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "partner group updated successfully", pg)
}

// Delete godoc
// @Summary Delete partner group
// @Tags Partner Groups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Partner Group ID"
// @Success 200 {object} response.APIResponse
// @Router /partners/partner-groups/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete partner group", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "partner group deleted successfully", nil)
}
