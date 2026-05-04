package sales_representative

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, sr *SalesRepresentative) error {
	return s.repo.Create(ctx, sr)
}

func (s *Service) Update(ctx context.Context, sr *SalesRepresentative) error {
	return s.repo.Update(ctx, sr)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) FindByID(ctx context.Context, id uint) (*SalesRepresentative, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) FindAll(ctx context.Context) ([]SalesRepresentative, error) {
	return s.repo.FindAll(ctx)
}
