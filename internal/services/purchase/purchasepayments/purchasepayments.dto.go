package purchasepayments

import "time"

// ─────────────────────────────────────────────────────────────────────────────
// PaymentByGRN DTOs
// ─────────────────────────────────────────────────────────────────────────────

// CreatePaymentByGRNRequest is the request body for creating a PaymentByGRN record.
type CreatePaymentByGRNRequest struct {
	PaymentDate      string  `json:"payment_date"       validate:"required"`
	MoneyReceiptNo   string  `json:"money_receipt_no"   validate:"required,max=100"`
	GRNID            int     `json:"grn_id"             validate:"required,gt=0"`
	PayableAmount    float64 `json:"payable_amount"     validate:"required,gte=0"`
	PayingAmount     float64 `json:"paying_amount"      validate:"required,gte=0"`
	AdjustmentAmount float64 `json:"adjustment_amount"  validate:"gte=0"`
	OfficeID         int     `json:"office_id"          validate:"required,gt=0"`
	SupplierID       int     `json:"supplier_id"        validate:"required,gt=0"`
	OfficeHeadID     int     `json:"office_head_id"     validate:"required,gt=0"`
	PaymentModeID    uint    `json:"payment_mode_id"    validate:"required,gt=0"`
}

// PaymentByGRNResponse is the response body for a PaymentByGRN record.
type PaymentByGRNResponse struct {
	PaymentByGRNID   int         `json:"payment_by_grn_id"`
	PaymentDate      time.Time   `json:"payment_date"`
	MoneyReceiptNo   string      `json:"money_receipt_no"`
	GRNID            int         `json:"grn_id"`
	PayableAmount    float64     `json:"payable_amount"`
	PayingAmount     float64     `json:"paying_amount"`
	AdjustmentAmount float64     `json:"adjustment_amount"`
	OfficeID         int         `json:"office_id"`
	SupplierID       int         `json:"supplier_id"`
	OfficeHeadID     int         `json:"office_head_id"`
	PaymentModeID    uint        `json:"payment_mode_id"`
	GRN              interface{} `json:"grn,omitempty"`
	Office           interface{} `json:"office,omitempty"`
	Supplier         interface{} `json:"supplier,omitempty"`
	OfficeHead       interface{} `json:"office_head,omitempty"`
	PaymentMode      interface{} `json:"payment_mode,omitempty"`
}

// ─────────────────────────────────────────────────────────────────────────────
// AdvancePayments DTOs
// ─────────────────────────────────────────────────────────────────────────────

// CreateAdvancePaymentRequest is the request body for creating an AdvancePayment.
type CreateAdvancePaymentRequest struct {
	PaymentDate    string  `json:"payment_date"     validate:"required"`
	CashAmount     float64 `json:"cash_amount"      validate:"required,gte=0"`
	Narration      string  `json:"narration"        validate:"max=1000"`
	Amount         float64 `json:"amount"           validate:"required,gte=0"`
	LcNo           string  `json:"lc_no"            validate:"max=100"`
	OfficeID       int     `json:"office_id"        validate:"required,gt=0"`
	AccountHeadID  int     `json:"account_head_id"  validate:"required,gt=0"`
	SupplierHeadID int     `json:"supplier_head_id" validate:"required,gt=0"`
	POID           int     `json:"po_id"            validate:"required,gt=0"`
	PaymentModeID  uint    `json:"payment_mode_id"  validate:"required,gt=0"`
}

// AdvancePaymentResponse is the response body for an AdvancePayment.
type AdvancePaymentResponse struct {
	AdvancePaymentID int                              `json:"advance_payment_id"`
	PaymentDate      time.Time                        `json:"payment_date"`
	CashAmount       float64                          `json:"cash_amount"`
	Narration        string                           `json:"narration"`
	Amount           float64                          `json:"amount"`
	LcNo             string                           `json:"lc_no"`
	OfficeID         int                              `json:"office_id"`
	AccountHeadID    int                              `json:"account_head_id"`
	SupplierHeadID   int                              `json:"supplier_head_id"`
	POID             int                              `json:"po_id"`
	PaymentModeID    uint                             `json:"payment_mode_id"`
	Office           interface{}                      `json:"office,omitempty"`
	AccountHead      interface{}                      `json:"account_head,omitempty"`
	SupplierHead     interface{}                      `json:"supplier_head,omitempty"`
	PurchaseOrder    interface{}                      `json:"purchase_order,omitempty"`
	PaymentMode      interface{}                      `json:"payment_mode,omitempty"`
}

