package partner_sub_group

import (
	"time"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/partner_group"
)

type PartnerSubGroup struct {
	ID             uint                        `gorm:"primaryKey" json:"id"`
	PartnerGroupID uint                        `gorm:"not null" json:"partner_group_id"`
	PartnerGroup   partner_group.PartnerGroup `gorm:"foreignKey:PartnerGroupID" json:"partner_group"`
	Name           string                      `gorm:"size:150;not null" json:"name" validate:"required,min=2,max=150"`
	CreatedAt      time.Time                   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time                   `gorm:"autoUpdateTime" json:"updated_at"`
}
