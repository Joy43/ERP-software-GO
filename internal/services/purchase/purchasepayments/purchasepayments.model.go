package purchasepayments

import (
	"time"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/office"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/payment_mode"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/user"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/supplier"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/grn"
	purchaseorder "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/order"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/uploads"
)

// PaymentByGRN records a supplier payment linked to a specific GRN (Goods Receipt Note).
type PaymentByGRN struct {
	PaymentByGRNID   int       `gorm:"primaryKey;column:payment_by_grn_id"                                              json:"payment_by_grn_id"`
	PaymentDate      time.Time `gorm:"not null;column:payment_date"                                                     json:"payment_date"`
	MoneyReceiptNo   string    `gorm:"size:100;not null;uniqueIndex;column:money_receipt_no"                            json:"money_receipt_no"`
	GRNID            int       `gorm:"not null;index;column:grn_id"                                                     json:"grn_id"`
	PayableAmount    float64   `gorm:"type:decimal(18,4);not null;default:0;column:payable_amount"                      json:"payable_amount"`
	PayingAmount     float64   `gorm:"type:decimal(18,4);not null;default:0;column:paying_amount"                       json:"paying_amount"`
	AdjustmentAmount float64   `gorm:"type:decimal(18,4);not null;default:0;column:adjustment_amount"                   json:"adjustment_amount"`
	OfficeID         int       `gorm:"not null;index;column:office_id"                                                  json:"office_id"`
	SupplierID       int       `gorm:"not null;index;column:supplier_id"                                                json:"supplier_id"`
	OfficeHeadID     int       `gorm:"not null;index;column:office_head_id"                                             json:"office_head_id"`
	PaymentModeID    uint      `gorm:"not null;index;column:payment_mode_id"                                            json:"payment_mode_id"`

	//------------ Associations------------
	GRN         *grn.GRNItem              `gorm:"foreignKey:GRNID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"         json:"grn"`
	Office      *office.Office            `gorm:"foreignKey:OfficeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"      json:"office"`
	Supplier    *supplier.Supplier        `gorm:"foreignKey:SupplierID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"    json:"supplier"`
	OfficeHead  *user.User                `gorm:"foreignKey:OfficeHeadID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"  json:"office_head"`
	PaymentMode payment_mode.PaymentMode  `gorm:"foreignKey:PaymentModeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"payment_mode"`
}

// TableName overrides the default GORM table name.
func (PaymentByGRN) TableName() string { return "payment_by_grns" }

// ─────────────────────────────────────────────────────────────────────────────

// AdvancePayments records advance payments made to suppliers against a Purchase Order.
type AdvancePayments struct {
	AdvancePaymentID int       `gorm:"primaryKey;column:advance_payment_id"                                             json:"advance_payment_id"`
	PaymentDate      time.Time `gorm:"not null;column:payment_date"                                                     json:"payment_date"`
	CashAmount       float64   `gorm:"type:decimal(18,4);not null;default:0;column:cash_amount"                         json:"cash_amount"`
	Narration        string    `gorm:"type:text;column:narration"                                                       json:"narration"`
	Amount           float64   `gorm:"type:decimal(18,4);not null;default:0;column:amount"                              json:"amount"`
	LcNo             string    `gorm:"size:100;index;column:lc_no"                                                      json:"lc_no"`
	OfficeID         int       `gorm:"not null;index;column:office_id"                                                  json:"office_id"`
	AccountHeadID    int       `gorm:"not null;index;column:account_head_id"                                            json:"account_head_id"`
	SupplierHeadID   int       `gorm:"not null;index;column:supplier_head_id"                                           json:"supplier_head_id"`
	POID             int       `gorm:"not null;index;column:po_id"                                                      json:"po_id"`
	PaymentModeID    uint      `gorm:"not null;index;column:payment_mode_id"                                            json:"payment_mode_id"`

	//----------  Associations ----------
	Office              *office.Office                      `gorm:"foreignKey:OfficeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"         json:"office"`
	AccountHead         *user.User                          `gorm:"foreignKey:AccountHeadID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"    json:"account_head"`
	SupplierHead        *user.User                          `gorm:"foreignKey:SupplierHeadID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"   json:"supplier_head"`
	PurchaseOrderItem   *purchaseorder.PurchaseOrder        `gorm:"foreignKey:POID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"             json:"purchase_order"`
	PaymentMode         payment_mode.PaymentMode            `gorm:"foreignKey:PaymentModeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"    json:"payment_mode"`
}

// TableName overrides the default GORM table name.
func (AdvancePayments) TableName() string { return "advance_payments" }

// ─────────────────────────────────────────────────────────────────────────────

// CreateSupplierBill records a formal bill issued by a supplier.
type CreateSupplierBill struct {
	SupplierBillID int       `gorm:"primaryKey;column:supplier_bill_id"                                               json:"supplier_bill_id"`
	CreateDate     time.Time `gorm:"not null;column:create_date"                                                      json:"create_date"`
	TentPayDate    time.Time `gorm:"column:tent_pay_date"                                                             json:"tent_pay_date"`
	BillNo         string    `gorm:"size:100;not null;uniqueIndex;column:bill_no"                                     json:"bill_no"         validate:"required"`
	BillAmount     float64   `gorm:"type:decimal(18,4);not null;default:0;column:bill_amount"                         json:"bill_amount"`
	Discount       float64   `gorm:"type:decimal(18,4);not null;default:0;column:discount"                            json:"discount"`
	Advance        float64   `gorm:"type:decimal(18,4);not null;default:0;column:advance"                             json:"advance"`
	NetPay         float64   `gorm:"type:decimal(18,4);not null;default:0;column:net_pay"                             json:"net_pay"`
	VatChallanNo   string    `gorm:"size:100;column:vat_challan_no"                                                   json:"vat_challan_no"`
	Vat            float64   `gorm:"type:decimal(18,4);not null;default:0;column:vat"                                 json:"vat"`
	Remarks        string    `gorm:"type:text;column:remarks"                                                         json:"remarks"`
	Sd             float64   `gorm:"type:decimal(18,4);not null;default:0;column:sd"                                  json:"sd"`
	OfficeID       int       `gorm:"not null;index;column:office_id"                                                  json:"office_id"`
	SupplierID     int       `gorm:"not null;index;column:supplier_id"                                                json:"supplier_id"`
	FileID         *uint     `gorm:"index;column:file_id"                                                             json:"file_id"         binding:"omitempty,gt=0"`

	// Associations
	Office   *office.Office      `gorm:"foreignKey:OfficeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"   json:"office"`
	Supplier *supplier.Supplier  `gorm:"foreignKey:SupplierID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"supplier"`
	File     *uploads.File       `gorm:"foreignKey:FileID"                                                   json:"file,omitempty"`
}

// TableName overrides the default GORM table name.
func (CreateSupplierBill) TableName() string { return "supplier_bills" }