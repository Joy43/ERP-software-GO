package sub_category

import (
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/category"
)

type SubCategory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CategoryID uint      `gorm:"not null;index" json:"category_id" validate:"required"`
	Name       string    `gorm:"size:150;not null" json:"name" validate:"required,min=2,max=150"`
	Category   *category.Category `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for SubCategory
func (SubCategory) TableName() string {
	return "sub_categories"
}
