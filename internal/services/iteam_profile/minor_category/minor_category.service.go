package minor_category

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

func (s *Service) Create(ctx context.Context, d *MinorCategory) error {
	// Validate input
	if d == nil {
		return errors.New("minor category cannot be nil")
	}
	if d.Name == "" {
		return errors.New("minor category name is required")
	}

	// Create the minor category
	if err := s.repo.Create(ctx, d); err != nil {
		return fmt.Errorf("failed to create minor category: %w", err)
	}
	return nil
}

func (s *Service) FindAll(ctx context.Context) ([]MinorCategory, error) {
	minorCategories, err := s.repo.FindAllWithSubCategory(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch minor categories: %w", err)
	}
	return minorCategories, nil
}

func (s *Service) FindByID(ctx context.Context, id uint) (*MinorCategory, error) {
	if id == 0 {
		return nil, errors.New("invalid minor category id")
	}

	minorCategory, err := s.repo.FindByIDWithSubCategory(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("minor category with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch minor category: %w", err)
	}
	return minorCategory, nil
}

func (s *Service) Update(ctx context.Context, d *MinorCategory) error {
	// Validate input
	if d == nil {
		return errors.New("minor category cannot be nil")
	}
	if d.ID == 0 {
		return errors.New("minor category id is required for update")
	}
	if d.Name == "" {
		return errors.New("minor category name is required")
	}

	// Check if record exists
	_, err := s.repo.FindByID(ctx, d.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("minor category with id %d not found", d.ID)
		}
		return fmt.Errorf("failed to verify minor category: %w", err)
	}

	// Update the minor category
	if err := s.repo.Update(ctx, d); err != nil {
		return fmt.Errorf("failed to update minor category: %w", err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid minor category id")
	}

	// Check if record exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("minor category with id %d not found", id)
		}
		return fmt.Errorf("failed to verify minor category: %w", err)
	}

	// Delete the minor category
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete minor category: %w", err)
	}
	return nil
}
