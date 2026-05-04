package requisitions

import "time"

// ========================================
//-----------  REQUEST DTOs -----------
// ========================================


type CreateRequisitionRequest struct {
	RequisitionType string `json:"requisition_type" binding:"required,oneof=EMPLOYEE DEPARTMENT PROJECT"`
	ExpectedDate    string `json:"expected_date" binding:"required"` 

	// --------- Party: RequisitionType ---------
	EmployeeID   *uint `json:"employee_id"`
	DepartmentID *uint `json:"department_id"`
	ProjectID    *uint `json:"project_id"`

	// --------- Optional fields ---------
	BuyerID         *uint   `json:"buyer_id"`
	OfficeID        *uint   `json:"office_id"`
	LocationID      *uint   `json:"location_id"`
	InventoryTypeID *uint   `json:"inventory_type_id"`
	SupplierID      *uint   `json:"supplier_id"`
	Remarks         *string `json:"remarks"`
	Description     *string `json:"description"`

	// --------- Items ---------
	Items []CreateRequisitionItemRequest `json:"items" binding:"required,min=1,dive"`
}

type CreateRequisitionItemRequest struct {
	ItemID          uint     `json:"item_id" binding:"required,gt=0"`
	RequestQuantity float64  `json:"request_quantity" binding:"required,gt=0"`
	CurrentStock    *float64 `json:"current_stock"`
	LastCost        *float64 `json:"last_cost"`
	AverageCost     *float64 `json:"average_cost"`
	Description     *string  `json:"description"`

	//-----------  Classification ---------
	ItemTypeID      *uint `json:"item_type_id"`
	CategoryID      *uint `json:"category_id"`
	SubCategoryID   *uint `json:"sub_category_id"`
	MinorCategoryID *uint `json:"minor_category_id"`
}

// UpdateRequisitionItemRequest - for partial updates, all fields are optional
type UpdateRequisitionItemRequest struct {
	ItemID          *uint    `json:"item_id"`
	RequestQuantity *float64 `json:"request_quantity"`
	CurrentStock    *float64 `json:"current_stock"`
	LastCost        *float64 `json:"last_cost"`
	AverageCost     *float64 `json:"average_cost"`
	Description     *string  `json:"description"`

	//-----------  Classification ---------
	ItemTypeID      *uint `json:"item_type_id"`
	CategoryID      *uint `json:"category_id"`
	SubCategoryID   *uint `json:"sub_category_id"`
	MinorCategoryID *uint `json:"minor_category_id"`
}

// UpdateRequisitionRequest -only update DRAFT status requisition can be updated. Status change is handled separately
type UpdateRequisitionRequest struct {
	ExpectedDate    *string `json:"expected_date"`
	BuyerID         *uint   `json:"buyer_id"`
	OfficeID        *uint   `json:"office_id"`
	LocationID      *uint   `json:"location_id"`
	InventoryTypeID *uint   `json:"inventory_type_id"`
	SupplierID      *uint   `json:"supplier_id"`
	Remarks         *string `json:"remarks"`
	Description     *string `json:"description"`

	//------ Items update CAN be update but no optional---
	Items []UpdateRequisitionItemRequest `json:"items"`
}

//------- UpdateRequisitionStatusRequest - Status ----
type UpdateRequisitionStatusRequest struct {
	Status          string  `json:"status" binding:"required"`
	Remarks         *string `json:"remarks"`
	RejectionReason *string `json:"rejection_reason"`
}

// ListHistoryRequest - filter/pagination for all history
type ListHistoryRequest struct {
	Page          int     `form:"page,default=1"`
	PageSize      int     `form:"page_size,default=20"`
	RequisitionID *uint   `form:"requisition_id"`
	ActionType    *string `form:"action_type"`
	UserID        *uint   `form:"user_id"`
}

// PaginatedHistoryResponse
type PaginatedHistoryResponse struct {
	Data       []RequisitionStatusHistoryResponse `json:"data"`
	Pagination PaginationMeta                     `json:"pagination"`
}

// ListRequisitionRequest - Filter/pagination
type ListRequisitionRequest struct {
	Page            int     `form:"page,default=1"`
	PageSize        int     `form:"page_size,default=10"`
	Status          *string `form:"status"`
	RequisitionType *string `form:"requisition_type"`
	DepartmentID    *uint   `form:"department_id"`
	ProjectID       *uint   `form:"project_id"`
	EmployeeID      *uint   `form:"employee_id"`
	Search          string  `form:"search"`
	SortBy          string  `form:"sort_by,default=created_at"`
	SortOrder       string  `form:"sort_order,default=desc"`
}

// ========================================
// RESPONSE DTOs
// ========================================

