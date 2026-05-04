package purchasereturn
import "time"

// ========================================
// REQUEST DTOs
// ========================================

type CreatePurchaseReturnRequest struct {
	Quantity     float64 `json:"quantity"      binding:"required,gt=0"`
	SellingPrice float64 `json:"selling_price" binding:"required,gt=0"`
	Remarks      string  `json:"remarks"`

	OfficeID   uint `json:"office_id"   binding:"required,gt=0"`
	LocationID uint `json:"location_id" binding:"required,gt=0"`
	SupplierID uint `json:"supplier_id" binding:"required,gt=0"`
	ItemID     uint `json:"item_id"     binding:"required,gt=0"`
}

type ListPurchaseReturnRequest struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
	SupplierID *uint  `form:"supplier_id"`
	OfficeID   *uint  `form:"office_id"`
	LocationID *uint  `form:"location_id"`
	ItemID     *uint  `form:"item_id"`
	Search     string `form:"search"`
	SortBy     string `form:"sort_by,default=created_at"`
	SortOrder  string `form:"sort_order,default=desc"`
}

// ========================================
// ----RESPONSE DTOs----
// ========================================

type PurchaseReturnResponse struct {
	ID           uint    `json:"id"`
	ReturnNumber string  `json:"return_number"`
	Quantity     float64 `json:"quantity"`
	SellingPrice float64 `json:"selling_price"`
	Remarks      string  `json:"remarks"`

	OfficeID   uint         `json:"office_id"`
	Office     *OfficeDTO   `json:"office,omitempty"`
	LocationID uint         `json:"location_id"`
	Location   *LocationDTO `json:"location,omitempty"`
	SupplierID uint         `json:"supplier_id"`
	Supplier   *SupplierDTO `json:"supplier,omitempty"`
	ItemID     uint         `json:"item_id"`
	Item       *ItemDTO     `json:"item,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginatedPurchaseReturnResponse struct {
	Data       []PurchaseReturnResponse `json:"data"`
	Pagination PaginationMeta           `json:"pagination"`
}

type PaginationMeta struct {
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// ========================================
// SHARED DTOs
// ========================================

type OfficeDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type LocationDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	
}

type SupplierDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ItemDTO struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	SKU     string `json:"sku,omitempty"`
	Barcode string `json:"barcode,omitempty"`
}