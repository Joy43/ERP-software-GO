package minor_category

import (
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/sub_category"
)

type MinorCategory struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	SubCategoryID uint      `gorm:"not null;index" json:"sub_category_id" validate:"required"`
	Name          string    `gorm:"size:150;not null" json:"name" validate:"required,min=2,max=150"`
	SubCategory   *sub_category.SubCategory `gorm:"foreignKey:SubCategoryID;references:ID" json:"sub_category"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for MinorCategory
func (MinorCategory) TableName() string {
	return "minor_categories"
}
