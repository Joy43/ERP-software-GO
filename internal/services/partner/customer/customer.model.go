package customer

import (
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/office"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/user"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/district"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/partner_group"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/partner_sub_group"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/tax_bracket"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/thana"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/price_type"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/sales_representative"
	"gorm.io/gorm"
)

type Customer struct {
	ID                        uint                               `gorm:"primaryKey" json:"id"`
	OfficeID                  uint                               `gorm:"not null" json:"office_id"`
	Office                    office.Office                     `gorm:"foreignKey:OfficeID" json:"office"`
	Name                      string                             `gorm:"size:255;not null" json:"name"`
	BillingName               string                             `gorm:"size:255;not null" json:"billing_name"`
	MobileNo                  string                             `gorm:"size:20;not null" json:"mobile_no"`
	Email                     string                             `gorm:"size:190" json:"email"`
	CreditDays                int                                `gorm:"default:0" json:"credit_days"`
	MaritalStatus             string                             `gorm:"size:50" json:"marital_status"`
	MarriageDate              *time.Time                         `json:"marriage_date"`
	BirthdayDate              *time.Time                         `json:"birthday_date"`
	Gender                    string                             `gorm:"size:20" json:"gender"`
	NID                       string                             `gorm:"column:nid;size:50" json:"nid"`
	PartnerGroupID           *uint                              `json:"partner_group_id"`
	PartnerGroup             *partner_group.PartnerGroup       `gorm:"foreignKey:PartnerGroupID" json:"partner_group"`
	PartnerSubGroupID        *uint                              `json:"partner_sub_group_id"`
	PartnerSubGroup          *partner_sub_group.PartnerSubGroup `gorm:"foreignKey:PartnerSubGroupID" json:"partner_sub_group"`
	CreditLimit              float64                            `gorm:"type:decimal(15,2);default:0.00" json:"credit_limit"`
	DefaultSalesRepID        *uint                                            `json:"default_sales_rep_id"`
	DefaultSalesRep          *sales_representative.SalesRepresentative        `gorm:"foreignKey:DefaultSalesRepID" json:"default_sales_rep"`
	Tolerance                float64                            `gorm:"type:decimal(5,2);default:0.00" json:"tolerance"`
	BinNumber                string                             `gorm:"size:100" json:"bin_number"`
	PriceTypeID              *uint                                            `json:"price_type_id"`
	PriceType                *price_type.PriceType                            `gorm:"foreignKey:PriceTypeID" json:"price_type"`
	TCSPercentage            float64                            `gorm:"type:decimal(5,2);default:0.00" json:"tcs_percentage"`
	VATRegNoCentral           string                             `gorm:"size:100" json:"vat_reg_no_central"`
	TINNumber                string                             `gorm:"size:100" json:"tin_number"`
	TradeLicenseNumber        string                             `gorm:"size:100" json:"trade_license_number"`
	IsForeignCustomer         bool                               `gorm:"default:false" json:"is_foreign_customer"`
	IsDistributor             bool                               `gorm:"default:false" json:"is_distributor"`
	IsEmployeeCustomer        bool                               `gorm:"default:false" json:"is_employee_customer"`
	UserID                    *uint                              `json:"user_id"`
	User                      *user.User                         `gorm:"foreignKey:UserID" json:"user_customer"`
	
	// Billing Details
	BillingContactPerson      string                             `gorm:"size:150" json:"billing_contact_person"`
	BillingContactNo          string                             `gorm:"size:20" json:"billing_contact_no"`
	DistrictID                *uint                              `json:"district_id"`
	District                  *district.District                 `gorm:"foreignKey:DistrictID" json:"district"`
	ThanaID                   *uint                              `json:"thana_id"`
	Thana                     *thana.Thana                       `gorm:"foreignKey:ThanaID" json:"thana"`
	BillingAddress            string                             `gorm:"type:text;not null" json:"billing_address"`
	PointExpiryDays           int                                `gorm:"default:0" json:"point_expiry_days"`
	TaxBracketID              *uint                              `json:"tax_bracket_id"`
	TaxBracket                *tax_bracket.TaxBracket           `gorm:"foreignKey:TaxBracketID" json:"tax_bracket"`
	Remarks                   string                             `gorm:"type:text" json:"remarks"`
	Attachment                string                             `gorm:"size:255" json:"attachment"`
	
	// Relationships
	ShippingAddresses         []CustomerShippingAddress          `gorm:"foreignKey:CustomerID" json:"shipping_addresses"`
	BankInfos                 []CustomerBankInfo                 `gorm:"foreignKey:CustomerID" json:"bank_infos"`
	
	CreatedAt                 time.Time                          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt                 time.Time                          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt                 gorm.DeletedAt                     `gorm:"index" json:"-"`
}

type CustomerShippingAddress struct {
	ID             uint               `gorm:"primaryKey" json:"id"`
	CustomerID     uint               `gorm:"not null" json:"customer_id"`
	Code           string             `gorm:"size:50" json:"code"`
	RefOutlet      string             `gorm:"size:150" json:"ref_outlet"`
	Address        string             `gorm:"type:text;not null" json:"address"`
	GPS            string             `gorm:"size:255" json:"gps"`
	ContactPerson  string             `gorm:"size:150" json:"contact_person"`
	ContactNo      string             `gorm:"size:20" json:"contact_no"`
	VATRegNo       string             `gorm:"size:100" json:"vat_reg_no"`
	DistrictID     *uint              `json:"district_id"`
	District       *district.District `gorm:"foreignKey:DistrictID" json:"district"`
	ThanaID        *uint              `json:"thana_id"`
	Thana          *thana.Thana       `gorm:"foreignKey:ThanaID" json:"thana"`
	CreatedAt      time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
}

type CustomerBankInfo struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CustomerID  uint      `gorm:"not null" json:"customer_id"`
	BankName    string    `gorm:"size:255" json:"bank_name"`
	AccountNo   string    `gorm:"size:100" json:"account_no"`
	AccountName string    `gorm:"size:255" json:"account_name"`
	BranchName  string    `gorm:"size:255" json:"branch_name"`
	RoutingNo   string    `gorm:"size:50" json:"routing_no"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
