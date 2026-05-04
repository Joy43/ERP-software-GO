package sales_supply_type

import (
	"time"
)

type SalesSupplyType struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:150;not null;unique" json:"name" validate:"required,min=2,max=150"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for SalesSupplyType
func (SalesSupplyType) TableName() string {
	return "sales_supply_types"
}