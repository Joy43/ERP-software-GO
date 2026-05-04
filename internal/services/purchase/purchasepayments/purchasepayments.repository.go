package purchasepayments

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Repository is the single repository for all purchasepayments entities.
type Repository interface {
	CreatePaymentByGRN(ctx context.Context, m *PaymentByGRN) error
	GetAllPaymentByGRN(ctx context.Context, q PaginationQuery) ([]PaymentByGRN, int64, error)

	CreateAdvancePayment(ctx context.Context, m *AdvancePayments) error
	GetAllAdvancePayments(ctx context.Context, q PaginationQuery) ([]AdvancePayments, int64, error)

	CreateSupplierBill(ctx context.Context, m *CreateSupplierBill) error
	GetAllSupplierBills(ctx context.Context, q PaginationQuery) ([]CreateSupplierBill, int64, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository wires the DB and returns a Repository — called from app.go.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ─── PaymentByGRN ────────────────────────────────────────────────────────────

func (r *repository) CreatePaymentByGRN(ctx context.Context, m *PaymentByGRN) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("repository.CreatePaymentByGRN: %w", err)
	}
	return nil
}

func (r *repository) GetAllPaymentByGRN(ctx context.Context, q PaginationQuery) ([]PaymentByGRN, int64, error) {
	var (
		records []PaymentByGRN
		total   int64
	)
	base := r.db.WithContext(ctx).Model(&PaymentByGRN{})

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("repository.GetAllPaymentByGRN count: %w", err)
	}
	if err := base.
		Preload("GRN").
		Preload("Office").
		Preload("Supplier").
		Preload("OfficeHead").
		Preload("PaymentMode").
		Offset(q.Offset()).
		Limit(q.PageSize).
		Order("payment_by_grn_id DESC").
		Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("repository.GetAllPaymentByGRN find: %w", err)
	}
	return records, total, nil
}

// ─── AdvancePayments ─────────────────────────────────────────────────────────

func (r *repository) CreateAdvancePayment(ctx context.Context, m *AdvancePayments) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("repository.CreateAdvancePayment: %w", err)
	}
	return nil
}

func (r *repository) GetAllAdvancePayments(ctx context.Context, q PaginationQuery) ([]AdvancePayments, int64, error) {
	var (
		records []AdvancePayments
		total   int64
	)
	base := r.db.WithContext(ctx).Model(&AdvancePayments{})

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("repository.GetAllAdvancePayments count: %w", err)
	}
	if err := base.
		Preload("Office").
		Preload("AccountHead").
		Preload("SupplierHead").
		Preload("PurchaseOrderItem").
		Preload("PaymentMode").
		Offset(q.Offset()).
		Limit(q.PageSize).
		Order("advance_payment_id DESC").
		Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("repository.GetAllAdvancePayments find: %w", err)
	}
	return records, total, nil
}

// ─── SupplierBill ────────────────────────────────────────────────────────────

func (r *repository) CreateSupplierBill(ctx context.Context, m *CreateSupplierBill) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("repository.CreateSupplierBill: %w", err)
	}
	return nil
}

func (r *repository) GetAllSupplierBills(ctx context.Context, q PaginationQuery) ([]CreateSupplierBill, int64, error) {
	var (
		records []CreateSupplierBill
		total   int64
	)
	base := r.db.WithContext(ctx).Model(&CreateSupplierBill{})

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("repository.GetAllSupplierBills count: %w", err)
	}
	if err := base.
		Preload("Office").
		Preload("Supplier").
		Preload("File").
		Offset(q.Offset()).
		Limit(q.PageSize).
		Order("supplier_bill_id DESC").
		Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("repository.GetAllSupplierBills find: %w", err)
	}
	return records, total, nil
}