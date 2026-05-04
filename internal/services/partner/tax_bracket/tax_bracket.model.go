package tax_bracket

import (
	"time"
)

type TaxBracket struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:150;not null;unique" json:"name" validate:"required,min=2,max=150"`
	Percentage float64   `gorm:"type:decimal(5,2);default:0.00" json:"percentage"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
