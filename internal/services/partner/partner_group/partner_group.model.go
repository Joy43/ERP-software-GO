package partner_group

import (
	"time"
)

type PartnerGroup struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:150;not null;unique" json:"name" validate:"required,min=2,max=150"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
