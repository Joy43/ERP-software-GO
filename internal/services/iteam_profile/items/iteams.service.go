package items

import (
	"context"
	"errors"
	"fmt"
	"time"
	"gorm.io/gorm"
)

// Service handles business logic for create item operations
type Service struct {
	repo Repository
	db   *gorm.DB
}

// NewService creates and returns a new service instance
func NewService(repo Repository, db *gorm.DB) *Service {
	return &Service{repo: repo, db: db}
}


// -------------------------------
//......... creates a new item with optional pharmacy data ...........
// --------------------------------
func (s *Service) Create(ctx context.Context, item *Items, pharmacy *ItemPharmacy) error {
	if item == nil {
		return errors.New("item cannot be nil")
	}

	if err := validateItems(item); err != nil {
		return err
	}

	// If item is a pharmacy item, validate pharmacy data
	if item.IsPharmacy && pharmacy == nil {
		return errors.New("pharmacy data is required when is_pharmacy is true")
	}

	if item.IsPharmacy && pharmacy != nil {
		if err := validateItemPharmacy(pharmacy); err != nil {
			return err
		}
	}

	// Use transaction for creating item and pharmacy data together
	if item.IsPharmacy {
		return s.createItemWithPharmacyTransaction(ctx, item, pharmacy)
	}

	// Create regular item without pharmacy
	if err := s.repo.Create(ctx, item); err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return nil
}

// createItemWithPharmacyTransaction creates an item and its pharmacy record within a transaction
func (s *Service) createItemWithPharmacyTransaction(ctx context.Context, item *Items, pharmacy *ItemPharmacy) error {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Create the item first
	if err := tx.Create(item).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create item in transaction: %w", err)
	}

	// Set the pharmacy's item ID to the newly created item's ID
	pharmacy.ItemID = item.ID

	// Create the pharmacy record
	if err := tx.Create(pharmacy).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create pharmacy record in transaction: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

//-----------generate a unique barcode for the item (this is a placeholder, implement your own logic)--------
func generateBarcode() string {
	return fmt.Sprintf("BC-%d", time.Now().UnixNano())
}
// ----------FindAll retrieves all items without pagination----------
func (s *Service) FindAll(ctx context.Context) ([]Items, error) {
	items, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}
	return items, nil
}

//----------- FindAllPaginated retrieves items with pagination support--------
func (s *Service) FindAllPaginated(ctx context.Context, page, limit int) ([]Items, int64, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Get total count
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count items: %w", err)
	}

	// Get paginated items
	items, err := s.repo.FindWithOffset(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch items: %w", err)
	}

	return items, total, nil
}

// -------------- FindByID retrieves a specific item by ID--------------------
func (s *Service) FindByID(ctx context.Context, id uint) (*Items, error) {
	if id == 0 {
		return nil, errors.New("invalid item id")
	}

	item, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("item not found")
		}
		return nil, fmt.Errorf("failed to fetch item: %w", err)
	}

	return item, nil
}

//-------- Update validates and updates an existing item----------
func (s *Service) Update(ctx context.Context, item *Items) error {
	if item == nil {
		return errors.New("item cannot be nil")
	}

	if item.ID == 0 {
		return errors.New("item id is required for update")
	}

	// --------- Validate required fields ------------
	if err := validateItems(item); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, item); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

//---------- Delete removes an item by ID (soft delete) ------------
func (s *Service) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid item id")
	}

	// ------------Check if item exists before attempting to delete -------------
	if _, err := s.FindByID(ctx, id); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("item not found")
		}
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}



// ========== validateItems validates required fields for item creation/update==================
func validateItems(item *Items) error {
	if item.Name == "" {
		return errors.New("item name is required")
	}
	// Only require barcode when AutoBarcode is FALSE
	if !item.AutoBarcode && item.Barcode == "" {
		return errors.New("barcode is required when auto_barcode is false")
	}
	if item.AutoBarcode {
		item.Barcode = generateBarcode()
	}
	if len(item.Name) < 2 || len(item.Name) > 150 {
		return errors.New("item name must be between 2 and 150 characters")
	}

	// Barcode is only required when auto_barcode is false
	if !item.AutoBarcode && item.Barcode == "" {
		return errors.New("barcode is required when auto_barcode is false")
	}

	if item.SKU == "" {
		return errors.New("sku is required")
	}

	if item.CategoryID == 0 {
		return errors.New("category_id is required")
	}

	if item.SubCategoryID == 0 {
		return errors.New("sub_category_id is required")
	}

	if item.MinorCategoryID == 0 {
		return errors.New("minor_category_id is required")
	}

	if item.ItemTypeID == 0 {
		return errors.New("item_type_id is required")
	}

	if item.DepartmentID == 0 {
		return errors.New("department_id is required")
	}

	if item.UomId == 0 {
		return errors.New("uom_id is required")
	}

	if item.SupplierID == 0 {
		return errors.New("supplier_id is required")
	}

	if item.CostPrice < 0 {
		return errors.New("cost_price cannot be negative")
	}

	if item.StandardSalesPrice < 0 {
		return errors.New("standard_sales_price cannot be negative")
	}

	if item.StandardSalesPrice < item.CostPrice {
		return errors.New("standard_sales_price cannot be less than cost_price")
	}

	if item.SalesSupplyTypeID == 0 {
		return errors.New("sales_supply_type_id is required")
	}

	if item.MaxDiscount < 0 || item.MaxDiscount > 100 {
		return errors.New("max_discount must be between 0 and 100")
	}

	return nil
}

// ========== validateItemPharmacy validates pharmacy-specific fields ==================
func validateItemPharmacy(pharmacy *ItemPharmacy) error {
	if pharmacy == nil {
		return errors.New("pharmacy data cannot be nil")
	}

	if pharmacy.GenericName == "" {
		return errors.New("generic_name is required for pharmacy items")
	}

	if len(pharmacy.GenericName) > 200 {
		return errors.New("generic_name must not exceed 200 characters")
	}

	if len(pharmacy.BrandName) > 200 {
		return errors.New("brand_name must not exceed 200 characters")
	}

	if len(pharmacy.Strength) > 100 {
		return errors.New("strength must not exceed 100 characters")
	}

	if len(pharmacy.DosageForm) > 100 {
		return errors.New("dosage_form must not exceed 100 characters")
	}

	if len(pharmacy.DrugRegistrationNo) > 100 {
		return errors.New("drug_registration_no must not exceed 100 characters")
	}

	if pharmacy.ShelfLifeDays != nil && *pharmacy.ShelfLifeDays < 0 {
		return errors.New("shelf_life_days cannot be negative")
	}

	if pharmacy.ReorderAlertDays != nil && *pharmacy.ReorderAlertDays < 0 {
		return errors.New("reorder_alert_days cannot be negative")
	}

	return nil
}