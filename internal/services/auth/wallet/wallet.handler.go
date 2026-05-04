package wallet

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create wallet
// @Description Create a new wallet
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateWalletRequest true "Create wallet request"
// @Success 201 {object} response.APIResponse
// @Router /auth/wallets [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "INVALID_REQUEST", err.Error())
		return
	}

	result, err := h.service.Create(ctx.Request.Context(), req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create wallet", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "wallet created", result)
}

// List godoc
// @Summary List wallets
// @Description Get all wallets
// @Tags Wallet
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Router /auth/wallets [get]
func (h *Handler) List(ctx *gin.Context) {
	results, err := h.service.List(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to list wallets", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "wallets retrieved", results)
}

// Get godoc
// @Summary Get wallet
// @Description Get wallet by ID
// @Tags Wallet
// @Produce json
// @Security BearerAuth
// @Param id path int true "Wallet ID"
// @Success 200 {object} response.APIResponse
// @Router /auth/wallets/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 0, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid ID", "INVALID_ID", err.Error())
		return
	}

	result, err := h.service.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to get wallet", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "wallet retrieved", result)
}

// Update godoc
// @Summary Update wallet
// @Description Update an existing wallet
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Wallet ID"
// @Param request body UpdateWalletRequest true "Update wallet request"
// @Success 200 {object} response.APIResponse
// @Router /auth/wallets/{id} [put]
func (h *Handler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 0, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid ID", "INVALID_ID", err.Error())
		return
	}

	var req UpdateWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "INVALID_REQUEST", err.Error())
		return
	}

	result, err := h.service.Update(ctx.Request.Context(), uint(id), req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update wallet", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "wallet updated", result)
}

// Delete godoc
// @Summary Delete wallet
// @Description Delete a wallet
// @Tags Wallet
// @Produce json
// @Security BearerAuth
// @Param id path int true "Wallet ID"
// @Success 200 {object} response.APIResponse
// @Router /auth/wallets/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 0, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid ID", "INVALID_ID", err.Error())
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete wallet", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "wallet deleted", nil)
}
