package office

import (
	"time"
)

type Office struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;unique" json:"name" validate:"required,min=2,max=100"`
	Location  string    `gorm:"size:255" json:"location"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
