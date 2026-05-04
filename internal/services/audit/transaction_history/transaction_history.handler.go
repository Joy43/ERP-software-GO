package transaction_history

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
// @Summary Create transaction delete history
// @Description Create a new record of a deleted transaction
// @Tags Transaction History
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateDeleteHistoryRequest true "Create history request"
// @Success 201 {object} response.APIResponse
// @Router /audit/transaction-history [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateDeleteHistoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "INVALID_REQUEST", err.Error())
		return
	}

	userID, _ := ctx.Get("user_id")
	deletedBy, ok := userID.(uint)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED", "user context missing")
		return
	}

	result, err := h.service.CreateHistory(ctx.Request.Context(), deletedBy, req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create history", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "transaction delete history created", result)
}

// List godoc
// @Summary List transaction delete history
// @Description Get transition delete history with optional filters
// @Tags Transaction History
// @Produce json
// @Security BearerAuth
// @Param transaction_type query string false "Transaction Type"
// @Param deleted_by query int false "Deleted By User ID"
// @Param from_date query string false "From Date (YYYY-MM-DD)"
// @Param to_date query string false "To Date (YYYY-MM-DD)"
// @Success 200 {object} response.APIResponse
// @Router /audit/transaction-history [get]
func (h *Handler) List(ctx *gin.Context) {
	var filter HistoryFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid query parameters", "INVALID_QUERY", err.Error())
		return
	}

	results, err := h.service.ListHistory(ctx.Request.Context(), filter)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to list history", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "transaction delete history retrieved", results)
}

// Delete godoc
// @Summary Delete transaction history record
// @Description Delete a record from the transaction history log
// @Tags Transaction History
// @Produce json
// @Security BearerAuth
// @Param id path int true "History ID"
// @Success 200 {object} response.APIResponse
// @Router /audit/transaction-history/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid ID", "INVALID_ID", err.Error())
		return
	}

	if err := h.service.DeleteHistory(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete record", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "transaction history record deleted", nil)
}
