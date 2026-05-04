package projects

import (
	"time"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/office"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/user"
	"gorm.io/gorm"
)

// =============================
// ENUM TYPE
// =============================
type ProjectStatus string

const (
	StatusPlanning  ProjectStatus = "PLANNING"
	StatusActive    ProjectStatus = "ACTIVE"
	StatusOnHold    ProjectStatus = "ON_HOLD"
	StatusCompleted ProjectStatus = "COMPLETED"
	StatusCancelled ProjectStatus = "CANCELLED"
)

// =============================
// --------MODEL----------
// =============================
type Project struct {
	ID uint `gorm:"primaryKey"`

	// =============================
	// BASIC INFO
	// =============================
	ProjectName string  `gorm:"size:255;not null" json:"project_name"`
	Description *string `gorm:"type:text" json:"description,omitempty"`
	ProjectCode string  `gorm:"size:50;not null;uniqueIndex" json:"project_code"`

	// =============================
	// PROJECT DETAILS
	// =============================
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Budget    *float64   `gorm:"type:decimal(15,2)" json:"budget,omitempty"`

	Status ProjectStatus `gorm:"type:enum('PLANNING','ACTIVE','ON_HOLD','COMPLETED','CANCELLED');default:'PLANNING'" json:"status"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	// =============================
	// RELATIONS
	// =============================
	ManagerID *uint      `gorm:"index" json:"manager_id,omitempty"`
	Manager   *user.User `gorm:"foreignKey:ManagerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	OfficeID uint          `gorm:"not null;index" json:"office_id"`
	Office   office.Office `gorm:"foreignKey:OfficeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// =============================
	// TIMESTAMPS
	// =============================
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// =============================
// TABLE NAME
// =============================
func (Project) TableName() string {
	return "projects"
}