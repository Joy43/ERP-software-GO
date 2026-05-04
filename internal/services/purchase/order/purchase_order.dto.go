package purchaseorder

import "time"

// =============================================
// REQUEST DTOs
// =============================================

type CreatePORequest struct {
	PODate       string      `json:"po_date" binding:"required"`
	DeliveryDate *string     `json:"delivery_date"`
	OrderType    POOrderType `json:"order_type" binding:"required,oneof=DIRECT REQUISITION_BASED CONTRACT"`
	RequisitionID *uint `json:"requisition_id"` 
	OfficeID   uint `json:"office_id" binding:"required,gt=0"`
	LocationID uint `json:"location_id" binding:"required,gt=0"`
	SupplierID uint `json:"supplier_id" binding:"required,gt=0"`
	PaymentTerms    *string `json:"payment_terms"`
	GeneralRemarks  *string `json:"general_remarks"`
	ShippingAddress *string `json:"shipping_address"`

	// Items: required for DIRECT/CONTRACT, auto-populated from requisition for REQUISITION_BASED
	Items []CreatePOItemRequest `json:"items"`
}

type CreatePOItemRequest struct {
	ItemID            uint    `json:"item_id" binding:"required,gt=0"`
	RequisitionItemID *uint   `json:"requisition_item_id"`
	OrderQuantity     float64 `json:"order_quantity" binding:"required,gt=0"`
	UOMID             *uint   `json:"uom_id"`
	UnitPrice         float64 `json:"unit_price" binding:"required,gte=0"`

	VatPercentage      float64 `json:"vat_percentage"`
	DiscountPercentage float64 `json:"discount_percentage"`
	Remarks            *string `json:"remarks"`
}

type UpdatePORequest struct {
	DeliveryDate    *string `json:"delivery_date"`
	PaymentTerms    *string `json:"payment_terms"`
	GeneralRemarks  *string `json:"general_remarks"`
	ShippingAddress *string `json:"shipping_address"`
	Items           []CreatePOItemRequest `json:"items"`
}

type UpdatePOStatusRequest struct {
	Status  POStatus `json:"status" binding:"required"`
	Remarks *string  `json:"remarks"`
}

type ListPORequest struct {
	Page       int     `form:"page,default=1"`
	PageSize   int     `form:"page_size,default=10"`
	Status     *string `form:"status"`
	SupplierID *uint   `form:"supplier_id"`
	OrderType  *string `form:"order_type"`
	Search     string  `form:"search"`
	SortBy     string  `form:"sort_by,default=created_at"`
	SortOrder  string  `form:"sort_order,default=desc"`
}

// =============================================
// RESPONSE DTOs
// =============================================

type POResponse struct {
	ID           uint        `json:"id"`
	PONumber     string      `json:"po_number"`
	PODate       time.Time   `json:"po_date"`
	DeliveryDate *time.Time  `json:"delivery_date,omitempty"`
	OrderType    POOrderType `json:"order_type"`
	Status       POStatus    `json:"status"`

	RequisitionID *uint `json:"requisition_id,omitempty"`

	OfficeID   uint `json:"office_id"`
	LocationID uint `json:"location_id"`
	SupplierID uint `json:"supplier_id"`

	PaymentTerms    *string `json:"payment_terms,omitempty"`
	GeneralRemarks  *string `json:"general_remarks,omitempty"`
	ShippingAddress *string `json:"shipping_address,omitempty"`

	Subtotal       float64 `json:"subtotal"`
	VatAmount      float64 `json:"vat_amount"`
	DiscountAmount float64 `json:"discount_amount"`
	TotalAmount    float64 `json:"total_amount"`

	Office   *PORefDTO `json:"office,omitempty"`
	Location *PORefDTO `json:"location,omitempty"`
	Supplier *PORefDTO `json:"supplier,omitempty"`

	CreatedByID  *uint      `json:"created_by_id,omitempty"`
	ApprovedByID *uint      `json:"approved_by_id,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`

	Items []POItemResponse `json:"items"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type POItemResponse struct {
	ID                uint    `json:"id"`
	POID              uint    `json:"po_id"`
	ItemID            uint    `json:"item_id"`
	RequisitionItemID *uint   `json:"requisition_item_id,omitempty"`
	OrderQuantity     float64 `json:"order_quantity"`
	ReceivedQuantity  float64 `json:"received_quantity"`
	UOMID             *uint   `json:"uom_id,omitempty"`
	UnitPrice         float64 `json:"unit_price"`

	VatPercentage      float64 `json:"vat_percentage"`
	VatAmount          float64 `json:"vat_amount"`
	DiscountPercentage float64 `json:"discount_percentage"`
	DiscountAmount     float64 `json:"discount_amount"`
	TotalAmount        float64 `json:"total_amount"`

	Remarks *string  `json:"remarks,omitempty"`
	Item    *PORefDTO `json:"item,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

type OrderedRequisitionRow struct {
	RequisitionID     uint   `json:"requisition_id"`
	RequisitionNumber string `json:"requisition_number"`
	Status            string `json:"status"`
	POID              *uint  `json:"po_id"`
	PONumber          string `json:"po_number"`
}

type PORefDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PaginatedPOResponse struct {
	Data       []POResponse   `json:"data"`
	Pagination POPaginationMeta `json:"pagination"`
}

type POPaginationMeta struct {
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}
