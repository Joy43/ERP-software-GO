package responsibility_transfer

import (
	"time"
)

type ResponsibilityTransfer struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	FromUserID   uint       `gorm:"not null" json:"from_user_id"`
	ToUserID     uint       `gorm:"not null" json:"to_user_id"`
	FromDate     time.Time  `gorm:"type:date;not null" json:"from_date"`
	ToDate       time.Time  `gorm:"type:date;not null" json:"to_date"`
	Remarks      string     `gorm:"type:text" json:"remarks"`
	Status       string     `gorm:"type:enum('Pending','Approved','Rejected');default:'Pending'" json:"status"`
	ApprovedByID *uint      `json:"approved_by_id"`
	ApprovedAt   *time.Time `json:"approved_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relations
	FromUser *User `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUser   *User `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
	ApprovedBy *User `gorm:"foreignKey:ApprovedByID" json:"approved_by,omitempty"`
}

// User is a minimal representation of user for relations to avoid circular dependencies
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (User) TableName() string {
	return "users"
}

type CreateTransferRequest struct {
	FromUserID uint      `json:"from_user_id" binding:"required"`
	ToUserID   uint      `json:"to_user_id" binding:"required"`
	FromDate   time.Time `json:"from_date" binding:"required"`
	ToDate     time.Time `json:"to_date" binding:"required"`
	Remarks    string    `json:"remarks"`
}

type ApproveTransferRequest struct {
	Status string `json:"status" binding:"required,oneof=Approved Rejected"`
}

type TransferFilter struct {
	FromUserID uint       `form:"from_user_id"`
	ToUserID   uint       `form:"to_user_id"`
	FromDate   *time.Time `form:"from_date"`
	ToDate     *time.Time `form:"to_date"`
}
