package sales_representative

import (
	"time"
	"gorm.io/gorm"
)

type SalesRepresentative struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:150;not null" json:"name"`
	Mobile         string         `gorm:"size:20;not null" json:"mobile"`
	Email          string         `gorm:"size:190" json:"email"`
	CommissionRate float64        `gorm:"type:decimal(5,2);default:0.00" json:"commission_rate"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
