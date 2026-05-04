package requisitions

import (
	"time"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/department"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/office"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/auth/user"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/item_type"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/items"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/minor_category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/iteam_profile/sub_category"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/partner/supplier"
	inventorytypes "github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/inventory_types"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/location"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/projects"
)


// ========================================
//  REQUISITION MODEL
// ========================================

type Requisition struct {
	ID                  uint              `gorm:"primaryKey" json:"id"`
	RequisitionNumber   string            `gorm:"uniqueIndex;not null;size:50" json:"requisition_number"`
	RequisitionType     RequisitionType   `gorm:"type:enum('EMPLOYEE','DEPARTMENT','PROJECT');default:'EMPLOYEE';not null;index" json:"requisition_type"`
	Status              RequisitionStatus `gorm:"type:enum('PENDING','DEPARTMENT_APPROVED','FINANCE_APPROVED','APPROVED','REJECTED','CANCELLED','ORDERED');default:'PENDING';not null;index" json:"status"`
	RejectionReason     *string           `gorm:"type:text" json:"rejection_reason,omitempty"`

	// --------- Dates ---------
	ExpectedDate time.Time  `gorm:"type:date;not null" json:"expected_date"`
	CreatedDate  time.Time  `gorm:"type:date;not null;default:CURRENT_DATE" json:"created_date"`

	// --------- Remarks ---------
	Remarks     *string `gorm:"type:text" json:"remarks,omitempty"`
	Description *string `gorm:"type:text" json:"description,omitempty"`

	// --------- Party: only ONE of these will be set based on RequisitionType ---------
	EmployeeID   *uint                  `gorm:"index" json:"employee_id,omitempty"`
	Employee     *user.User             `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"employee,omitempty"`
	DepartmentID *uint                  `gorm:"index" json:"department_id,omitempty"`
	Department   *department.Department `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"department,omitempty"`
	ProjectID    *uint                  `gorm:"index" json:"project_id,omitempty"`
	Project      *projects.Project      `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"project,omitempty"`

	// --------- Buyer (assigned by admin/procurement) ---------
	BuyerID *uint      `gorm:"index" json:"buyer_id,omitempty"`
	Buyer   *user.User `gorm:"foreignKey:BuyerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"buyer,omitempty"`

	// --------- Location ---------
	OfficeID   *uint              `gorm:"index" json:"office_id,omitempty"`
	Office     *office.Office     `gorm:"foreignKey:OfficeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"office,omitempty"`
	LocationID *uint              `gorm:"index" json:"location_id,omitempty"`
	Location   *location.Location `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"location,omitempty"`

	// --------- Inventory Type ---------
	InventoryTypeID *uint                          `json:"inventory_type_id,omitempty"`
	InventoryType   *inventorytypes.InventoryType  `gorm:"foreignKey:InventoryTypeID" json:"inventory_type,omitempty"`

	// --------- Supplier  ---------
	SupplierID *uint              `json:"supplier_id,omitempty"`
	Supplier   *supplier.Supplier `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`

	// --------- Audit ---------
	CreatedByID *uint      `gorm:"index" json:"created_by_id,omitempty"`
	CreatedBy   *user.User `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"created_by,omitempty"`
	UpdatedByID *uint      `json:"updated_by_id,omitempty"`
	UpdatedBy   *user.User `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"updated_by,omitempty"`

	// --------- Items  ---------
	Items []RequisitionItem `gorm:"foreignKey:RequisitionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"items,omitempty"`

	// --------- Status History ---------
	StatusHistory []RequisitionStatusHistory `gorm:"foreignKey:RequisitionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"status_history,omitempty"`

	// --------- Timestamps ---------
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Requisition) TableName() string {
	return "requisitions"
}

// ========================================
// REQUISITION ITEM MODEL (separate table)
// ========================================

