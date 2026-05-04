package requisitions

import (
	"fmt"
	"time"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// preloadAll adds all standard preloads
func (r *Repository) preloadAll(query *gorm.DB) *gorm.DB {
	return query.
		Preload("Employee").
		Preload("Buyer").
		Preload("Department").
		Preload("Project").
		Preload("Office").
		Preload("Location").
		Preload("InventoryType").
		Preload("Supplier").
		Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("Items").
		Preload("Items.Item").
		Preload("Items.ItemType").
		Preload("Items.Category").
		Preload("Items.SubCategory").
		Preload("Items.MinorCategory").
		Preload("StatusHistory").
		Preload("StatusHistory.User")
}

func (r *Repository) applyFilters(query *gorm.DB, filter ListRequisitionRequest) *gorm.DB {
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.RequisitionType != nil && *filter.RequisitionType != "" {
		query = query.Where("requisition_type = ?", *filter.RequisitionType)
	}
	if filter.DepartmentID != nil && *filter.DepartmentID > 0 {
		query = query.Where("department_id = ?", *filter.DepartmentID)
	}
	if filter.ProjectID != nil && *filter.ProjectID > 0 {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.EmployeeID != nil && *filter.EmployeeID > 0 {
		query = query.Where("employee_id = ?", *filter.EmployeeID)
	}
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("remarks LIKE ? OR description LIKE ? OR requisition_number LIKE ?", like, like, like)
	}
	return query
}

func (r *Repository) Count(filter ListRequisitionRequest) (int64, error) {
	var count int64
	query := r.applyFilters(r.db.Model(&Requisition{}), filter)
	return count, query.Count(&count).Error
}

func (r *Repository) FindAll(filter ListRequisitionRequest) ([]Requisition, error) {
	var requisitions []Requisition
	offset := (filter.Page - 1) * filter.PageSize

	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := filter.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}

	query := r.applyFilters(r.db, filter)
	query = r.preloadAll(query)

	err := query.
		Order(sortBy + " " + sortOrder).
		Offset(offset).
		Limit(filter.PageSize).
		Find(&requisitions).Error

	return requisitions, err
}

func (r *Repository) FindByID(id uint) (*Requisition, error) {
	var req Requisition
	err := r.preloadAll(r.db).Where("id = ?", id).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *Repository) Create(req *Requisition) (*Requisition, error) {
	if err := r.db.Create(req).Error; err != nil {
		return nil, err
	}
	return r.FindByID(req.ID)
}

func (r *Repository) Update(req *Requisition) (*Requisition, error) {
	err := r.db.Model(&Requisition{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"requisition_number":  req.RequisitionNumber,
		"requisition_type":    req.RequisitionType,
		"status":              req.Status,
		"rejection_reason":    req.RejectionReason,
		"expected_date":       req.ExpectedDate,
		"remarks":             req.Remarks,
		"description":         req.Description,
		"employee_id":         req.EmployeeID,
		"department_id":       req.DepartmentID,
		"project_id":          req.ProjectID,
		"buyer_id":            req.BuyerID,
		"office_id":           req.OfficeID,
		"location_id":         req.LocationID,
		"inventory_type_id":   req.InventoryTypeID,
		"supplier_id":         req.SupplierID,
		"updated_by_id":       req.UpdatedByID,
	}).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(req.ID)
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Requisition{}, id).Error
}

// AddStatusHistory records a status change
func (r *Repository) AddStatusHistory(history *RequisitionStatusHistory) error {
	return r.db.Create(history).Error
}

// GetStatusHistory returns all history records for a requisition ordered by created_at asc
func (r *Repository) GetStatusHistory(requisitionID uint) ([]RequisitionStatusHistory, error) {
	var history []RequisitionStatusHistory
	err := r.db.Where("requisition_id = ?", requisitionID).
		Preload("User").
		Order("created_at asc").
		Find(&history).Error
	return history, err
}

