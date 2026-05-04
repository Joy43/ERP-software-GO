package locationstock

import (
	"time"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/items"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/location"
	"gorm.io/gorm"
)
type LocationStock struct {
	ID uint `gorm:"primaryKey"`

	// =============================
	//------- STOCK DATA-------
	// =============================
	Quantity         float64 `gorm:"type:decimal(15,3);not null;default:0"`
	ReservedQuantity float64 `gorm:"type:decimal(15,3);default:0"`
	LastCost         float64 `gorm:"type:decimal(15,4);default:0"`
	AverageCost      float64 `gorm:"type:decimal(15,4);default:0"`

	// =============================
	//------ RELATIONS ------
	// =============================
	ItemID uint        `gorm:"not null;index;uniqueIndex:idx_item_location"`
	Item   items.Items `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	LocationID uint              `gorm:"not null;index;uniqueIndex:idx_item_location"`
	Location   location.Location `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// =============================
	// TIMESTAMPS
	// =============================
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}



