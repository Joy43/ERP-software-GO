package category

import (
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/department"
)

type Category struct {
	ID           uint                        `gorm:"primaryKey" json:"id"`
	Name         string                      `gorm:"size:150;not null;unique" json:"name" validate:"required,min=2,max=150"`
	DepartmentID *uint                       `json:"department_id"`
	Department   *department.Department      `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	CreatedAt    time.Time                   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time                   `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Category) TableName() string {
	return "categories"
}