// GetAllHistory returns all history records across all requisitions with optional filters
func (r *Repository) GetAllHistory(filter ListHistoryRequest) ([]RequisitionStatusHistory, int64, error) {
	var history []RequisitionStatusHistory
	var count int64

	query := r.db.Model(&RequisitionStatusHistory{})

	if filter.RequisitionID != nil {
		query = query.Where("requisition_id = ?", *filter.RequisitionID)
	}
	if filter.ActionType != nil && *filter.ActionType != "" {
		query = query.Where("action_type = ?", *filter.ActionType)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	query.Count(&count)

	offset := (filter.Page - 1) * filter.PageSize
	err := query.Preload("User").
		Order("created_at desc").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&history).Error

	return history, count, err
}

// Transaction methods
func (r *Repository) Begin() *Repository {
	tx := r.db.Begin()
	return &Repository{db: tx}
}

func (r *Repository) Commit() error {
	return r.db.Commit().Error
}

func (r *Repository) Rollback() {
	r.db.Rollback()
}

// CreateItems creates multiple requisition items in a transaction
func (r *Repository) CreateItems(items []RequisitionItem) error {
	return r.db.Create(&items).Error
}

// GetExistingItemIDs validates that the provided item IDs exist in the database
func (r *Repository) GetExistingItemIDs(itemIDs []uint) ([]uint, error) {
	var existingIDs []uint
	err := r.db.Table("items").Where("id IN ?", itemIDs).Pluck("id", &existingIDs).Error
	return existingIDs, err
}

// DeleteItems removes existing items before re-saving
func (r *Repository) DeleteItems(requisitionID uint) error {
	return r.db.Where("requisition_id = ?", requisitionID).Delete(&RequisitionItem{}).Error
}

// GetRequisitionSummary returns total requisitions, approved count, total value, and breakdown by category
func (r *Repository) GetRequisitionSummary() (*RequisitionSummaryResponse, error) {
	var total, approved int64
	var totalValue float64

	r.db.Model(&Requisition{}).Count(&total)
	r.db.Model(&Requisition{}).Where("status = ?", StatusApproved).Count(&approved)
	r.db.Table("requisition_items").
		Select("COALESCE(SUM(request_quantity * COALESCE(average_cost, last_cost, 0)), 0)").
		Scan(&totalValue)

	type categoryRow struct {
		Category string
		Total    int64
		Approved int64
		Value    float64
	}
	var rows []categoryRow
	r.db.Table("requisition_items ri").
		Select(`
			COALESCE(c.name, 'Uncategorized') AS category,
			COUNT(DISTINCT ri.requisition_id) AS total,
			COUNT(DISTINCT CASE WHEN req.status = ? THEN ri.requisition_id END) AS approved,
			COALESCE(SUM(ri.request_quantity * COALESCE(ri.average_cost, ri.last_cost, 0)), 0) AS value
		`, StatusApproved).
		Joins("JOIN requisitions req ON req.id = ri.requisition_id").
		Joins("LEFT JOIN categories c ON c.id = ri.category_id").
		Group("c.name").
		Scan(&rows)

	byCategory := make([]RequisitionCategoryRow, 0, len(rows))
	for _, row := range rows {
		byCategory = append(byCategory, RequisitionCategoryRow{
			Category: row.Category,
			Total:    row.Total,
			Approved: row.Approved,
			Value:    row.Value,
		})
	}

	return &RequisitionSummaryResponse{
		TotalRequisitions: total,
		Approved:          approved,
		TotalValue:        totalValue,
		ByCategory:        byCategory,
	}, nil
}

// GenerateRequisitionNumber generates a unique requisition number like REQ-20240601-0001
func (r *Repository) GenerateRequisitionNumber() (string, error) {
	var count int64
	today := time.Now().Format("20060102")
	prefix := "REQ-" + today + "-"
	r.db.Model(&Requisition{}).Where("requisition_number LIKE ?", prefix+"%").Count(&count)
	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}