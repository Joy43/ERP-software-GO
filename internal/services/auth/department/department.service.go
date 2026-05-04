package department

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, d *Department) error {
	return s.repo.Create(ctx, d)
}

func (s *Service) FindAll(ctx context.Context) ([]Department, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) FindByID(ctx context.Context, id uint) (*Department, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, d *Department) error {
	return s.repo.Update(ctx, d)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
