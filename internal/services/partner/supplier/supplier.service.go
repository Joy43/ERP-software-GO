package supplier

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, sup *Supplier) error {
	return s.repo.Create(ctx, sup)
}

func (s *Service) FindAll(ctx context.Context) ([]Supplier, error) {
	return s.repo.FindAllWithPreload(ctx)
}

func (s *Service) FindByID(ctx context.Context, id uint) (*Supplier, error) {
	return s.repo.FindByIDWithPreload(ctx, id)
}

func (s *Service) Update(ctx context.Context, sup *Supplier) error {
	return s.repo.Update(ctx, sup)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
