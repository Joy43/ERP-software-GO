package category

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, d *Category) error {
	return s.repo.Create(ctx, d)
}

func (s *Service) FindAll(ctx context.Context) ([]Category, error) {
	return s.repo.FindAll(ctx)
}

// FindAllWithHierarchy retrieves all categories with full hierarchy
func (s *Service) FindAllWithHierarchy(ctx context.Context) ([]interface{}, error) {
	categories, err := s.repo.FindAllWithDepartment(ctx)
	if err != nil {
		return nil, err
	}

	var result []interface{}
	for _, category := range categories {
		// Get sub-categories for each category
		subCategoriesWithMinor, err := s.repo.GetSubCategoriesWithMinorByCategory(ctx, category.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch sub-categories: %w", err)
		}

		// Build response with hierarchy
		response := map[string]interface{}{
			"id":              category.ID,
			"name":            category.Name,
			"department_id":   category.DepartmentID,
			"department":      category.Department,
			"sub_categories": subCategoriesWithMinor,
			"created_at":      category.CreatedAt.Format("2006-01-02T15:04:05Z"),
			"updated_at":      category.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		result = append(result, response)
	}

	return result, nil
}

func (s *Service) FindByID(ctx context.Context, id uint) (*Category, error) {
	return s.repo.FindByID(ctx, id)
}

// FindByIDWithHierarchy retrieves category with full hierarchy (Department>Category>SubCategory>MinorCategory)
func (s *Service) FindByIDWithHierarchy(ctx context.Context, id uint) (interface{}, error) {

	category, err := s.repo.FindByIDWithHierarchy(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get all sub-categories for this category
	subCategoriesWithMinor, err := s.repo.GetSubCategoriesWithMinorByCategory(ctx, category.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sub-categories: %w", err)
	}

	// Build response with hierarchy
	response := map[string]interface{}{
		"id":              category.ID,
		"name":            category.Name,
		"department_id":   category.DepartmentID,
		"department":      category.Department,
		"sub_categories": subCategoriesWithMinor,
		"created_at":      category.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":      category.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	return response, nil
}

func (s *Service) Update(ctx context.Context, d *Category) error {
	return s.repo.Update(ctx, d)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
