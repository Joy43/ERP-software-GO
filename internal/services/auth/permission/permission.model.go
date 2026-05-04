package permission

import "time"

type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	GroupName   string    `gorm:"size:100;not null" json:"group_name"`
	Name        string    `gorm:"size:150;not null" json:"name"`
	Slug        string    `gorm:"size:150;uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
