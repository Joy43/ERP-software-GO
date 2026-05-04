package responsibility_transfer

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/response"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type createTransferRequest struct {
	FromUserID uint   `json:"from_user_id" binding:"required" example:"1"`
	ToUserID   uint   `json:"to_user_id" binding:"required" example:"2"`
	FromDate   string `json:"from_date" binding:"required" example:"2026-04-04"`
	ToDate     string `json:"to_date" binding:"required" example:"2026-04-10"`
	Remarks    string `json:"remarks" example:"Transferring duties"`
}

type listFilterRequest struct {
	FromUserID uint   `form:"from_user_id"`
	ToUserID   uint   `form:"to_user_id"`
	FromDate   string `form:"from_date"`
	ToDate     string `form:"to_date"`
}

// Create godoc
// @Summary Create responsibility transfer
// @Description Create a new responsibility transfer request
// @Tags Responsibility Transfer
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createTransferRequest true "Create transfer request"
// @Success 201 {object} response.APIResponse
// @Router /auth/responsibility-transfers [post]
func (h *Handler) Create(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "INVALID_REQUEST", err.Error())
		return
	}

	fromDate, err := utils.ParseDate(req.FromDate)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid from_date format", "INVALID_DATE", err.Error())
		return
	}

	toDate, err := utils.ParseDate(req.ToDate)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid to_date format", "INVALID_DATE", err.Error())
		return
	}

	result, err := h.service.CreateTransfer(ctx.Request.Context(), CreateTransferRequest{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		FromDate:   fromDate,
		ToDate:     toDate,
		Remarks:    req.Remarks,
	})
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create transfer", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "responsibility transfer created", result)
}

// List godoc
// @Summary List responsibility transfers
// @Description Get all responsibility transfers with optional filters
// @Tags Responsibility Transfer
// @Produce json
// @Security BearerAuth
// @Param from_user_id query int false "From User ID"
// @Param to_user_id query int false "To User ID"
// @Param from_date query string false "From Date (YYYY-MM-DD)"
// @Param to_date query string false "To Date (YYYY-MM-DD)"
// @Success 200 {object} response.APIResponse
// @Router /auth/responsibility-transfers [get]
func (h *Handler) List(ctx *gin.Context) {
	var req listFilterRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid query parameters", "INVALID_QUERY", err.Error())
		return
	}

	filter := TransferFilter{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
	}

	if req.FromDate != "" {
		fromDate, err := utils.ParseDate(req.FromDate)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid from_date format", "INVALID_DATE", err.Error())
			return
		}
		filter.FromDate = &fromDate
	}

	if req.ToDate != "" {
		toDate, err := utils.ParseDate(req.ToDate)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, "invalid to_date format", "INVALID_DATE", err.Error())
			return
		}
		filter.ToDate = &toDate
	}

	results, err := h.service.ListTransfers(ctx.Request.Context(), filter)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to list transfers", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "responsibility transfers retrieved", results)
}

// Approve godoc
// @Summary Approve/Reject responsibility transfer
// @Description Approve or reject a responsibility transfer request
// @Tags Responsibility Transfer
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transfer ID"
// @Param request body ApproveTransferRequest true "Approval request"
// @Success 200 {object} response.APIResponse
// @Router /auth/responsibility-transfers/{id}/approve [post]
func (h *Handler) Approve(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid ID", "INVALID_ID", err.Error())
		return
	}

	var req ApproveTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request", "INVALID_REQUEST", err.Error())
		return
	}

	// Get current user ID from context (assuming it's set by auth middleware)
	userID, _ := ctx.Get("user_id")
	approvedBy, ok := userID.(uint)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED", "user context missing")
		return
	}

	result, err := h.service.ApproveTransfer(ctx.Request.Context(), uint(id), approvedBy, req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to approve transfer", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "responsibility transfer status updated", result)
}

// Delete godoc
// @Summary Delete responsibility transfer
// @Description Delete a responsibility transfer
// @Tags Responsibility Transfer
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.APIResponse
// @Router /auth/responsibility-transfers/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid ID", "INVALID_ID", err.Error())
		return
	}

	if err := h.service.DeleteTransfer(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete transfer", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "responsibility transfer deleted", nil)
}

// Get godoc
// @Summary Get responsibility transfer
// @Description Get a responsibility transfer by ID
// @Tags Responsibility Transfer
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.APIResponse
// @Router /auth/responsibility-transfers/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid ID", "INVALID_ID", err.Error())
		return
	}

	result, err := h.service.GetTransfer(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to get transfer", "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(ctx, "responsibility transfer retrieved", result)
}
