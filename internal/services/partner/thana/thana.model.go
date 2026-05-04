package thana

import (
	"time"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/district"
)

type Thana struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	DistrictID uint              `gorm:"not null" json:"district_id"`
	District   district.District `gorm:"foreignKey:DistrictID" json:"district"`
	Name       string            `gorm:"size:150;not null" json:"name" validate:"required,min=2,max=150"`
	CreatedAt  time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}
