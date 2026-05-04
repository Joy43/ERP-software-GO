package purchasereturn


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

func (r *Repository) preloadAll(query *gorm.DB) *gorm.DB {
	return query.
		Preload("Office").
		Preload("Location").
		Preload("Supplier").
		Preload("Item")
}

func (r *Repository) applyFilters(query *gorm.DB, filter ListPurchaseReturnRequest) *gorm.DB {
	if filter.SupplierID != nil && *filter.SupplierID > 0 {
		query = query.Where("supplier_id = ?", *filter.SupplierID)
	}
	if filter.OfficeID != nil && *filter.OfficeID > 0 {
		query = query.Where("office_id = ?", *filter.OfficeID)
	}
	if filter.LocationID != nil && *filter.LocationID > 0 {
		query = query.Where("location_id = ?", *filter.LocationID)
	}
	if filter.ItemID != nil && *filter.ItemID > 0 {
		query = query.Where("item_id = ?", *filter.ItemID)
	}
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("remarks LIKE ? OR return_number LIKE ?", like, like)
	}
	return query
}

func (r *Repository) Count(filter ListPurchaseReturnRequest) (int64, error) {
	var count int64
	query := r.applyFilters(r.db.Model(&PurchaseReturn{}), filter)
	return count, query.Count(&count).Error
}

func (r *Repository) FindAll(filter ListPurchaseReturnRequest) ([]PurchaseReturn, error) {
	var records []PurchaseReturn
	offset := (filter.Page - 1) * filter.PageSize

	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := filter.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := r.applyFilters(r.db, filter)
	query = r.preloadAll(query)

	err := query.
		Order(sortBy + " " + sortOrder).
		Offset(offset).
		Limit(filter.PageSize).
		Find(&records).Error

	return records, err
}

func (r *Repository) FindByID(id uint) (*PurchaseReturn, error) {
	var record PurchaseReturn
	err := r.preloadAll(r.db).Where("id = ?", id).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *Repository) Create(record *PurchaseReturn) (*PurchaseReturn, error) {
	if err := r.db.Create(record).Error; err != nil {
		return nil, err
	}
	return r.FindByID(record.ID)
}


func (r *Repository) GenerateReturnNumber() (string, error) {
	prefix := fmt.Sprintf("PR-%s-", time.Now().Format("200601"))

	var count int64
	if err := r.db.Model(&PurchaseReturn{}).
		Where("return_number LIKE ?", prefix+"%").
		Count(&count).Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}