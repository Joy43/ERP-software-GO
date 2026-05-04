package role

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

// Create godoc
// @Summary Create a new role
// @Description Create a new role with optional permissions
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateRoleRequest true "role payload"
// @Success 201 {object} response.APIResponse
// @Router /auth/roles [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "BAD_REQUEST", err.Error())
		return
	}

	role, err := h.service.CreateRole(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create role", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "role created successfully", role)
}

// List godoc
// @Summary List all roles
// @Description Get a list of all roles
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /auth/roles [get]
func (h *Handler) List(ctx *gin.Context) {
	roles, err := h.service.ListRoles(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch roles", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "roles fetched successfully", roles)
}

// Get godoc
// @Summary Get role by ID
// @Description Get details of a role by its ID
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} response.APIResponse
// @Router /auth/roles/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	role, err := h.service.GetRole(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "role not found", "NOT_FOUND", err.Error())
		return
	}

	response.Success(ctx, "role fetched successfully", role)
}

// Update godoc
// @Summary Update role
// @Description Update role details and permissions
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param payload body UpdateRoleRequest true "role payload"
// @Success 200 {object} response.APIResponse
// @Router /auth/roles/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "BAD_REQUEST", err.Error())
		return
	}

	role, err := h.service.UpdateRole(ctx.Request.Context(), uint(id), req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update role", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "role updated successfully", role)
}

// Delete godoc
// @Summary Delete role
// @Description Delete a role by its ID
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 204 {object} response.APIResponse
// @Router /auth/roles/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.service.DeleteRole(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete role", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "role deleted successfully", nil)
}
