package items

// PaginationParams defines pagination query parameters
type PaginationParams struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
	Sort  string `form:"sort" binding:"omitempty"`
}

// CreateItemRequest defines the request payload for creating a new item
type CreateItemRequest struct {
	// Basic Fields
	Name                string  `json:"name" binding:"required,min=2,max=150" example:"Product Name required"`
	Barcode             string  `json:"barcode" binding:"omitempty,max=100" example:"EAN123456"`
	SKU                 string  `json:"sku" binding:"required,max=100" example:"SKU123"`

	// Category Fields
	CategoryID          uint    `json:"category_id" binding:"required" example:"1"`
	SubCategoryID       uint    `json:"sub_category_id" binding:"required" example:"1"`
	MinorCategoryID     uint    `json:"minor_category_id" binding:"required" example:"1"`
	ItemTypeID          uint    `json:"item_type_id" binding:"required" example:"1"`
	DepartmentID        uint    `json:"department_id" binding:"required" example:"1"`

	// UOM & Supplier
	UomID               uint    `json:"uom_id" binding:"required" example:"1"`
	SupplierID          uint    `json:"supplier_id" binding:"required" example:"1"`
	TagID               uint    `json:"tag_id" binding:"omitempty" example:"1"`

	// Pricing Fields
	CostPrice           float64 `json:"cost_price" binding:"required,min=0" example:"100.00"`
	StandardSalesPrice  float64 `json:"standard_sales_price" binding:"required,min=0" example:"200.00"`
	CostOnGP            float64 `json:"cost_on_gp" binding:"min=0" example:"0"`
	LastCost            float64 `json:"last_cost" binding:"min=0" example:"100.00"`
	AvgCost             float64 `json:"avg_cost" binding:"min=0" example:"100.00"`
	PriceIncreasing     float64 `json:"price_increasing" binding:"min=0" example:"0"`
	AltUMO              float64 `json:"alt_umo" binding:"min=0" example:"1"`
	AltUnit             string  `json:"alt_unit" binding:"max=50" example:"BOX"`
	MaxDiscount         float64 `json:"max_discount" binding:"min=0,max=100" example:"20"`
	ReorderMinimumQty   float64 `json:"reorder_minimum_qty" binding:"min=0" example:"10"`
	CalculateBasePrice  float64 `json:"calculate_base_price" binding:"min=0" example:"100.00"`

	// Tax & Supply Setup
	SalesSupplyTypeID   uint    `json:"sales_supply_type_id" binding:"required" example:"1"`
	SalesTaxSetupID     uint    `json:"sales_tax_setup_id" example:"1"`
	SalesSetupID        uint    `json:"sales_setup_id" binding:"omitempty" example:"1"`

	// Media & Metadata
	FileID              *uint   `json:"file_id" example:"1"`
	

	// Flags
	IsChildBarcode      bool    `json:"is_child_barcode" example:"false"`
	AutoBarcode         bool    `json:"auto_barcode" example:"false"`
	CanBeSold           bool    `json:"can_be_sold" example:"true"`
	CanBeProduced       bool    `json:"can_be_produced" example:"false"`
	CanBeRented         bool    `json:"can_be_rented" example:"false"`
	CanBePurchased      bool    `json:"can_be_purchased" example:"true"`
	IsVATRebatable      bool    `json:"is_vat_rebatable" example:"false"`
	IsNotAllowDecimal   bool    `json:"is_not_allow_decimal" example:"false"`
	IsActive            bool    `json:"is_active" example:"true"`
	IsStyle             bool    `json:"is_style" example:"false"`
	IsPercentage        bool    `json:"is_percentage" example:"false"`
	IsPharmacy          bool    `json:"is_pharmacy" example:"false"`

	// Pharmacy Fields (required if IsPharmacy is true)
	Pharmacy            *CreateItemPharmacyRequest `json:"pharmacy" binding:"omitempty"`
}

