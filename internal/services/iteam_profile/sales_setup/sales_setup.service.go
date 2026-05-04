package sales_setup

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"gorm.io/gorm"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// validateSalesSetup performs common validation for SalesSetup
func (s *Service) validateSalesSetup(d *SalesSetup, isUpdate bool) error {
	if d == nil {
		return errors.New("sales setup cannot be nil")
	}

	if isUpdate && d.ID == 0 {
		return errors.New("sales setup id is required for update")
	}

	if strings.TrimSpace(d.Name) == "" {
		return errors.New("sales setup name is required")
	}

	// Trim whitespace from name
	d.Name = strings.TrimSpace(d.Name)

	// Additional business rule: name should not exceed max length
	if len(d.Name) > 150 {
		return errors.New("sales setup name must not exceed 150 characters")
	}

	// Additional business rule: name should have at least 2 characters
	if len(d.Name) < 2 {
		return errors.New("sales setup name must have at least 2 characters")
	}

	return nil
}

// Create creates a new sales setup with full business logic
func (s *Service) Create(ctx context.Context, d *SalesSetup) error {
	// Validate input
	if err := s.validateSalesSetup(d, false); err != nil {
		return err
	}

	// Create the sales setup in the repository
	if err := s.repo.Create(ctx, d); err != nil {
		return fmt.Errorf("failed to create sales setup: %w", err)
	}

	return nil
}

// FindAll retrieves all sales setups with business logic
func (s *Service) FindAll(ctx context.Context) ([]SalesSetup, error) {
	salesSetups, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sales setups: %w", err)
	}

	// Return empty slice instead of nil if no records found
	if salesSetups == nil {
		return []SalesSetup{}, nil
	}

	return salesSetups, nil
}

// FindByID retrieves a sales setup by ID with business logic
func (s *Service) FindByID(ctx context.Context, id uint) (*SalesSetup, error) {
	// Validate ID
	if id == 0 {
		return nil, errors.New("invalid sales setup id")
	}

	// Fetch from repository
	salesSetup, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("sales setup with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch sales setup: %w", err)
	}

	return salesSetup, nil
}

// Update updates an existing sales setup with full business logic
func (s *Service) Update(ctx context.Context, d *SalesSetup) error {
	// Validate input
	if err := s.validateSalesSetup(d, true); err != nil {
		return err
	}

	// Check if record exists before updating
	existing, err := s.repo.FindByID(ctx, d.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("sales setup with id %d not found", d.ID)
		}
		return fmt.Errorf("failed to verify sales setup: %w", err)
	}

	// Preserve fields that shouldn't be updated
	if d.CreatedAt.IsZero() {
		d.CreatedAt = existing.CreatedAt
	}

	// Update the sales setup in the repository
	if err := s.repo.Update(ctx, d); err != nil {
		return fmt.Errorf("failed to update sales setup: %w", err)
	}

	return nil
}

// Delete removes a sales setup with business logic
func (s *Service) Delete(ctx context.Context, id uint) error {
	// Validate ID
	if id == 0 {
		return errors.New("invalid sales setup id")
	}

	// Check if record exists before deleting
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("sales setup with id %d not found", id)
		}
		return fmt.Errorf("failed to verify sales setup: %w", err)
	}

	// Delete the sales setup from the repository
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete sales setup: %w", err)
	}

	return nil
}
