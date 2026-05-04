package sales_tax_setup

import (
	"time"
)

type SalesTaxSetup struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:150;not null;unique" json:"name" validate:"required,min=2,max=150"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for SalesTaxSetup
func (SalesTaxSetup) TableName() string {
	return "sales_tax_setups"
}