type RequisitionResponse struct {
	ID                uint              `json:"id"`
	RequisitionNumber string            `json:"requisition_number"`
	RequisitionType   RequisitionType   `json:"requisition_type"`
	Status            RequisitionStatus `json:"status"`
	RejectionReason   *string           `json:"rejection_reason,omitempty"`
	ExpectedDate      time.Time         `json:"expected_date"`
	CreatedDate       time.Time         `json:"created_date"`
	Remarks           *string           `json:"remarks,omitempty"`
	Description       *string           `json:"description,omitempty"`

	//-------------  IDs --------------------
	EmployeeID      *uint `json:"employee_id,omitempty"`
	DepartmentID    *uint `json:"department_id,omitempty"`
	ProjectID       *uint `json:"project_id,omitempty"`
	BuyerID         *uint `json:"buyer_id,omitempty"`
	OfficeID        *uint `json:"office_id,omitempty"`
	LocationID      *uint `json:"location_id,omitempty"`
	InventoryTypeID *uint `json:"inventory_type_id,omitempty"`
	SupplierID      *uint `json:"supplier_id,omitempty"`
	CreatedByID     *uint `json:"created_by_id,omitempty"`
	UpdatedByID     *uint `json:"updated_by_id,omitempty"`

	//------------------ Relations (populated when available) -----------------
	Employee      *UserDTO          `json:"employee,omitempty"`
	Department    *DepartmentDTO    `json:"department,omitempty"`
	Project       *ProjectDTO       `json:"project,omitempty"`
	Buyer         *UserDTO          `json:"buyer,omitempty"`
	Office        *OfficeDTO        `json:"office,omitempty"`
	Location      *LocationDTO      `json:"location,omitempty"`
	InventoryType *InventoryTypeDTO `json:"inventory_type,omitempty"`
	Supplier      *SupplierDTO      `json:"supplier,omitempty"`
	CreatedBy     *UserDTO          `json:"created_by,omitempty"`
	UpdatedBy     *UserDTO          `json:"updated_by,omitempty"`

	Items         []RequisitionItemResponse          `json:"items"`
	StatusHistory []RequisitionStatusHistoryResponse `json:"status_history,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RequisitionItemResponse struct {
	ID              uint     `json:"id"`
	RequisitionID   uint     `json:"requisition_id"`
	ItemID          uint     `json:"item_id"`
	RequestQuantity float64  `json:"request_quantity"`
	ApprovedQuantity *float64 `json:"approved_quantity,omitempty"`
	CurrentStock    *float64 `json:"current_stock,omitempty"`
	LastCost        *float64 `json:"last_cost,omitempty"`
	AverageCost     *float64 `json:"average_cost,omitempty"`
	Description     *string  `json:"description,omitempty"`
	Item          *ItemDTO          `json:"item,omitempty"`
	ItemType      *ItemTypeDTO      `json:"item_type,omitempty"`
	Category      *CategoryDTO      `json:"category,omitempty"`
	SubCategory   *SubCategoryDTO   `json:"sub_category,omitempty"`
	MinorCategory *MinorCategoryDTO `json:"minor_category,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type RequisitionStatusHistoryResponse struct {
	ID            uint      `json:"id"`
	RequisitionID uint      `json:"requisition_id"`
	ActionType    string    `json:"action_type"`
	FromStatus    string    `json:"from_status"`
	ToStatus      string    `json:"to_status"`
	Remarks       *string   `json:"remarks,omitempty"`
	User          *UserDTO  `json:"user,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// ========================================
// PAGINATION
// ========================================

type PaginatedRequisitionResponse struct {
	Data       []RequisitionResponse `json:"data"`
	Pagination PaginationMeta        `json:"pagination"`
}

type PaginationMeta struct {
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// ========================================
//--------- SHARED DTOs -----------
// ========================================

type UserDTO struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile,omitempty"`
	EmployeeID string `json:"employee_id,omitempty"`
}

type DepartmentDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProjectDTO struct {
	ID          uint   `json:"id"`
	ProjectName string `json:"project_name"`
}

type OfficeDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type LocationDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`

}

type CategoryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SubCategoryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type MinorCategoryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SupplierDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type InventoryTypeDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ItemTypeDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ItemDTO struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	SKU     string `json:"sku,omitempty"`
	Barcode string `json:"barcode,omitempty"`
}

type RequisitionSummaryResponse struct {
	TotalRequisitions int64                    `json:"total_requisitions"`
	Approved          int64                    `json:"approved"`
	TotalValue        float64                  `json:"total_value"`
	ByCategory        []RequisitionCategoryRow `json:"by_category"`
}

type RequisitionCategoryRow struct {
	Category string  `json:"category"`
	Total    int64   `json:"total"`
	Approved int64   `json:"approved"`
	Value    float64 `json:"value"`
}