// ─────────────────────────────────────────────────────────────────────────────
// CreateSupplierBill DTOs
// ─────────────────────────────────────────────────────────────────────────────

// CreateSupplierBillRequest is the request body for creating a SupplierBill.
type CreateSupplierBillRequest struct {
	CreateDate   string  `json:"create_date"    validate:"required"`
	TentPayDate  string  `json:"tent_pay_date"`
	BillNo       string  `json:"bill_no"        validate:"required,max=100"`
	BillAmount   float64 `json:"bill_amount"    validate:"required,gte=0"`
	Discount     float64 `json:"discount"       validate:"gte=0"`
	Advance      float64 `json:"advance"        validate:"gte=0"`
	NetPay       float64 `json:"net_pay"        validate:"required,gte=0"`
	VatChallanNo string  `json:"vat_challan_no" validate:"max=100"`
	Vat          float64 `json:"vat"            validate:"gte=0"`
	Remarks      string  `json:"remarks"        validate:"max=2000"`
	Sd           float64 `json:"sd"             validate:"gte=0"`
	OfficeID     int     `json:"office_id"      validate:"required,gt=0"`
	SupplierID   int     `json:"supplier_id"    validate:"required,gt=0"`
	FileID       *uint   `json:"file_id"        binding:"omitempty,gt=0"`
}

// SupplierBillResponse is the response body for a SupplierBill.
type SupplierBillResponse struct {
	SupplierBillID int         `json:"supplier_bill_id"`
	CreateDate     time.Time   `json:"create_date"`
	TentPayDate    time.Time   `json:"tent_pay_date"`
	BillNo         string      `json:"bill_no"`
	BillAmount     float64     `json:"bill_amount"`
	Discount       float64     `json:"discount"`
	Advance        float64     `json:"advance"`
	NetPay         float64     `json:"net_pay"`
	VatChallanNo   string      `json:"vat_challan_no"`
	Vat            float64     `json:"vat"`
	Remarks        string      `json:"remarks"`
	Sd             float64     `json:"sd"`
	OfficeID       int         `json:"office_id"`
	SupplierID     int         `json:"supplier_id"`
	FileID         *uint       `json:"file_id"`
	Office         interface{} `json:"office,omitempty"`
	Supplier       interface{} `json:"supplier,omitempty"`
	File           interface{} `json:"file,omitempty"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Shared pagination / filter DTOs
// ─────────────────────────────────────────────────────────────────────────────

// PaginationQuery holds common query-string params for list endpoints.
type PaginationQuery struct {
	Page     int `form:"page"      validate:"min=1"`
	PageSize int `form:"page_size" validate:"min=1,max=100"`
}

// DefaultPagination sets safe defaults when params are absent.
func (p *PaginationQuery) DefaultPagination() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
}

// Offset calculates the DB offset from page / page_size.
func (p *PaginationQuery) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// PaginatedResponse wraps any list response with meta.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	TotalCount int64       `json:"total_count"`
	Page       int         `json:"page"`
	PageLimit  int         `json:"page_limit"`
	TotalPages int         `json:"total_pages"`
}

// NewPaginatedResponse constructs a PaginatedResponse.
func NewPaginatedResponse(data interface{}, total int64, q PaginationQuery) PaginatedResponse {
	totalPages := int(total) / q.PageSize
	if int(total)%q.PageSize != 0 {
		totalPages++
	}
	return PaginatedResponse{
		Data:       data,
		TotalCount: total,
		Page:       q.Page,
		PageLimit:  q.PageSize,
		TotalPages: totalPages,
	}
}