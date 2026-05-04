package payment_mode

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
// @Summary List all payment modes
// @Tags Payment Modes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /auth/payment-modes [get]
func (h *Handler) List(ctx *gin.Context) {
	paymentModes, err := h.service.FindAll(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to fetch payment modes", "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(ctx, "payment modes fetched successfully", paymentModes)
}
