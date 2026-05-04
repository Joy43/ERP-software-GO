package purchaseorder

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
		Preload("CreatedBy").
		Preload("ApprovedBy").
		Preload("Items").
		Preload("Items.Item").
		Preload("Items.UOM")
}

func (r *Repository) applyFilters(q *gorm.DB, f ListPORequest) *gorm.DB {
	if f.Status != nil && *f.Status != "" {
		q = q.Where("status = ?", *f.Status)
	}
	if f.SupplierID != nil && *f.SupplierID > 0 {
		q = q.Where("supplier_id = ?", *f.SupplierID)
	}
	if f.OrderType != nil && *f.OrderType != "" {
		q = q.Where("order_type = ?", *f.OrderType)
	}
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("po_number LIKE ? OR general_remarks LIKE ?", like, like)
	}
	return q
}

func (r *Repository) Count(f ListPORequest) (int64, error) {
	var count int64
	return count, r.applyFilters(r.db.Model(&PurchaseOrder{}), f).Count(&count).Error
}

func (r *Repository) FindAll(f ListPORequest) ([]PurchaseOrder, error) {
	var pos []PurchaseOrder
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
		Find(&pos).Error
	return pos, err
}

func (r *Repository) FindByID(id uint) (*PurchaseOrder, error) {
	var po PurchaseOrder
	err := r.preload(r.db).Where("id = ?", id).First(&po).Error
	if err != nil {
		return nil, err
	}
	return &po, nil
}

func (r *Repository) Create(po *PurchaseOrder) (*PurchaseOrder, error) {
	if err := r.db.Create(po).Error; err != nil {
		return nil, err
	}
	return r.FindByID(po.ID)
}

func (r *Repository) Update(po *PurchaseOrder) (*PurchaseOrder, error) {
	err := r.db.Model(&PurchaseOrder{}).Where("id = ?", po.ID).Updates(map[string]interface{}{
		"po_number":        po.PONumber,
		"po_date":          po.PODate,
		"delivery_date":    po.DeliveryDate,
		"requisition_id":   po.RequisitionID,
		"order_type":       po.OrderType,
		"office_id":        po.OfficeID,
		"location_id":      po.LocationID,
		"supplier_id":      po.SupplierID,
		"payment_terms":    po.PaymentTerms,
		"general_remarks":  po.GeneralRemarks,
		"shipping_address": po.ShippingAddress,
		"subtotal":         po.Subtotal,
		"vat_amount":       po.VatAmount,
		"discount_amount":  po.DiscountAmount,
		"total_amount":     po.TotalAmount,
		"status":           po.Status,
		"approved_by_id":   po.ApprovedByID,
		"approved_at":      po.ApprovedAt,
	}).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(po.ID)
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&PurchaseOrder{}, id).Error
}

func (r *Repository) DeleteItems(poID uint) error {
	return r.db.Where("po_id = ?", poID).Delete(&PurchaseOrderItem{}).Error
}

// UpdateItemReceivedQty increments received_quantity on a PO item (called during GRN)
func (r *Repository) UpdateItemReceivedQty(tx *gorm.DB, poItemID uint, qty float64) error {
	return tx.Model(&PurchaseOrderItem{}).
		Where("id = ?", poItemID).
		UpdateColumn("received_quantity", gorm.Expr("received_quantity + ?", qty)).Error
}

// UpdatePOReceiveStatus recalculates PO status after GRN (PARTIALLY_RECEIVED or FULLY_RECEIVED)
func (r *Repository) UpdatePOReceiveStatus(tx *gorm.DB, poID uint) error {
	type result struct {
		Total    float64
		Received float64
	}
	var res result
	tx.Model(&PurchaseOrderItem{}).
		Select("SUM(order_quantity) as total, SUM(received_quantity) as received").
		Where("po_id = ?", poID).
		Scan(&res)

	status := POStatusPartiallyReceived
	if res.Total > 0 && res.Received >= res.Total {
		status = POStatusFullyReceived
	}
	return tx.Model(&PurchaseOrder{}).Where("id = ?", poID).
		Update("status", status).Error
}

func (r *Repository) GeneratePONumber() (string, error) {
	var count int64
	today := time.Now().Format("20060102")
	prefix := "PO-" + today + "-"
	r.db.Model(&PurchaseOrder{}).Where("po_number LIKE ?", prefix+"%").Count(&count)
	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}

// GetOrderedRequisitions returns requisitions with status ORDERED and their linked PO id
func (r *Repository) GetOrderedRequisitions() ([]OrderedRequisitionRow, error) {
	var rows []OrderedRequisitionRow
	err := r.db.Table("requisitions r").
		Select("r.id as requisition_id, r.requisition_number, r.status, po.id as po_id, po.po_number").
		Joins("LEFT JOIN purchase_orders po ON po.requisition_id = r.id AND po.deleted_at IS NULL").
		Where("r.status = ? AND r.deleted_at IS NULL", "ORDERED").
		Order("r.created_at DESC").
		Scan(&rows).Error
	return rows, err
}

// GetDB exposes db for transaction use in service layer
func (r *Repository) GetDB() *gorm.DB {
	return r.db
}
