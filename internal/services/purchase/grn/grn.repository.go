package grn

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

func (r *Repository) preload(q *gorm.DB) *gorm.DB {
	return q.
		Preload("Office").
		Preload("Location").
		Preload("Supplier").
		Preload("PaymentMethod").
		Preload("File").
		Preload("CreatedBy").
		Preload("ReceivedBy").
		Preload("Items").
		Preload("Items.Item").
		Preload("Items.UOM").
		Preload("Items.Category").
		Preload("Items.SubCategory").
		Preload("Items.MinorCategory")
}

func (r *Repository) applyFilters(q *gorm.DB, f ListGRNRequest) *gorm.DB {
	if f.Status != nil && *f.Status != "" {
		q = q.Where("status = ?", *f.Status)
	}
	if f.ReceiveType != nil && *f.ReceiveType != "" {
		q = q.Where("receive_type = ?", *f.ReceiveType)
	}
	if f.SupplierID != nil && *f.SupplierID > 0 {
		q = q.Where("supplier_id = ?", *f.SupplierID)
	}
	if f.POID != nil && *f.POID > 0 {
		q = q.Where("po_id = ?", *f.POID)
	}
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("grn_number LIKE ? OR challan_no LIKE ?", like, like)
	}
	return q
}

func (r *Repository) Count(f ListGRNRequest) (int64, error) {
	var count int64
	return count, r.applyFilters(r.db.Model(&GoodsReceiptNote{}), f).Count(&count).Error
}

func (r *Repository) FindAll(f ListGRNRequest) ([]GoodsReceiptNote, error) {
	var grns []GoodsReceiptNote
	sortBy := f.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := f.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}
	offset := (f.Page - 1) * f.PageSize
	err := r.preload(r.applyFilters(r.db, f)).
		Order(sortBy + " " + sortOrder).
		Offset(offset).Limit(f.PageSize).
		Find(&grns).Error
	return grns, err
}

func (r *Repository) FindByID(id uint) (*GoodsReceiptNote, error) {
	var g GoodsReceiptNote
	err := r.preload(r.db).Where("id = ?", id).First(&g).Error
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *Repository) Create(g *GoodsReceiptNote) (*GoodsReceiptNote, error) {
	if err := r.db.Create(g).Error; err != nil {
		return nil, err
	}
	return r.FindByID(g.ID)
}

func (r *Repository) UpdateStatus(tx *gorm.DB, id uint, status GRNStatus) error {
	return tx.Model(&GoodsReceiptNote{}).Where("id = ?", id).Update("status", status).Error
}

func (r *Repository) GenerateGRNNumber() (string, error) {
	var count int64
	today := time.Now().Format("20060102")
	prefix := "GRN-" + today + "-"
	r.db.Model(&GoodsReceiptNote{}).Where("grn_number LIKE ?", prefix+"%").Count(&count)
	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}

func (r *Repository) GetDB() *gorm.DB {
	return r.db
}
