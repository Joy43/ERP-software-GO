package items

import (
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/department"
	sales_supply_type "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/SalesSupplyType"
	sales_tax_setup "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/SalesTaxSetup"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/item_type"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/minor_category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/sales_setup"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/sub_category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/tags"
	umo_measurement "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/uom"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/supplier"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/uploads"
	"gorm.io/gorm"
)

type Items struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// ---------------Foreign Keys-------------------------
	CategoryID      uint  `json:"category_id" binding:"required,gt=0"`
	Category        *category.Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	SubCategoryID   uint  `json:"sub_category_id" binding:"required,gt=0"`
	SubCategory     *sub_category.SubCategory `gorm:"foreignKey:SubCategoryID" json:"sub_category,omitempty"`

	MinorCategoryID uint  `json:"minor_category_id" binding:"required,gt=0"`
	MinorCategory   *minor_category.MinorCategory `gorm:"foreignKey:MinorCategoryID" json:"minor_category,omitempty"`

	SalesTaxSetupID   uint `json:"sales_tax_setup_id" binding:"omitempty,gt=0"`
	SalesTaxSetup     *sales_tax_setup.SalesTaxSetup `gorm:"foreignKey:SalesTaxSetupID" json:"sales_tax_setup,omitempty"`

	SalesSupplyTypeID uint `json:"sales_supply_type_id" binding:"omitempty,gt=0"`
	SalesSupplyType   *sales_supply_type.SalesSupplyType `gorm:"foreignKey:SalesSupplyTypeID" json:"sales_supply_type,omitempty"`

	ItemTypeID        uint `json:"item_type_id" binding:"omitempty,gt=0"`
	ItemType          *item_type.ItemType `gorm:"foreignKey:ItemTypeID" json:"item_type,omitempty"`

	DepartmentID      uint `json:"department_id" binding:"omitempty,gt=0"`
	Department        *department.Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`

	UomId             uint `json:"uom_id" binding:"required,gt=0"`
	Uom               *umo_measurement.Uom `gorm:"foreignKey:UomId" json:"uom,omitempty"`

	SupplierID        uint `json:"supplier_id" binding:"omitempty,gt=0"`
	Supplier          *supplier.Supplier `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`

	TagID             uint `json:"tag_id" binding:"omitempty,gt=0"`
	Tag               *tags.Tags `gorm:"foreignKey:TagID" json:"tag,omitempty"`

	SalesSetupID      uint `json:"sales_setup_id" binding:"omitempty,gt=0"`
	SalesSetup        *sales_setup.SalesSetup `gorm:"foreignKey:SalesSetupID" json:"sales_setup,omitempty"`

FileID            *uint `json:"file_id" binding:"omitempty,gt=0"`
File              *uploads.File `gorm:"foreignKey:FileID" json:"file,omitempty"`

	//------------- Item Basic Fields-----------------
	Name     string `gorm:"size:150;not null" json:"name" binding:"required,min=2,max=150"`
	Barcode  string `gorm:"size:400" json:"barcode" binding:"omitempty,max=400"`
	SKU      string `gorm:"size:100;unique" json:"sku" binding:"omitempty,max=100"`
	// -------------Pricing Fields----------------
	CostPrice          float64 `json:"cost_price" binding:"omitempty,gte=0"`
	StandardSalesPrice float64 `json:"standard_sales_price" binding:"omitempty,gte=0"`
	CostOnGP           float64 `json:"cost_on_gp" binding:"omitempty,gte=0"`
	LastCost           float64 `json:"last_cost" binding:"omitempty,gte=0"`
	AvgCost            float64 `json:"avg_cost" binding:"omitempty,gte=0"`
	PriceIncreasing    float64 `json:"price_increasing" binding:"omitempty,gte=0"`
	AltUMO             float64 `json:"alt_umo" binding:"omitempty,gte=0"`
	AltUnit            string  `gorm:"size:50" json:"alt_unit" binding:"omitempty,max=50"`
	MaxDiscount        float64 `json:"max_discount" binding:"omitempty,gte=0,lte=100"`
	ReorderMinimumQty  float64 `json:"reorder_minimum_qty" binding:"omitempty,gte=0"`
	CalculateBasePrice float64 `json:"calculate_base_price" binding:"omitempty,gte=0"`

	//-------------- Boolean Flags---------------
	IsChildBarcode    bool `json:"is_child_barcode" gorm:"default:false"`
	AutoBarcode       bool `json:"auto_barcode" gorm:"default:false"`
	CanBeSold         bool `json:"can_be_sold" gorm:"default:true"`
	CanBeProduced     bool `json:"can_be_produced" gorm:"default:false"`
	CanBeRented       bool `json:"can_be_rented" gorm:"default:false"`
	CanBePurchased    bool `json:"can_be_purchased" gorm:"default:true"`
	IsVATRebatable    bool `json:"is_vat_rebatable" gorm:"default:false"`
	IsNotAllowDecimal bool `json:"is_not_allow_decimal" gorm:"default:false"`
	IsActive          bool `json:"is_active" gorm:"default:true"`
	IsStyle           bool `json:"is_style" gorm:"default:false"`
	IsPercentage      bool `json:"is_percentage" gorm:"default:false"`
	IsPharmacy        bool `json:"is_pharmacy" gorm:"default:false;index"`
	
	//------- Pharmacy relationship------------------
	Pharmacy          *ItemPharmacy `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE" json:"pharmacy,omitempty"`

	//------------ Timestamps---------------------
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

//--------- TableName specifies the table name for Items----------
func (Items) TableName() string {
	return "items"
}

type ItemPharmacy struct {
	ItemID uint   `gorm:"primaryKey" json:"item_id"`
	Item   *Items `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE" json:"item,omitempty"`
	GenericName string `gorm:"size:200;not null" json:"generic_name" binding:"required,max=200"`
	BrandName   string `gorm:"size:200" json:"brand_name" binding:"omitempty,max=200"`
	Strength    string `gorm:"size:100" json:"strength" binding:"omitempty,max=100"`
	DosageForm  string `gorm:"size:100" json:"dosage_form" binding:"omitempty,max=100"`

	ScheduleType           string `gorm:"size:50" json:"schedule_type" binding:"omitempty,max=50"`
	IsPrescriptionRequired bool   `gorm:"default:false" json:"is_prescription_required"`
	IsControlledDrug       bool   `gorm:"default:false" json:"is_controlled_drug"`

	StorageCondition string `gorm:"size:200" json:"storage_condition" binding:"omitempty,max=200"`
	MaxDailyDose     string `gorm:"size:100" json:"max_daily_dose" binding:"omitempty,max=100"`

	ShelfLifeDays    *int `json:"shelf_life_days" binding:"omitempty,gte=0"`
	ReorderAlertDays *int `json:"reorder_alert_days" binding:"omitempty,gte=0"`

	ManufacturerName      string `gorm:"size:200" json:"manufacturer_name" binding:"omitempty,max=200"`
	DrugRegistrationNo    string `gorm:"size:100;uniqueIndex" json:"drug_registration_no" binding:"omitempty,max=100"`
	RouteOfAdministration string `gorm:"size:100" json:"route_of_administration" binding:"omitempty,max=100"`
	TherapeuticClass      string `gorm:"size:150;index" json:"therapeutic_class" binding:"omitempty,max=150"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ItemPharmacy) TableName() string {
    return "item_pharmacies"
}
