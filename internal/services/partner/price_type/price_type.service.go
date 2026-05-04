package price_type

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, pt *PriceType) error {
	return s.repo.Create(ctx, pt)
}

func (s *Service) Update(ctx context.Context, pt *PriceType) error {
	return s.repo.Update(ctx, pt)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) FindByID(ctx context.Context, id uint) (*PriceType, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) FindAll(ctx context.Context) ([]PriceType, error) {
	return s.repo.FindAll(ctx)
}
