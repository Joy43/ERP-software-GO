package inventorytypes

import (
	"time"

	"gorm.io/gorm"
)

type InventoryType struct {
    ID          uint  `gorm:"primaryKey" json:"id"`
    TypeCode    string         `json:"type_code" gorm:"type:varchar(150);uniqueIndex;not null"`
	TypeName    string         `json:"type_name" gorm:"type:varchar(100);not null"`
	Description string         `json:"description" gorm:"type:text"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (InventoryType) TableName() string {
	return "inventory_types"
}

	