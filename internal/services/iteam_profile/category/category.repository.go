package category

import (
	"context"
	"time"

	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/shared/repository"
	"gorm.io/gorm"
)

type Repository interface {
	repository.BaseRepository[Category]
	FindByIDWithHierarchy(ctx context.Context, id uint) (*Category, error)
	FindAllWithDepartment(ctx context.Context) ([]Category, error)
	GetSubCategoriesWithMinorByCategory(ctx context.Context, categoryID uint) ([]map[string]interface{}, error)
}

type repositoryImpl struct {
	repository.BaseRepository[Category]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{
		BaseRepository: repository.NewBaseRepository[Category](db),
		db:             db,
	}
}

// FindByIDWithHierarchy retrieves category with Department and all related SubCategories and MinorCategories
func (r *repositoryImpl) FindByIDWithHierarchy(ctx context.Context, id uint) (*Category, error) {
	var category Category
	if err := r.db.WithContext(ctx).
		Preload("Department").
		First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// FindAllWithDepartment retrieves all categories with Department preloaded
func (r *repositoryImpl) FindAllWithDepartment(ctx context.Context) ([]Category, error) {
	var categories []Category
	if err := r.db.WithContext(ctx).
		Preload("Department").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetSubCategoriesWithMinorByCategory retrieves all sub-categories with their minor categories for a given category
func (r *repositoryImpl) GetSubCategoriesWithMinorByCategory(ctx context.Context, categoryID uint) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	// First get all sub-categories for this category
	type SubCatRow struct {
		ID         uint
		Name       string
		CategoryID uint
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}

	var subCats []SubCatRow
	if err := r.db.WithContext(ctx).
		Table("sub_categories").
		Where("category_id = ?", categoryID).
		Scan(&subCats).Error; err != nil {
		return nil, err
	}

	// For each sub-category, get its minor categories
	for _, sc := range subCats {
		type MinorCatRow struct {
			ID            uint
			Name          string
			SubCategoryID uint
			CreatedAt     time.Time
			UpdatedAt     time.Time
		}

		var minorCats []MinorCatRow
		if err := r.db.WithContext(ctx).
			Table("minor_categories").
			Where("sub_category_id = ?", sc.ID).
			Scan(&minorCats).Error; err != nil {
			return nil, err
		}

		// Convert minor categories to map format
		minorCatsData := make([]map[string]interface{}, len(minorCats))
		for i, mc := range minorCats {
			minorCatsData[i] = map[string]interface{}{
				"id":              mc.ID,
				"name":            mc.Name,
				"sub_category_id": mc.SubCategoryID,
				"created_at":      mc.CreatedAt.Format("2006-01-02T15:04:05Z"),
				"updated_at":      mc.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			}
		}

		// Build sub-category response
		scData := map[string]interface{}{
			"id":               sc.ID,
			"name":             sc.Name,
			"category_id":      sc.CategoryID,
			"created_at":       sc.CreatedAt.Format("2006-01-02T15:04:05Z"),
			"updated_at":       sc.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			"minor_categories": minorCatsData,
		}

		result = append(result, scData)
	}

	return result, nil
}
