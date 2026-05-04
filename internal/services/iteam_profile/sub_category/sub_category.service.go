package sub_category

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, d *SubCategory) error {
	if d == nil {
		return errors.New("sub-category cannot be nil")
	}
	if d.Name == "" {
		return errors.New("sub-category name is required")
	}
	if d.CategoryID == 0 {
		return errors.New("category id is required")
	}
	return s.repo.Create(ctx, d)
}

func (s *Service) FindAll(ctx context.Context) ([]SubCategory, error) {
	return s.repo.FindAllWithCategory(ctx)
}

func (s *Service) FindByID(ctx context.Context, id uint) (*SubCategory, error) {
	if id == 0 {
		return nil, errors.New("invalid sub-category id")
	}

	subCategory, err := s.repo.FindByIDWithCategory(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("sub-category with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch sub-category: %w", err)
	}
	return subCategory, nil
}

func (s *Service) Update(ctx context.Context, d *SubCategory) error {
	if d == nil {
		return errors.New("sub-category cannot be nil")
	}
	if d.ID == 0 {
		return errors.New("sub-category id is required for update")
	}
	if d.Name == "" {
		return errors.New("sub-category name is required")
	}
	if d.CategoryID == 0 {
		return errors.New("category id is required")
	}

	// Check if record exists
	_, err := s.repo.FindByID(ctx, d.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sub-category with id %d not found", d.ID)
		}
		return fmt.Errorf("failed to verify sub-category: %w", err)
	}

	return s.repo.Update(ctx, d)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid sub-category id")
	}

	// Check if record exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sub-category with id %d not found", id)
		}
		return fmt.Errorf("failed to verify sub-category: %w", err)
	}

	return s.repo.Delete(ctx, id)
}
