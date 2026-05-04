package wallet

import (
	"time"
)

type Wallet struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"size:100;not null;unique" json:"name" validate:"required,min=2,max=100"`
	CommissionPercent float64   `gorm:"type:decimal(5,2);default:0.00" json:"commission_percent"`
	BankAccount       string    `gorm:"size:50" json:"bank_account"`
	ReferenceCode     string    `gorm:"size:50" json:"reference_code"`
	IsRounding        bool      `gorm:"default:false" json:"is_rounding"`
	IsCoupon          bool      `gorm:"default:false" json:"is_coupon"`
	IsWalletCharge    bool      `gorm:"default:false" json:"is_wallet_charge"`
	IsActive          bool      `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CreateWalletRequest struct {
	Name              string  `json:"name" binding:"required,min=2,max=100"`
	CommissionPercent float64 `json:"commission_percent"`
	BankAccount       string  `json:"bank_account"`
	ReferenceCode     string  `json:"reference_code"`
	IsRounding        bool    `json:"is_rounding"`
	IsCoupon          bool    `json:"is_coupon"`
	IsWalletCharge    bool    `json:"is_wallet_charge"`
	IsActive          *bool   `json:"is_active"` // Use pointer to handle boolean false in binding
}

type UpdateWalletRequest struct {
	Name              string   `json:"name" binding:"omitempty,min=2,max=100"`
	CommissionPercent *float64 `json:"commission_percent"`
	BankAccount       *string  `json:"bank_account"`
	ReferenceCode     *string  `json:"reference_code"`
	IsRounding        *bool    `json:"is_rounding"`
	IsCoupon          *bool    `json:"is_coupon"`
	IsWalletCharge    *bool    `json:"is_wallet_charge"`
	IsActive          *bool    `json:"is_active"`
}
