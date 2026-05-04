package grn

import (
	"time"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/office"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/payment_mode"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/user"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/items"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/minor_category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/sub_category"
	umo_measurement "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/uom"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/supplier"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/uploads"
	purchaseorder "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/order"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/location"
)

// =============================================
//----------- ENUMS-----------------
// =============================================

type GRNReceiveType string

const (
	GRNDirect    GRNReceiveType = "DIRECT"
	GRNAgainstPO GRNReceiveType = "AGAINST_PO"
)

type GRNStatus string

const (
	GRNStatusPending   GRNStatus = "PENDING"
	GRNStatusConfirmed GRNStatus = "CONFIRMED"
	GRNStatusCancelled GRNStatus = "CANCELLED"
)

// =============================================
// GOODS RECEIPT NOTE MODEL
// =============================================

type GoodsReceiptNote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GRNNumber string    `gorm:"size:50;uniqueIndex;not null" json:"grn_number"`
	GRNDate   time.Time `gorm:"type:date;not null" json:"grn_date"`

	ReceiveType GRNReceiveType `gorm:"type:enum('DIRECT','AGAINST_PO');not null" json:"receive_type"`
	Status      GRNStatus      `gorm:"type:enum('PENDING','CONFIRMED','CANCELLED');default:'PENDING';index" json:"status"`

	POID          *uint                        `gorm:"index" json:"po_id,omitempty"`
	PO            *purchaseorder.PurchaseOrder `gorm:"foreignKey:POID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"po,omitempty"`
	RequisitionID *uint                        `gorm:"index" json:"requisition_id,omitempty"`

	OfficeID   uint               `gorm:"not null" json:"office_id"`
	Office     *office.Office     `gorm:"foreignKey:OfficeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"office,omitempty"`
	LocationID uint               `gorm:"not null" json:"location_id"`
	Location   *location.Location `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"location,omitempty"`

	SupplierID uint               `gorm:"not null;index" json:"supplier_id"`
	Supplier   *supplier.Supplier `gorm:"foreignKey:SupplierID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"supplier,omitempty"`

	ChallanNo              *string    `gorm:"size:100" json:"challan_no,omitempty"`
	ChallanDate            *time.Time `gorm:"type:date" json:"challan_date,omitempty"`
	SalesInvoiceNumber     *string    `gorm:"size:100" json:"sales_invoice_number,omitempty"`
	VATChallanNumber       *string    `gorm:"size:100" json:"vat_challan_number,omitempty"`
	DeliveryNumber         *string    `gorm:"size:100" json:"delivery_number,omitempty"`
	ShippingAddress        *string    `gorm:"type:text" json:"shipping_address,omitempty"`
	ShipmentDocumentNumber *string    `gorm:"size:100" json:"shipment_document_number,omitempty"`

	PaymentMethodID *uint                    `json:"payment_method_id,omitempty"`
	PaymentMethod   *payment_mode.PaymentMode `gorm:"foreignKey:PaymentMethodID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"payment_method,omitempty"`

	Remarks  *string `gorm:"type:text" json:"remarks,omitempty"`
	FileID   *uint          `gorm:"index;column:file_id"                                    json:"file_id,omitempty"  binding:"omitempty,gt=0"`
	File     *uploads.File  `gorm:"foreignKey:FileID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"file,omitempty"`

	CreatedByID  *uint      `gorm:"index" json:"created_by_id,omitempty"`
	CreatedBy    *user.User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"created_by,omitempty"`
	ReceivedByID *uint      `json:"received_by_id,omitempty"`
	ReceivedBy   *user.User `gorm:"foreignKey:ReceivedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"received_by,omitempty"`

	Items []GRNItem `gorm:"foreignKey:GRNID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"items,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (GoodsReceiptNote) TableName() string { return "goods_receipt_notes" }

// =============================================
//---------  GRN ITEM MODEL ----------
// =============================================

type GRNItem struct {
	ID    uint `gorm:"primaryKey" json:"id"`
	GRNID uint `gorm:"not null;index" json:"grn_id"`

	POItemID *uint                            `gorm:"index" json:"po_item_id,omitempty"`
	POItem   *purchaseorder.PurchaseOrderItem `gorm:"foreignKey:POItemID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"po_item,omitempty"`

	ItemID uint        `gorm:"not null;index" json:"item_id"`
	Item   items.Items `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"item,omitempty"`

	ReceivedQuantity float64  `gorm:"type:decimal(15,2);not null" json:"received_quantity"`
	UOMID            *uint                    `json:"uom_id,omitempty"`
	UOM              *umo_measurement.Uom     `gorm:"foreignKey:UOMID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"uom,omitempty"`
	PurchasePrice    float64  `gorm:"type:decimal(15,2);not null" json:"purchase_price"`

	VatPercentage      float64 `gorm:"type:decimal(5,2);default:0" json:"vat_percentage"`
	VatAmount          float64 `gorm:"type:decimal(15,2);default:0" json:"vat_amount"`
	DiscountPercentage float64 `gorm:"type:decimal(5,2);default:0" json:"discount_percentage"`
	DiscountAmount     float64 `gorm:"type:decimal(15,2);default:0" json:"discount_amount"`
	TotalAmount        float64 `gorm:"type:decimal(15,2);default:0" json:"total_amount"`

	CategoryID      *uint                         `json:"category_id,omitempty"`
	Category        *category.Category            `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	SubCategoryID   *uint                         `json:"sub_category_id,omitempty"`
	SubCategory     *sub_category.SubCategory     `gorm:"foreignKey:SubCategoryID" json:"sub_category,omitempty"`
	MinorCategoryID *uint                         `json:"minor_category_id,omitempty"`
	MinorCategory   *minor_category.MinorCategory `gorm:"foreignKey:MinorCategoryID" json:"minor_category,omitempty"`

	Remarks *string `gorm:"type:text" json:"remarks,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (GRNItem) TableName() string { return "grn_items" }