type RequisitionItem struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	RequisitionID uint    `gorm:"not null;index" json:"requisition_id"`

	// --------- Item ---------
	ItemID uint        `gorm:"not null;index" json:"item_id"`
	Item   items.Items `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"item,omitempty"`

	// --------- Quantities ---------
	RequestQuantity  float64  `gorm:"type:decimal(15,2);not null" json:"request_quantity"`
	ApprovedQuantity *float64 `gorm:"type:decimal(15,2)" json:"approved_quantity,omitempty"`
	CurrentStock     *float64 `gorm:"type:decimal(15,2);comment:Stock at requisition time" json:"current_stock,omitempty"`

	// --------- Costing ---------
	LastCost    *float64 `gorm:"type:decimal(15,2)" json:"last_cost,omitempty"`
	AverageCost *float64 `gorm:"type:decimal(15,2)" json:"average_cost,omitempty"`

	// --------- Description ---------
	Description *string `gorm:"type:text" json:"description,omitempty"`

	// --------- Item Classification (snapshot at requisition time) ---------
	ItemTypeID      *uint                        `json:"item_type_id,omitempty"`
	ItemType        *item_type.ItemType           `gorm:"foreignKey:ItemTypeID" json:"item_type,omitempty"`
	CategoryID      *uint                        `json:"category_id,omitempty"`
	Category        *category.Category           `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	SubCategoryID   *uint                        `json:"sub_category_id,omitempty"`
	SubCategory     *sub_category.SubCategory    `gorm:"foreignKey:SubCategoryID" json:"sub_category,omitempty"`
	MinorCategoryID *uint                        `json:"minor_category_id,omitempty"`
	MinorCategory   *minor_category.MinorCategory `gorm:"foreignKey:MinorCategoryID" json:"minor_category,omitempty"`

	// --------- Timestamps ---------
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (RequisitionItem) TableName() string {
	return "requisition_items"
}

// ========================================
// ------- STATUS HISTORY MODEL -------
// ========================================

type RequisitionStatusHistory struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	RequisitionID uint      `gorm:"not null;index" json:"requisition_id"`
	UserID        *uint     `gorm:"index" json:"user_id,omitempty"`
	User          *user.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user,omitempty"`
	ActionType    string    `gorm:"type:enum('CREATED','UPDATED','STATUS_CHANGED');not null" json:"action_type"`
	FromStatus    string    `gorm:"size:50" json:"from_status"`
	ToStatus      string    `gorm:"not null;size:50" json:"to_status"`
	Remarks       *string   `gorm:"type:text" json:"remarks,omitempty"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

const (
	ActionCreated       = "CREATED"
	ActionUpdated       = "UPDATED"
	ActionStatusChanged = "STATUS_CHANGED"
)

func (RequisitionStatusHistory) TableName() string {
	return "requisition_status_history"
}



// ========================================
//--------------  ENUMS -------------------
// ========================================

type RequisitionType string

const (
	EmployeeRequisition   RequisitionType = "EMPLOYEE"
	DepartmentRequisition RequisitionType = "DEPARTMENT"
	ProjectRequisition    RequisitionType = "PROJECT"
)

type RequisitionStatus string

const (
	StatusPending             RequisitionStatus = "PENDING"
	StatusDepartmentApproved  RequisitionStatus = "DEPARTMENT_APPROVED"
	StatusFinanceApproved     RequisitionStatus = "FINANCE_APPROVED"
	StatusApproved            RequisitionStatus = "APPROVED"
	StatusRejected            RequisitionStatus = "REJECTED"
	StatusCancelled           RequisitionStatus = "CANCELLED"
	StatusOrdered             RequisitionStatus = "ORDERED"
)

// ValidStatusTransitions defines allowed status changes
// Key = current status, Value = allowed next statuses
var ValidStatusTransitions = map[RequisitionStatus][]RequisitionStatus{
	StatusPending:            {StatusDepartmentApproved, StatusRejected, StatusCancelled},
	StatusDepartmentApproved: {StatusFinanceApproved, StatusRejected, StatusCancelled},
	StatusFinanceApproved:    {StatusApproved, StatusRejected, StatusCancelled},
	StatusApproved:           {StatusOrdered, StatusCancelled},
	StatusRejected:           {StatusPending},
	StatusCancelled:          {},
	StatusOrdered:            {},
}

// ---------- IsValidTransition checks if status change is allowed---------
func (current RequisitionStatus) IsValidTransition(next RequisitionStatus) bool {
	allowed, ok := ValidStatusTransitions[current]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}
