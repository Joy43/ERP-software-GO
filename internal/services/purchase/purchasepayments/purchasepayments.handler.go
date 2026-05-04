package purchasepayments

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	svc *Service
}

// NewHandler returns a *Handler — called from app.go.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// ─── Shared helpers ───────────────────────────────────────────────────────────

func successResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func createdResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": data})
}

func errorResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"success": false, "message": msg})
}

func validationErrorResponse(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	fields := make(map[string]string, len(errs))
	for _, e := range errs {
		fields[e.Field()] = e.Tag()
	}
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"success": false,
		"message": "validation failed",
		"errors":  fields,
	})
}

// ───--------- PaymentByGRN ───────────────────────────────

// CreatePaymentByGRN godoc
// @Summary  Create a payment by GRN
// @Tags    purchase------------------------- purchase-payments
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    request body CreatePaymentByGRNRequest true "Payload"
// @Success  201  {object} PaymentByGRNResponse
// @Router   /purchase/payment-by-grn [post]
func (h *Handler) CreatePaymentByGRN(c *gin.Context) {
	var req CreatePaymentByGRNRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	if err := validate.Struct(req); err != nil {
		validationErrorResponse(c, err)
		return
	}
	resp, err := h.svc.CreatePaymentByGRN(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	createdResponse(c, resp)
}

// GetAllPaymentByGRN godoc
// @Summary  List payments by GRN
// @Tags    purchase------------------------- purchase-payments
// @Produce  json
// @Security BearerAuth
// @Param    page      query int false "Page"           default(1)
// @Param    page_size query int false "Records/page"   default(10)
// @Success  200  {object} PaginatedResponse
// @Router   /purchase/payment-by-grn [get]
func (h *Handler) GetAllPaymentByGRN(c *gin.Context) {
	var q PaginationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid query params: "+err.Error())
		return
	}
	result, err := h.svc.GetAllPaymentByGRN(c.Request.Context(), q)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	successResponse(c, result)
}

// ─── AdvancePayments ──────────────────────────────────────────────────────────

// CreateAdvancePayment godoc
// @Summary  Create an advance payment
// @Tags    purchase------------------------- purchase-payments
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    request body CreateAdvancePaymentRequest true "Payload"
// @Success  201  {object} AdvancePaymentResponse
// @Router   /purchase/advance-payments [post]
func (h *Handler) CreateAdvancePayment(c *gin.Context) {
	var req CreateAdvancePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	if err := validate.Struct(req); err != nil {
		validationErrorResponse(c, err)
		return
	}
	resp, err := h.svc.CreateAdvancePayment(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	createdResponse(c, resp)
}

// GetAllAdvancePayments godoc
// @Summary  List advance payments
// @Tags    purchase------------------------- purchase-payments
// @Produce  json
// @Security BearerAuth
// @Param    page      query int false "Page"           default(1)
// @Param    page_size query int false "Records/page"   default(10)
// @Success  200  {object} PaginatedResponse
// @Router   /purchase/advance-payments [get]
func (h *Handler) GetAllAdvancePayments(c *gin.Context) {
	var q PaginationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid query params: "+err.Error())
		return
	}
	result, err := h.svc.GetAllAdvancePayments(c.Request.Context(), q)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	successResponse(c, result)
}

// ─── SupplierBill ─────────────────────────────────────────────────────────────

// CreateSupplierBill godoc
// @Summary  Create a supplier bill
// @Tags    purchase------------------------- purchase-payments
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    request body CreateSupplierBillRequest true "Payload"
// @Success  201  {object} SupplierBillResponse
// @Router   /purchase/supplier-bills [post]
func (h *Handler) CreateSupplierBill(c *gin.Context) {
	var req CreateSupplierBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	if err := validate.Struct(req); err != nil {
		validationErrorResponse(c, err)
		return
	}
	resp, err := h.svc.CreateSupplierBill(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	createdResponse(c, resp)
}

// GetAllSupplierBills godoc
// @Summary  List supplier bills
// @Tags    purchase------------------------- purchase-payments
// @Produce  json
// @Security BearerAuth
// @Param    page      query int false "Page"           default(1)
// @Param    page_size query int false "Records/page"   default(10)
// @Success  200  {object} PaginatedResponse
// @Router   /purchase/supplier-bills [get]
func (h *Handler) GetAllSupplierBills(c *gin.Context) {
	var q PaginationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid query params: "+err.Error())
		return
	}
	result, err := h.svc.GetAllSupplierBills(c.Request.Context(), q)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	successResponse(c, result)
}