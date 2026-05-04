package transaction_history

import (
	"time"
)

type TransactionDeleteHistory struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TransactionType string    `gorm:"type:varchar(100);not null" json:"transaction_type"`
	InvoiceRefNo    string    `gorm:"type:varchar(100);not null" json:"invoice_ref_no"`
	Amount          float64   `gorm:"type:decimal(15,2);default:0.00" json:"amount"`
	DeletedByID     *uint     `json:"deleted_by_id"`
	ReasonRemarks   string    `gorm:"type:text" json:"reason_remarks"`
	DeletedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"deleted_at"`

	// Relations
	DeletedBy *User `gorm:"foreignKey:DeletedByID" json:"deleted_by,omitempty"`
}

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (User) TableName() string {
	return "users"
}

type CreateDeleteHistoryRequest struct {
	TransactionType string  `json:"transaction_type" binding:"required"`
	InvoiceRefNo    string  `json:"invoice_ref_no" binding:"required"`
	Amount          float64 `json:"amount"`
	ReasonRemarks   string  `json:"reason_remarks"`
}

type HistoryFilter struct {
	TransactionType string `form:"transaction_type"`
	DeletedBy       uint   `form:"deleted_by"`
	FromDate        string `form:"from_date"`
	ToDate          string `form:"to_date"`
}
