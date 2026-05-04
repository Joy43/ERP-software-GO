package umo_measurement

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

func (s *Service) Create(ctx context.Context, d *Uom) error {

	if d == nil {
		return errors.New("uom cannot be nil")
	}
	if d.Name == "" {
		return errors.New("uom name is required")
	}

	// ---------------Create the uom---------
	if err := s.repo.Create(ctx, d); err != nil {
		return fmt.Errorf("failed to create uom: %w", err)
	}
	return nil
}

func (s *Service) FindAll(ctx context.Context) ([]Uom, error) {
	Uoms, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch uom: %w", err)
	}
	return Uoms, nil
}

func (s *Service) FindByID(ctx context.Context, id uint) (*Uom, error) {
	if id == 0 {
		return nil, errors.New("invalid uom id")
	}

	uomSetup, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("uom with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch uom: %w", err)
	}
	return uomSetup, nil
}

func (s *Service) Update(ctx context.Context, d *Uom) error {
	// ---------Validate input--------
	if d == nil {
		return errors.New("uom cannot be nil")
	}
	if d.ID == 0 {
		return errors.New("uom id is required for update")
	}
	if d.Name == "" {
		return errors.New("uom name is required")
	}

	// ------------- Check if record exists-----------
	_, err := s.repo.FindByID(ctx, d.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("uom with id %d not found", d.ID)
		}
		return fmt.Errorf("failed to verify uom: %w", err)
	}

	// Update the uom
	if err := s.repo.Update(ctx, d); err != nil {
		return fmt.Errorf("failed to update uom: %w", err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid uom id")
	}

	//------------- Check if record exists-----------
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("uom with id %d not found", id)
		}
		return fmt.Errorf("failed to verify uom: %w", err)
	}

	// ------------- Delete the uom ------------
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete uom: %w", err)
	}
	return nil
}