// CreateItemPharmacyRequest defines pharmacy-specific fields for items
type CreateItemPharmacyRequest struct {
	GenericName           string `json:"generic_name" binding:"required_if=IsPharmacy true,max=200" example:"Paracetamol"`
	BrandName             string `json:"brand_name" binding:"omitempty,max=200" example:"Napa"`
	Strength              string `json:"strength" binding:"omitempty,max=100" example:"500mg"`
	DosageForm            string `json:"dosage_form" binding:"omitempty,max=100" example:"Tablet"`
	ScheduleType          string `json:"schedule_type" binding:"omitempty,max=50" example:"Schedule I"`
	IsPrescriptionRequired bool  `json:"is_prescription_required" example:"false"`
	IsControlledDrug      bool  `json:"is_controlled_drug" example:"false"`
	StorageCondition      string `json:"storage_condition" binding:"omitempty,max=200" example:"Room Temperature"`
	MaxDailyDose          string `json:"max_daily_dose" binding:"omitempty,max=100" example:"4000mg"`
	ShelfLifeDays         *int  `json:"shelf_life_days" binding:"omitempty,gte=0"`
	ReorderAlertDays      *int  `json:"reorder_alert_days" binding:"omitempty,gte=0"`
	ManufacturerName      string `json:"manufacturer_name" binding:"omitempty,max=200" example:"Beximco Pharmaceuticals"`
	DrugRegistrationNo    string `json:"drug_registration_no" binding:"omitempty,max=100" example:"REG-123456"`
	RouteOfAdministration string `json:"route_of_administration" binding:"omitempty,max=100" example:"Oral"`
	TherapeuticClass      string `json:"therapeutic_class" binding:"omitempty,max=150" example:"Analgesic"`
}

// UpdateItemRequest defines the request payload for updating an item
type UpdateItemRequest struct {
	// Basic Fields
	Name                *string  `json:"name" binding:"omitempty,min=2,max=150"`
	Barcode             *string  `json:"barcode" binding:"omitempty,max=100"`
	SKU                 *string  `json:"sku" binding:"omitempty,max=100"`
	// Category Fields
	CategoryID          *uint    `json:"category_id"`
	SubCategoryID       *uint    `json:"sub_category_id"`
	MinorCategoryID     *uint    `json:"minor_category_id"`
	ItemTypeID          *uint    `json:"item_type_id"`
	DepartmentID        *uint    `json:"department_id"`

	//---------- UOM & Supplier----------
	UomID               *uint    `json:"uom_id"`
	SupplierID          *uint    `json:"supplier_id"`
	TagID               *uint    `json:"tag_id"`

	//-------------- Pricing Fields----------------
	CostPrice           *float64 `json:"cost_price" binding:"omitempty,min=0"`
	StandardSalesPrice  *float64 `json:"standard_sales_price" binding:"omitempty,min=0"`
	CostOnGP            *float64 `json:"cost_on_gp" binding:"omitempty,min=0"`
	LastCost            *float64 `json:"last_cost" binding:"omitempty,min=0"`
	AvgCost             *float64 `json:"avg_cost" binding:"omitempty,min=0"`
	PriceIncreasing     *float64 `json:"price_increasing" binding:"omitempty,min=0"`
	AltUMO              *float64 `json:"alt_umo" binding:"omitempty,min=0"`
	AltUnit             *string  `json:"alt_unit" binding:"omitempty,max=50"`
	MaxDiscount         *float64 `json:"max_discount" binding:"omitempty,min=0,max=100"`
	ReorderMinimumQty   *float64 `json:"reorder_minimum_qty" binding:"omitempty,min=0"`
	CalculateBasePrice  *float64 `json:"calculate_base_price" binding:"omitempty,min=0"`

	//------------ Tax & Supply Setup----------------
	SalesSupplyTypeID   *uint    `json:"sales_supply_type_id"`
	SalesTaxSetupID     *uint    `json:"sales_tax_setup_id"`
	SalesSetupID        *uint    `json:"sales_setup_id"`

	//-------------- Media & Metadata---------------------
	FileID              *uint    `json:"file_id"`


	// --------------Flags----------------------
	IsChildBarcode      *bool    `json:"is_child_barcode"`
	AutoBarcode         *bool    `json:"auto_barcode"`
	CanBeSold           *bool    `json:"can_be_sold"`
	CanBeProduced       *bool    `json:"can_be_produced"`
	CanBeRented         *bool    `json:"can_be_rented"`
	CanBePurchased      *bool    `json:"can_be_purchased"`
	IsVATRebatable      *bool    `json:"is_vat_rebatable"`
	IsNotAllowDecimal   *bool    `json:"is_not_allow_decimal"`
	IsActive            *bool    `json:"is_active"`
	IsStyle             *bool    `json:"is_style"`
	IsPercentage        *bool    `json:"is_percentage"`
}

