package special_discount

import "time"
type SpecialDiscount struct {
	ID uint `gorm:"primaryKey"`

	OfferName string `gorm:"size:255;not null"`

	OfficeID  uint `gorm:"not null"`
	CounterID uint `gorm:"not null"`

	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"not null"`

	Remarks string `gorm:"type:text"`

	// Flags boolean fields to indicate various properties of the discount offer
	IsActive             bool `gorm:"default:true"`
	IsOpen               bool `gorm:"default:false"`
	IsItemWiseOpen       bool `gorm:"default:false"`
	IsOnlyDiscount       bool `gorm:"default:false"`
	IsAutoApplied        bool `gorm:"default:false"`
	IsRequiredComment    bool `gorm:"default:false"`
	IsCapped             bool `gorm:"default:false"`

	Attachment string `gorm:"size:700"`

	CreatedAt time.Time
	UpdatedAt time.Time

	
}
