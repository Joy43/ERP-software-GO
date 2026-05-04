package grn

import "time"

// =============================================
// REQUEST DTOs
// =============================================

type CreateGRNRequest struct {
	GRNDate     string         `json:"grn_date" binding:"required"`
	ReceiveType GRNReceiveType `json:"receive_type" binding:"required,oneof=DIRECT AGAINST_PO"`

	POID          *uint `json:"po_id"`           
	RequisitionID *uint `json:"requisition_id"`  

	OfficeID   uint `json:"office_id" binding:"required,gt=0"`
	LocationID uint `json:"location_id" binding:"required,gt=0"`
	SupplierID uint `json:"supplier_id" binding:"required,gt=0"`

	ChallanNo              *string `json:"challan_no"`
	ChallanDate            *string `json:"challan_date"`
	SalesInvoiceNumber     *string `json:"sales_invoice_number"`
	VATChallanNumber       *string `json:"vat_challan_number"`
	DeliveryNumber         *string `json:"delivery_number"`
	ShippingAddress        *string `json:"shipping_address"`
	ShipmentDocumentNumber *string `json:"shipment_document_number"`
	PaymentMethodID        *uint   `json:"payment_method_id"`
	Remarks                *string `json:"remarks"`
	FileID                 *uint   `json:"file_id" binding:"omitempty,gt=0"`
	ReceivedByID           *uint   `json:"received_by_id"`

	Items []CreateGRNItemRequest `json:"items" binding:"required,min=1,dive"`
}

type CreateGRNItemRequest struct {
	ItemID    uint  `json:"item_id" binding:"required,gt=0"`
	POItemID  *uint `json:"po_item_id"` 

	ReceivedQuantity   float64 `json:"received_quantity" binding:"required,gt=0"`
	UOMID              *uint   `json:"uom_id"`
	PurchasePrice      float64 `json:"purchase_price" binding:"required,gte=0"`
	VatPercentage      float64 `json:"vat_percentage"`
	DiscountPercentage float64 `json:"discount_percentage"`

	CategoryID      *uint   `json:"category_id"`
	SubCategoryID   *uint   `json:"sub_category_id"`
	MinorCategoryID *uint   `json:"minor_category_id"`
	Remarks         *string `json:"remarks"`
}

type ListGRNRequest struct {
	Page        int     `form:"page,default=1"`
	PageSize    int     `form:"page_size,default=10"`
	Status      *string `form:"status"`
	ReceiveType *string `form:"receive_type"`
	SupplierID  *uint   `form:"supplier_id"`
	POID        *uint   `form:"po_id"`
	Search      string  `form:"search"`
	SortBy      string  `form:"sort_by,default=created_at"`
	SortOrder   string  `form:"sort_order,default=desc"`
}

// =============================================
// RESPONSE DTOs
// =============================================

type GRNResponse struct {
	ID          uint           `json:"id"`
	GRNNumber   string         `json:"grn_number"`
	GRNDate     time.Time      `json:"grn_date"`
	ReceiveType GRNReceiveType `json:"receive_type"`
	Status      GRNStatus      `json:"status"`

	POID          *uint `json:"po_id,omitempty"`
	RequisitionID *uint `json:"requisition_id,omitempty"`

	OfficeID   uint `json:"office_id"`
	LocationID uint `json:"location_id"`
	SupplierID uint `json:"supplier_id"`

	ChallanNo              *string    `json:"challan_no,omitempty"`
	ChallanDate            *time.Time `json:"challan_date,omitempty"`
	SalesInvoiceNumber     *string    `json:"sales_invoice_number,omitempty"`
	VATChallanNumber       *string    `json:"vat_challan_number,omitempty"`
	DeliveryNumber         *string    `json:"delivery_number,omitempty"`
	ShippingAddress        *string    `json:"shipping_address,omitempty"`
	ShipmentDocumentNumber *string    `json:"shipment_document_number,omitempty"`
	PaymentMethodID        *uint      `json:"payment_method_id,omitempty"`
	Remarks                *string    `json:"remarks,omitempty"`
	FileID                 *uint      `json:"file_id,omitempty"`
	File                   *FileDTO   `json:"file,omitempty"`

	Office   *GRNRefDTO `json:"office,omitempty"`
	Location *GRNRefDTO `json:"location,omitempty"`
	Supplier *GRNRefDTO `json:"supplier,omitempty"`

	CreatedByID  *uint `json:"created_by_id,omitempty"`
	ReceivedByID *uint `json:"received_by_id,omitempty"`

	Items []GRNItemResponse `json:"items"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GRNItemResponse struct {
	ID               uint    `json:"id"`
	GRNID            uint    `json:"grn_id"`
	ItemID           uint    `json:"item_id"`
	POItemID         *uint   `json:"po_item_id,omitempty"`
	ReceivedQuantity float64 `json:"received_quantity"`
	UOMID            *uint   `json:"uom_id,omitempty"`
	PurchasePrice    float64 `json:"purchase_price"`

	VatPercentage      float64 `json:"vat_percentage"`
	VatAmount          float64 `json:"vat_amount"`
	DiscountPercentage float64 `json:"discount_percentage"`
	DiscountAmount     float64 `json:"discount_amount"`
	TotalAmount        float64 `json:"total_amount"`

	CategoryID      *uint `json:"category_id,omitempty"`
	SubCategoryID   *uint `json:"sub_category_id,omitempty"`
	MinorCategoryID *uint `json:"minor_category_id,omitempty"`
	Remarks         *string `json:"remarks,omitempty"`

	Item *GRNRefDTO `json:"item,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

type GRNRefDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FileDTO struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	Size         uint64 `json:"size"`
	StoragePath  string `json:"storage_path"`
	Extension    string `json:"extension,omitempty"`
	IsPublic     bool   `json:"is_public"`
}

type PaginatedGRNResponse struct {
	Data       []GRNResponse    `json:"data"`
	Pagination GRNPaginationMeta `json:"pagination"`
}

type GRNPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
