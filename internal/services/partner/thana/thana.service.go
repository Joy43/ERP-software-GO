package thana

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, t *Thana) error {
	return s.repo.Create(ctx, t)
}

func (s *Service) FindAll(ctx context.Context) ([]Thana, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) FindByID(ctx context.Context, id uint) (*Thana, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) FindByDistrictID(ctx context.Context, districtID uint) ([]Thana, error) {
	return s.repo.FindByDistrictID(ctx, districtID)
}

func (s *Service) Update(ctx context.Context, t *Thana) error {
	return s.repo.Update(ctx, t)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
