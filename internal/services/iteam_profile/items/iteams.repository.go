package items

import (
	"context"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

// ---------Repository defines the interface for create item data operations -------------
type Repository interface {
	repository.BaseRepository[Items]
	Count(ctx context.Context) (int64, error)
	FindWithOffset(ctx context.Context, offset, limit int) ([]Items, error)
}

//----------- repositoryImpl implements the Repository interface -------------
type repositoryImpl struct {
	db *gorm.DB
}

// -----------NewRepository creates and returns a new repository instance -------------
func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

//-------------- Create creates a new item -------------
func (r *repositoryImpl) Create(ctx context.Context, entity *Items) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

//--------------  FindAll retrieves all items without pagination -------------
func (r *repositoryImpl) FindAll(ctx context.Context) ([]Items, error) {
	var items []Items
	if err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("SubCategory").
		Preload("MinorCategory").
		Preload("ItemType").
		Preload("Department").
		Preload("Uom").
		Preload("Supplier").
		Preload("Tag").
		Preload("SalesSupplyType").
		Preload("SalesTaxSetup").
		Preload("SalesSetup").
		Preload("File").
		Preload("Pharmacy").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

//-------------  FindByID retrieves a single item by ID -------------
func (r *repositoryImpl) FindByID(ctx context.Context, id uint) (*Items, error) {
	var item Items
	if err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("SubCategory").
		Preload("MinorCategory").
		Preload("ItemType").
		Preload("Department").
		Preload("Uom").
		Preload("Supplier").
		Preload("Tag").
		Preload("SalesSupplyType").
		Preload("SalesTaxSetup").
		Preload("SalesSetup").
		Preload("File").
		Preload("Pharmacy").
		First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

//------------  Update updates an existing item -------------
func (r *repositoryImpl) Update(ctx context.Context, entity *Items) error {
	return r.db.WithContext(ctx).Model(entity).Updates(entity).Error
}

//------------  Delete soft deletes an item -------------
func (r *repositoryImpl) Delete(ctx context.Context, id uint) error {
	var item Items
	return r.db.WithContext(ctx).Delete(&item, id).Error
}

// ----------- Count returns the total number of items -------------
func (r *repositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&Items{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ----FindWithOffset retrieves items with pagination using offset and limit -------------
func (r *repositoryImpl) FindWithOffset(ctx context.Context, offset, limit int) ([]Items, error) {
	var items []Items
	if err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("SubCategory").
		Preload("MinorCategory").
		Preload("ItemType").
		Preload("Department").
		Preload("Uom").
		Preload("Supplier").
		Preload("Tag").
		Preload("SalesSupplyType").
		Preload("SalesTaxSetup").
		Preload("SalesSetup").
		Preload("File").
		Preload("Pharmacy").
		Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}