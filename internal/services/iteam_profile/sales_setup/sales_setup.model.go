package sales_setup

import (
	"time"
)

type SalesSetup struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:150;not null;unique" json:"name" validate:"required,min=2,max=150"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for SalesSetup
func (SalesSetup) TableName() string {
	return "sales_setups"
}
