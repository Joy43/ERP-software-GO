package minor_category

import (
	"context"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[MinorCategory]
	FindByIDWithSubCategory(ctx context.Context, id uint) (*MinorCategory, error)
	FindAllWithSubCategory(ctx context.Context) ([]MinorCategory, error)
}

type repositoryImpl struct {
	repository.BaseRepository[MinorCategory]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		BaseRepository: repository.NewBaseRepository[MinorCategory](db),
		db:             db,
	}
}

// FindByIDWithSubCategory retrieves a minor category by ID with its related sub-category
func (r *repositoryImpl) FindByIDWithSubCategory(ctx context.Context, id uint) (*MinorCategory, error) {
	var minorCategory MinorCategory
	if err := r.db.WithContext(ctx).Preload("SubCategory").First(&minorCategory, id).Error; err != nil {
		return nil, err
	}
	return &minorCategory, nil
}

// FindAllWithSubCategory retrieves all minor categories with their related sub-categories
func (r *repositoryImpl) FindAllWithSubCategory(ctx context.Context) ([]MinorCategory, error) {
	var minorCategories []MinorCategory
	if err := r.db.WithContext(ctx).Preload("SubCategory").Find(&minorCategories).Error; err != nil {
		return nil, err
	}
	return minorCategories, nil
}
