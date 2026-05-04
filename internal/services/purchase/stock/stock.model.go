package stock

import (
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/items"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/user"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/location"
)

// =============================================
// STOCK TRANSACTION MODEL
// =============================================

type TransactionType string

const (
	TxGRN        TransactionType = "GRN"
	TxIssue      TransactionType = "ISSUE"
	TxTransfer   TransactionType = "TRANSFER"
	TxAdjustment TransactionType = "ADJUSTMENT"
	TxReturn     TransactionType = "RETURN"
	TxDamage     TransactionType = "DAMAGE"
	TxExpired    TransactionType = "EXPIRED"
)

type StockTransaction struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	TransactionNumber string          `gorm:"size:50;uniqueIndex;not null" json:"transaction_number"`
	TransactionType   TransactionType `gorm:"type:enum('GRN','ISSUE','TRANSFER','ADJUSTMENT','RETURN','DAMAGE','EXPIRED');not null;index" json:"transaction_type"`

	ItemID     uint             `gorm:"not null;index" json:"item_id"`
	Item       items.Items      `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"item,omitempty"`
	LocationID uint             `gorm:"not null;index" json:"location_id"`
	Location   location.Location `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"location,omitempty"`

	QuantityChange float64 `gorm:"type:decimal(15,2);not null" json:"quantity_change"` // + add, - remove
	BeforeQuantity float64 `gorm:"type:decimal(15,2)" json:"before_quantity"`
	AfterQuantity  float64 `gorm:"type:decimal(15,2)" json:"after_quantity"`
	UnitCost       float64 `gorm:"type:decimal(15,2)" json:"unit_cost"`

	ReferenceType string `gorm:"size:50;index" json:"reference_type"`
	ReferenceID   uint   `gorm:"index" json:"reference_id"`
	GRNItemID     *uint  `gorm:"index" json:"grn_item_id,omitempty"`

	Remarks *string `gorm:"type:text" json:"remarks,omitempty"`

	CreatedByID *uint      `gorm:"index" json:"created_by_id,omitempty"`
	CreatedBy   *user.User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"created_by,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (StockTransaction) TableName() string { return "stock_transactions" }
