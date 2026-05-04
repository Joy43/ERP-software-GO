package stock

import (
	"fmt"
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/services/purchase/locationstock"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// UpsertStock adds qty to location_stocks using weighted average cost.
// Must be called inside a transaction (tx).
func (r *Repository) UpsertStock(tx *gorm.DB, itemID, locationID uint, qty, unitCost float64) (*locationstock.LocationStock, error) {
	var ls locationstock.LocationStock

	// Lock the row for update to prevent race conditions
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("item_id = ? AND location_id = ?", itemID, locationID).
		First(&ls).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		// First stock entry for this item+location
		ls = locationstock.LocationStock{
			ItemID:      itemID,
			LocationID:  locationID,
			Quantity:    qty,
			LastCost:    unitCost,
			AverageCost: unitCost,
		}
		if createErr := tx.Create(&ls).Error; createErr != nil {
			return nil, createErr
		}
	} else {
	
		newQty := ls.Quantity + qty
		newAvgCost := ls.AverageCost
		if newQty > 0 {
			newAvgCost = (ls.Quantity*ls.AverageCost + qty*unitCost) / newQty
		}
		if err := tx.Model(&ls).Updates(map[string]interface{}{
			"quantity":     newQty,
			"last_cost":    unitCost,
			"average_cost": newAvgCost,
		}).Error; err != nil {
			return nil, err
		}
		ls.Quantity = newQty
		ls.AverageCost = newAvgCost
		ls.LastCost = unitCost
	}
	return &ls, nil
}

// GetStock returns current stock for item+location (read-only, no lock)
func (r *Repository) GetStock(itemID, locationID uint) (*locationstock.LocationStock, error) {
	var ls locationstock.LocationStock
	err := r.db.Where("item_id = ? AND location_id = ?", itemID, locationID).First(&ls).Error
	if err != nil {
		return nil, err
	}
	return &ls, nil
}

// CreateTransaction records a stock movement. Must be called inside a transaction (tx).
func (r *Repository) CreateTransaction(tx *gorm.DB, t *StockTransaction) error {
	return tx.Create(t).Error
}

func (r *Repository) GenerateTxNumber() (string, error) {
	var count int64
	today := time.Now().Format("20060102")
	prefix := "STX-" + today + "-"
	r.db.Model(&StockTransaction{}).Where("transaction_number LIKE ?", prefix+"%").Count(&count)
	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}

func (r *Repository) GetDB() *gorm.DB {
	return r.db
}
