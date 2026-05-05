package transfer

import (
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/office"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/items"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/location"
	"gorm.io/gorm"
)

type Transfer struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// =========================
	// BASIC INFO
	// =========================
	ReferenceNo     string  `gorm:"type:varchar(100);unique" json:"reference_no"`
	TransferQty     float64 `json:"transfer_qty"`
	Status          string  `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft, approved, completed
	GeneralRemarks  string  `gorm:"type:text" json:"general_remarks"`
	Note            string  `gorm:"type:text" json:"note"`

	// =========================
	// RELATIONS
	// =========================
	ItemID uint        `json:"item_id"`
	Item   items.Items `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	FromOfficeID uint          `json:"from_office_id"`
	FromOffice   office.Office `gorm:"foreignKey:FromOfficeID"`

	ToOfficeID uint          `json:"to_office_id"`
	ToOffice   office.Office `gorm:"foreignKey:ToOfficeID"`

	ToLocationID uint              `json:"to_location_id"`
	ToLocation   location.Location `gorm:"foreignKey:ToLocationID"`

	// =========================
	// AUDIT
	// =========================
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CreatedByID uint `json:"created_by"`
	UpdatedByID uint `json:"updated_by"`
}