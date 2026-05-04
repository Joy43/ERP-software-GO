package purchaseorder

import (
	"time"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/office"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/user"
	umo_measurement "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/uom"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/items"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/supplier"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/location"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/requisitions"
)

// =============================================
// ENUMS
// =============================================

type POStatus string

const (
	POStatusPending           POStatus = "PENDING"
	POStatusIssued            POStatus = "ISSUED"
	POStatusConfirmed         POStatus = "CONFIRMED"
	POStatusPartiallyReceived POStatus = "PARTIALLY_RECEIVED"
	POStatusFullyReceived     POStatus = "FULLY_RECEIVED"
	POStatusCancelled         POStatus = "CANCELLED"
)

type POOrderType string

const (
	POTypeDirect           POOrderType = "DIRECT"
	POTypeRequisitionBased POOrderType = "REQUISITION_BASED"
	POTypeContract         POOrderType = "CONTRACT"
)

var ValidPOStatusTransitions = map[POStatus][]POStatus{
	POStatusPending:           {POStatusIssued, POStatusCancelled},
	POStatusIssued:            {POStatusConfirmed, POStatusCancelled},
	POStatusConfirmed:         {POStatusPartiallyReceived, POStatusFullyReceived, POStatusCancelled},
	POStatusPartiallyReceived: {POStatusFullyReceived, POStatusCancelled},
	POStatusFullyReceived:     {},
	POStatusCancelled:         {},
}

func (current POStatus) IsValidTransition(next POStatus) bool {
	allowed, ok := ValidPOStatusTransitions[current]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}

// =============================================
// PURCHASE ORDER MODEL
// =============================================

type PurchaseOrder struct {
	ID uint `gorm:"primaryKey" json:"id"`

	PONumber     string     `gorm:"size:50;uniqueIndex;not null" json:"po_number"`
	PODate       time.Time  `gorm:"type:date;not null" json:"po_date"`
	DeliveryDate *time.Time `gorm:"type:date" json:"delivery_date,omitempty"`

	RequisitionID *uint                      `gorm:"index" json:"requisition_id,omitempty"`
	Requisition   *requisitions.Requisition  `gorm:"foreignKey:RequisitionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"requisition,omitempty"`

	OrderType POOrderType `gorm:"type:enum('DIRECT','REQUISITION_BASED','CONTRACT');default:'DIRECT'" json:"order_type"`

	OfficeID   uint           `gorm:"not null" json:"office_id"`
	Office     *office.Office `gorm:"foreignKey:OfficeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"office,omitempty"`
	LocationID uint           `gorm:"not null" json:"location_id"`
	Location   *location.Location `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"location,omitempty"`

	SupplierID   uint               `gorm:"not null;index" json:"supplier_id"`
	Supplier     *supplier.Supplier `gorm:"foreignKey:SupplierID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"supplier,omitempty"`
	PaymentTerms *string            `gorm:"size:100" json:"payment_terms,omitempty"`

	Subtotal       float64 `gorm:"type:decimal(15,2);default:0" json:"subtotal"`
	VatAmount      float64 `gorm:"type:decimal(15,2);default:0" json:"vat_amount"`
	DiscountAmount float64 `gorm:"type:decimal(15,2);default:0" json:"discount_amount"`
	TotalAmount    float64 `gorm:"type:decimal(15,2);default:0" json:"total_amount"`

	GeneralRemarks  *string `gorm:"type:text" json:"general_remarks,omitempty"`
	ShippingAddress *string `gorm:"type:text" json:"shipping_address,omitempty"`

	Status POStatus `gorm:"type:enum('PENDING','ISSUED','CONFIRMED','PARTIALLY_RECEIVED','FULLY_RECEIVED','CANCELLED');default:'PENDING';index" json:"status"`

	CreatedByID *uint      `gorm:"index" json:"created_by_id,omitempty"`
	CreatedBy   *user.User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"created_by,omitempty"`
	ApprovedByID *uint     `json:"approved_by_id,omitempty"`
	ApprovedBy  *user.User `gorm:"foreignKey:ApprovedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`

	Items []PurchaseOrderItem `gorm:"foreignKey:POID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"items,omitempty"`

	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (PurchaseOrder) TableName() string { return "purchase_orders" }

// =============================================
// PURCHASE ORDER ITEM MODEL
// =============================================

type PurchaseOrderItem struct {
	ID   uint `gorm:"primaryKey" json:"id"`
	POID uint `gorm:"not null;index" json:"po_id"`

	RequisitionItemID *uint `gorm:"index" json:"requisition_item_id,omitempty"`

	ItemID uint        `gorm:"not null;index" json:"item_id"`
	Item   items.Items `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"item,omitempty"`

	OrderQuantity    float64 `gorm:"type:decimal(15,2);not null" json:"order_quantity"`
	ReceivedQuantity float64 `gorm:"type:decimal(15,2);default:0" json:"received_quantity"`

	UOMID     *uint                    `json:"uom_id,omitempty"`
	UOM       *umo_measurement.Uom     `gorm:"foreignKey:UOMID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"uom,omitempty"`
	UnitPrice float64   `gorm:"type:decimal(15,2);not null" json:"unit_price"`

	VatPercentage      float64 `gorm:"type:decimal(5,2);default:0" json:"vat_percentage"`
	VatAmount          float64 `gorm:"type:decimal(15,2);default:0" json:"vat_amount"`
	DiscountPercentage float64 `gorm:"type:decimal(5,2);default:0" json:"discount_percentage"`
	DiscountAmount     float64 `gorm:"type:decimal(15,2);default:0" json:"discount_amount"`
	TotalAmount        float64 `gorm:"type:decimal(15,2);default:0" json:"total_amount"`
	Remarks *string `gorm:"type:text" json:"remarks,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (PurchaseOrderItem) TableName() string { return "purchase_order_items" }
