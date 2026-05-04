package sales_supply_type

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

func (s *Service) Create(ctx context.Context, d *SalesSupplyType) error {
	// Validate input
	if d == nil {
		return errors.New("sales supply type cannot be nil")
	}
	if d.Name == "" {
		return errors.New("sales supply type name is required")
	}

	// Create the sales supply type
	if err := s.repo.Create(ctx, d); err != nil {
		return fmt.Errorf("failed to create sales supply type: %w", err)
	}
	return nil
}

func (s *Service) FindAll(ctx context.Context) ([]SalesSupplyType, error) {
	salesSupplyTypes, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sales supply types: %w", err)
	}
	return salesSupplyTypes, nil
}

func (s *Service) FindByID(ctx context.Context, id uint) (*SalesSupplyType, error) {
	if id == 0 {
		return nil, errors.New("invalid sales supply type id")
	}

	salesSupplyType, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("sales supply type with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch sales supply type: %w", err)
	}
	return salesSupplyType, nil
}

func (s *Service) Update(ctx context.Context, d *SalesSupplyType) error {
	// Validate input
	if d == nil {
		return errors.New("sales supply type cannot be nil")
	}
	if d.ID == 0 {
		return errors.New("sales supply type id is required for update")
	}
	if d.Name == "" {
		return errors.New("sales supply type name is required")
	}

	// Check if record exists
	_, err := s.repo.FindByID(ctx, d.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sales supply type with id %d not found", d.ID)
		}
		return fmt.Errorf("failed to verify sales supply type: %w", err)
	}

	// Update the sales supply type
	if err := s.repo.Update(ctx, d); err != nil {
		return fmt.Errorf("failed to update sales supply type: %w", err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid sales supply type id")
	}

	// Check if record exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sales supply type with id %d not found", id)
		}
		return fmt.Errorf("failed to verify sales supply type: %w", err)
	}

	// Delete the sales supply type
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete sales supply type: %w", err)
	}
	return nil
}