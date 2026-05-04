package sub_category

import (
	"context"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[SubCategory]
	FindByIDWithCategory(ctx context.Context, id uint) (*SubCategory, error)
	FindAllWithCategory(ctx context.Context) ([]SubCategory, error)
}

type repositoryImpl struct {
	repository.BaseRepository[SubCategory]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		BaseRepository: repository.NewBaseRepository[SubCategory](db),
		db:             db,
	}
}

// FindByIDWithCategory retrieves a sub-category by ID with its related category
func (r *repositoryImpl) FindByIDWithCategory(ctx context.Context, id uint) (*SubCategory, error) {
	var subCategory SubCategory
	if err := r.db.WithContext(ctx).Preload("Category").First(&subCategory, id).Error; err != nil {
		return nil, err
	}
	return &subCategory, nil
}

// FindAllWithCategory retrieves all sub-categories with their related categories
func (r *repositoryImpl) FindAllWithCategory(ctx context.Context) ([]SubCategory, error) {
	var subCategories []SubCategory
	if err := r.db.WithContext(ctx).Preload("Category").Find(&subCategories).Error; err != nil {
		return nil, err
	}
	return subCategories, nil
}
