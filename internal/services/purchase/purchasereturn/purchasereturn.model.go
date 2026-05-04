package purchasereturn

import (
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/office"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/items"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/supplier"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/location"
)

type PurchaseReturn struct {
	ID               uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	ReturnNumber     string  `gorm:"uniqueIndex;not null;size:50" json:"return_number"`
	Quantity         float64 `gorm:"type:decimal(15,2);not null" json:"quantity"`
	SellingPrice     float64 `gorm:"type:decimal(15,2);not null" json:"selling_price"`
	Remarks          string  `gorm:"type:text" json:"remarks"`

	// -------- Foreign keys --------
	OfficeID   uint             `gorm:"not null;index" json:"office_id"`
	Office     *office.Office   `gorm:"foreignKey:OfficeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"office,omitempty"`

	LocationID uint               `gorm:"not null;index" json:"location_id"`
	Location   *location.Location `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"location,omitempty"`

	SupplierID uint               `gorm:"not null;index" json:"supplier_id"`
	Supplier   *supplier.Supplier `gorm:"foreignKey:SupplierID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"supplier,omitempty"`

	ItemID uint         `gorm:"not null;index" json:"item_id"`
	Item   *items.Items `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"item,omitempty"`

	// -------- Timestamps --------
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (PurchaseReturn) TableName() string {
	return "purchase_returns"
}