// --------- ItemResponse defines the response payload for item -------------
type ItemResponse struct {
	ID                  uint                 `json:"id"`
	Name                string               `json:"name"`
	Barcode             string               `json:"barcode"`
	SKU                 string               `json:"sku"`
	CategoryID          uint                 `json:"category_id"`
	Category            interface{}          `json:"category,omitempty"`
	SubCategoryID       uint                 `json:"sub_category_id"`
	SubCategory         interface{}          `json:"sub_category,omitempty"`
	MinorCategoryID     uint                 `json:"minor_category_id"`
	MinorCategory       interface{}          `json:"minor_category,omitempty"`
	ItemTypeID          uint                 `json:"item_type_id"`
	ItemType            interface{}          `json:"item_type,omitempty"`
	DepartmentID        uint                 `json:"department_id"`
	Department          interface{}          `json:"department,omitempty"`
	UomID               uint                 `json:"uom_id"`
	Uom                 interface{}          `json:"uom,omitempty"`
	SupplierID          uint                 `json:"supplier_id"`
	Supplier            interface{}          `json:"supplier,omitempty"`
	TagID               uint                 `json:"tag_id"`
	Tag                 interface{}          `json:"tag,omitempty"`
	SalesSupplyTypeID   uint                 `json:"sales_supply_type_id"`
	SalesSupplyType     interface{}          `json:"sales_supply_type,omitempty"`
	SalesTaxSetupID     uint                 `json:"sales_tax_setup_id"`
	SalesTaxSetup       interface{}          `json:"sales_tax_setup,omitempty"`
	SalesSetupID        uint                 `json:"sales_setup_id"`
	SalesSetup          interface{}          `json:"sales_setup,omitempty"`
	FileID              *uint                `json:"file_id"`
	File                interface{}          `json:"file,omitempty"`
	CostPrice           float64              `json:"cost_price"`
	StandardSalesPrice  float64              `json:"standard_sales_price"`
	CostOnGP            float64              `json:"cost_on_gp"`
	LastCost            float64              `json:"last_cost"`
	AvgCost             float64              `json:"avg_cost"`
	PriceIncreasing     float64              `json:"price_increasing"`
	AltUMO              float64              `json:"alt_umo"`
	AltUnit             string               `json:"alt_unit"`
	MaxDiscount         float64              `json:"max_discount"`
	ReorderMinimumQty   float64              `json:"reorder_minimum_qty"`
	CalculateBasePrice  float64              `json:"calculate_base_price"`
	IsChildBarcode      bool                 `json:"is_child_barcode"`
	AutoBarcode         bool                 `json:"auto_barcode"`
	CanBeSold           bool                 `json:"can_be_sold"`
	CanBeProduced       bool                 `json:"can_be_produced"`
	CanBeRented         bool                 `json:"can_be_rented"`
	CanBePurchased      bool                 `json:"can_be_purchased"`
	IsVATRebatable      bool                 `json:"is_vat_rebatable"`
	IsNotAllowDecimal   bool                 `json:"is_not_allow_decimal"`
	IsActive            bool                 `json:"is_active"`
	IsStyle             bool                 `json:"is_style"`
	IsPercentage        bool                 `json:"is_percentage"`
	IsPharmacy          bool                 `json:"is_pharmacy"`
	Pharmacy            *ItemPharmacyResponse `json:"pharmacy,omitempty"`
	CreatedAt           string               `json:"created_at"`
	UpdatedAt           string               `json:"updated_at"`
}

// ItemPharmacyResponse defines the response structure for pharmacy data
type ItemPharmacyResponse struct {
	ItemID                 uint   `json:"item_id"`
	GenericName            string `json:"generic_name"`
	BrandName              string `json:"brand_name"`
	Strength               string `json:"strength"`
	DosageForm             string `json:"dosage_form"`
	ScheduleType           string `json:"schedule_type"`
	IsPrescriptionRequired bool   `json:"is_prescription_required"`
	IsControlledDrug       bool   `json:"is_controlled_drug"`
	StorageCondition       string `json:"storage_condition"`
	MaxDailyDose           string `json:"max_daily_dose"`
	ShelfLifeDays          *int   `json:"shelf_life_days"`
	ReorderAlertDays       *int   `json:"reorder_alert_days"`
	ManufacturerName       string `json:"manufacturer_name"`
	DrugRegistrationNo     string `json:"drug_registration_no"`
	RouteOfAdministration  string `json:"route_of_administration"`
	TherapeuticClass       string `json:"therapeutic_class"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
}

// ListResponse defines the response for list items with pagination
type ListResponse struct {
	Items      []ItemResponse `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}