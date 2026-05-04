 package sales_tax_setup


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

func (s *Service) Create(ctx context.Context, d *SalesTaxSetup) error {
	// Validate input
	if d == nil {
		return errors.New("sales tax setup cannot be nil")
	}
	if d.Name == "" {
		return errors.New("sales tax setup name is required")
	}

	// Create the sales tax setup
	if err := s.repo.Create(ctx, d); err != nil {
		return fmt.Errorf("failed to create sales tax setup: %w", err)
	}
	return nil
}

func (s *Service) FindAll(ctx context.Context) ([]SalesTaxSetup, error) {
	salesTaxSetups, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sales tax setups: %w", err)
	}
	return salesTaxSetups, nil
}

func (s *Service) FindByID(ctx context.Context, id uint) (*SalesTaxSetup, error) {
	if id == 0 {
		return nil, errors.New("invalid sales tax setup id")
	}

	salesTaxSetup, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("sales tax setup with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch sales tax setup: %w", err)
	}
	return salesTaxSetup, nil
}

func (s *Service) Update(ctx context.Context, d *SalesTaxSetup) error {
	// Validate input
	if d == nil {
		return errors.New("sales tax setup cannot be nil")
	}
	if d.ID == 0 {
		return errors.New("sales tax setup id is required for update")
	}
	if d.Name == "" {
		return errors.New("sales tax setup name is required")
	}

	// Check if record exists
	_, err := s.repo.FindByID(ctx, d.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sales tax setup with id %d not found", d.ID)
		}
		return fmt.Errorf("failed to verify sales tax setup: %w", err)
	}

	// Update the sales tax setup
	if err := s.repo.Update(ctx, d); err != nil {
		return fmt.Errorf("failed to update sales tax setup: %w", err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid sales tax setup id")
	}

	// Check if record exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sales tax setup with id %d not found", id)
		}
		return fmt.Errorf("failed to verify sales tax setup: %w", err)
	}

	// Delete the sales tax setup
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete sales tax setup: %w", err)
	}
	return nil
}