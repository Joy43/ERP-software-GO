package role

import (
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/permission"
)

type Role struct {
	ID          uint                    `gorm:"primaryKey" json:"id"`
	Name        string                  `gorm:"size:150;not null" json:"name"`
	Slug        string                  `gorm:"size:150;uniqueIndex;not null" json:"slug"`
	Description string                  `gorm:"type:text" json:"description"`
	IsActive    bool                    `gorm:"not null;default:true" json:"is_active"`
	UserCount   int64                   `gorm:"column:user_count;->;" json:"user_count"`
	Permissions []permission.Permission `gorm:"many2many:role_permissions;" json:"permissions"`
	CreatedAt   time.Time               `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time               `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}
