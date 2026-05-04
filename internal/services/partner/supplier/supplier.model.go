package supplier

import (
	"time"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/office"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/district"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/partner_group"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/partner_sub_group"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/tax_bracket"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/thana"
	"gorm.io/gorm"
)

type Supplier struct {
	ID                            uint                               `gorm:"primaryKey" json:"id"`
	OfficeID                      uint                               `gorm:"not null" json:"office_id"`
	Office                        office.Office                     `gorm:"foreignKey:OfficeID" json:"office"`
	Code                          string                             `gorm:"size:50;unique" json:"code"`
	Name                          string                             `gorm:"size:255;not null" json:"name"`
	VendorName                    string                             `gorm:"size:255" json:"vendor_name"`
	MobileNo                      string                             `gorm:"size:20;not null" json:"mobile_no"`
	Email                         string                             `gorm:"size:190" json:"email"`
	CreditDays                    int                                `gorm:"default:0" json:"credit_days"`
	PartnerGroupID               *uint                              `json:"partner_group_id"`
	PartnerGroup                 *partner_group.PartnerGroup       `gorm:"foreignKey:PartnerGroupID" json:"partner_group"`
	PartnerSubGroupID            *uint                              `json:"partner_sub_group_id"`
	PartnerSubGroup              *partner_sub_group.PartnerSubGroup `gorm:"foreignKey:PartnerSubGroupID" json:"partner_sub_group"`
	Tolerance                     float64                            `gorm:"type:decimal(5,2);default:0.00" json:"tolerance"`
	BinNumber                     string                             `gorm:"size:100" json:"bin_number"`
	TCSPercentage                 float64                            `gorm:"type:decimal(5,2);default:0.00" json:"tcs_percentage"`
	VATRegNoCentral               string                             `gorm:"size:100" json:"vat_reg_no_central"`
	TINNumber                     string                             `gorm:"size:100" json:"tin_number"`
	TradeLicenseNumber            string                             `gorm:"size:100" json:"trade_license_number"`
	IsForeignSupplier             bool                               `gorm:"default:false" json:"is_foreign_supplier"`
	VDSApplicable                 bool                               `gorm:"default:false" json:"vds_applicable"`
	IsAnonymous                   bool                               `gorm:"default:false" json:"is_anonymous"`
	IsVATAccountingNotApplicable bool                               `gorm:"default:false" json:"is_vat_accounting_not_applicable"`
	
	// Billing Details
	BillingContactPerson          string                             `gorm:"size:150" json:"billing_contact_person"`
	BillingContactNo              string                             `gorm:"size:20" json:"billing_contact_no"`
	DistrictID                    *uint                              `json:"district_id"`
	District                      *district.District                 `gorm:"foreignKey:DistrictID" json:"district"`
	ThanaID                       *uint                              `json:"thana_id"`
	Thana                         *thana.Thana                       `gorm:"foreignKey:ThanaID" json:"thana"`
	BillingAddress                string                             `gorm:"type:text" json:"billing_address"`
	TaxBracketID                  *uint                              `json:"tax_bracket_id"`
	TaxBracket                    *tax_bracket.TaxBracket           `gorm:"foreignKey:TaxBracketID" json:"tax_bracket"`
	Remarks                       string                             `gorm:"type:text" json:"remarks"`
	Attachment                    string                             `gorm:"size:255" json:"attachment"`
	
	// Extra Supplier Details
	Address                       string                             `gorm:"type:text" json:"address"`
	GPS                           string                             `gorm:"size:255" json:"gps"`
	ContactPerson                 string                             `gorm:"size:150" json:"contact_person"`
	VATRegNo                      string                             `gorm:"size:100" json:"vat_reg_no"`
	
	CreatedAt                     time.Time                          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt                     time.Time                          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt                     gorm.DeletedAt                     `gorm:"index" json:"-"`
}
