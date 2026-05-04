package purchasepayments

import (
	"context"
	"fmt"
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/utils"
)

// Service is the single service for all purchasepayments use-cases.
type Service struct {
	repo Repository
}

// NewService returns a *Service — called from app.go.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ─── PaymentByGRN ────────────────────────────────────────────────────────────

func (s *Service) CreatePaymentByGRN(ctx context.Context, req CreatePaymentByGRNRequest) (*PaymentByGRNResponse, error) {
	// Parse payment date
	paymentDate, err := utils.ParseDate(req.PaymentDate)
	if err != nil {
		return nil, fmt.Errorf("service.CreatePaymentByGRN: invalid payment_date: %w", err)
	}
	
	m := &PaymentByGRN{
		PaymentDate:      paymentDate,
		MoneyReceiptNo:   req.MoneyReceiptNo,
		GRNID:            req.GRNID,
		PayableAmount:    req.PayableAmount,
		PayingAmount:     req.PayingAmount,
		AdjustmentAmount: req.AdjustmentAmount,
		OfficeID:         req.OfficeID,
		SupplierID:       req.SupplierID,
		OfficeHeadID:     req.OfficeHeadID,
		PaymentModeID:    req.PaymentModeID,
	}
	if err := s.repo.CreatePaymentByGRN(ctx, m); err != nil {
		return nil, fmt.Errorf("service.CreatePaymentByGRN: %w", err)
	}
	return toPaymentByGRNResponse(m), nil
}

func (s *Service) GetAllPaymentByGRN(ctx context.Context, q PaginationQuery) (PaginatedResponse, error) {
	q.DefaultPagination()
	records, total, err := s.repo.GetAllPaymentByGRN(ctx, q)
	if err != nil {
		return PaginatedResponse{}, fmt.Errorf("service.GetAllPaymentByGRN: %w", err)
	}
	resp := make([]PaymentByGRNResponse, len(records))
	for i := range records {
		resp[i] = *toPaymentByGRNResponse(&records[i])
	}
	return NewPaginatedResponse(resp, total, q), nil
}

func toPaymentByGRNResponse(m *PaymentByGRN) *PaymentByGRNResponse {
	return &PaymentByGRNResponse{
		PaymentByGRNID:   m.PaymentByGRNID,
		PaymentDate:      m.PaymentDate,
		MoneyReceiptNo:   m.MoneyReceiptNo,
		GRNID:            m.GRNID,
		PayableAmount:    m.PayableAmount,
		PayingAmount:     m.PayingAmount,
		AdjustmentAmount: m.AdjustmentAmount,
		OfficeID:         m.OfficeID,
		SupplierID:       m.SupplierID,
		OfficeHeadID:     m.OfficeHeadID,
		PaymentModeID:    m.PaymentModeID,
		GRN:              m.GRN,
		Office:           m.Office,
		Supplier:         m.Supplier,
		OfficeHead:       m.OfficeHead,
		PaymentMode:      m.PaymentMode,
	}
}

// ─── AdvancePayments ─────────────────────────────────────────────────────────

func (s *Service) CreateAdvancePayment(ctx context.Context, req CreateAdvancePaymentRequest) (*AdvancePaymentResponse, error) {
	// Parse payment date
	paymentDate, err := utils.ParseDate(req.PaymentDate)
	if err != nil {
		return nil, fmt.Errorf("service.CreateAdvancePayment: invalid payment_date: %w", err)
	}
	
	m := &AdvancePayments{
		PaymentDate:    paymentDate,
		CashAmount:     req.CashAmount,
		Narration:      req.Narration,
		Amount:         req.Amount,
		LcNo:           req.LcNo,
		OfficeID:       req.OfficeID,
		AccountHeadID:  req.AccountHeadID,
		SupplierHeadID: req.SupplierHeadID,
		POID:           req.POID,
		PaymentModeID:  req.PaymentModeID,
	}
	if err := s.repo.CreateAdvancePayment(ctx, m); err != nil {
		return nil, fmt.Errorf("service.CreateAdvancePayment: %w", err)
	}
	return toAdvancePaymentResponse(m), nil
}

