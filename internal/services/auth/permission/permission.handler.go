package permission

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List godoc
// @Summary List all permissions
// @Description Get a list of all available permissions
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /auth/permissions [get]
func (h *Handler) List(ctx *gin.Context) {
	permissions, err := h.service.ListPermissions(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch permissions", "INTERNAL_SERVER_ERROR", nil)
		return
	}

	response.Success(ctx, "permissions fetched successfully", permissions)
}
