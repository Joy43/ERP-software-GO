package location

import (
	"time"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/office"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/user"

	
)
type Location struct {
	ID uint `gorm:"primaryKey"`

	// ====================
	// ------BASIC INFO----
	// ====================
	Code string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Name string `gorm:"type:varchar(200);not null"`

	Type LocationType `gorm:"type:enum('warehouse','store','showroom','outlet');not null"`

	// =============================
	//----------- RELATIONS ---------
	// =============================
	OfficeID uint `gorm:"not null" json:"office_id"`
	Office  office.Office `gorm:"foreignKey:OfficeID"`
	ParentID *uint     `gorm:"index" json:"parent_id"`
	Parent   *Location `gorm:"foreignKey:ParentID" `
	Children []Location `gorm:"foreignKey:ParentID"`

	// =============================
	// OTHER INFO
	// =============================
	Location  string `gorm:"type:text"`
	IsActive  bool   `gorm:"default:true"`
	ManagerID *uint `gorm:"index" json:"manager_id"`
	Manager   *user.User `gorm:"foreignKey:ManagerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// =============================
	// TIMESTAMPS
	// =============================
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