func (s *Service) GetAllAdvancePayments(ctx context.Context, q PaginationQuery) (PaginatedResponse, error) {
	q.DefaultPagination()
	records, total, err := s.repo.GetAllAdvancePayments(ctx, q)
	if err != nil {
		return PaginatedResponse{}, fmt.Errorf("service.GetAllAdvancePayments: %w", err)
	}
	resp := make([]AdvancePaymentResponse, len(records))
	for i := range records {
		resp[i] = *toAdvancePaymentResponse(&records[i])
	}
	return NewPaginatedResponse(resp, total, q), nil
}

func toAdvancePaymentResponse(m *AdvancePayments) *AdvancePaymentResponse {
	return &AdvancePaymentResponse{
		AdvancePaymentID: m.AdvancePaymentID,
		PaymentDate:      m.PaymentDate,
		CashAmount:       m.CashAmount,
		Narration:        m.Narration,
		Amount:           m.Amount,
		LcNo:             m.LcNo,
		OfficeID:         m.OfficeID,
		AccountHeadID:    m.AccountHeadID,
		SupplierHeadID:   m.SupplierHeadID,
		POID:             m.POID,
		PaymentModeID:    m.PaymentModeID,
		Office:           m.Office,
		AccountHead:      m.AccountHead,
		SupplierHead:     m.SupplierHead,
		PurchaseOrder:    m.PurchaseOrderItem,
		PaymentMode:      m.PaymentMode,
	}
}

// ─── SupplierBill ────────────────────────────────────────────────────────────

func (s *Service) CreateSupplierBill(ctx context.Context, req CreateSupplierBillRequest) (*SupplierBillResponse, error) {
	// Parse create date
	createDate, err := utils.ParseDate(req.CreateDate)
	if err != nil {
		return nil, fmt.Errorf("service.CreateSupplierBill: invalid create_date: %w", err)
	}
	
	// Parse tent pay date (optional)
	var tentPayDate time.Time
	if req.TentPayDate != "" {
		var parseErr error
		tentPayDate, parseErr = utils.ParseDate(req.TentPayDate)
		if parseErr != nil {
			return nil, fmt.Errorf("service.CreateSupplierBill: invalid tent_pay_date: %w", parseErr)
		}
	}
	
	m := &CreateSupplierBill{
		CreateDate:   createDate,
		TentPayDate:  tentPayDate,
		BillNo:       req.BillNo,
		BillAmount:   req.BillAmount,
		Discount:     req.Discount,
		Advance:      req.Advance,
		NetPay:       req.NetPay,
		VatChallanNo: req.VatChallanNo,
		Vat:          req.Vat,
		Remarks:      req.Remarks,
		Sd:           req.Sd,
		OfficeID:     req.OfficeID,
		SupplierID:   req.SupplierID,
		FileID:       req.FileID,
	}
	if err := s.repo.CreateSupplierBill(ctx, m); err != nil {
		return nil, fmt.Errorf("service.CreateSupplierBill: %w", err)
	}
	return toSupplierBillResponse(m), nil
}

func (s *Service) GetAllSupplierBills(ctx context.Context, q PaginationQuery) (PaginatedResponse, error) {
	q.DefaultPagination()
	records, total, err := s.repo.GetAllSupplierBills(ctx, q)
	if err != nil {
		return PaginatedResponse{}, fmt.Errorf("service.GetAllSupplierBills: %w", err)
	}
	resp := make([]SupplierBillResponse, len(records))
	for i := range records {
		resp[i] = *toSupplierBillResponse(&records[i])
	}
	return NewPaginatedResponse(resp, total, q), nil
}

func toSupplierBillResponse(m *CreateSupplierBill) *SupplierBillResponse {
	return &SupplierBillResponse{
		SupplierBillID: m.SupplierBillID,
		CreateDate:     m.CreateDate,
		TentPayDate:    m.TentPayDate,
		BillNo:         m.BillNo,
		BillAmount:     m.BillAmount,
		Discount:       m.Discount,
		Advance:        m.Advance,
		NetPay:         m.NetPay,
		VatChallanNo:   m.VatChallanNo,
		Vat:            m.Vat,
		Remarks:        m.Remarks,
		Sd:             m.Sd,
		OfficeID:       m.OfficeID,
		SupplierID:     m.SupplierID,
		FileID:         m.FileID,
		Office:         m.Office,
		Supplier:       m.Supplier,
		File:           m.File,
	}
}